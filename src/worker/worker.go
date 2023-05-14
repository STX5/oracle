package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	smartContract "oracle/smartContract"
	"oracle/util"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var paraIndex int

func init() {
	paraIndex = runtime.NumCPU()
}

// TODO: add ETH CLIENT(Maybe OracleWriter will include it)
type Worker struct {
	// GroupPrefix indecates which group the worker belongs to,
	// which determines the key range the worker watches
	mu          *sync.Mutex
	GroupPrefix string
	ID          string
	ETCDClient  *clientv3.Client
	Logger      *logrus.Logger
	// GetJobs() wirtes watch response to WatcherChan;
	// Work() read from WatcherChan to deal with jobs
	WatcherChan chan *Job
	// write to alterPrefixCh when alter GroupPrefix
	alterPrefixCh    chan string
	alterEndpointsCh chan []string
	CloseOnce        *sync.Once

	/*
		OracleWriter writes job result to Oracle smart contract
		TODO: add ETH CLIENT to OracleWriter
	*/
	smartContract.OracleWriter
}

func NewWoker(id string, prefix string, endpoints []string, ow smartContract.OracleWriter) (*Worker, error) {
	id, err := util.DecodeHex(id)
	if err != nil || len(id) != 160 {
		return nil, err
	}
	prefix, err = util.DecodeHex(prefix)
	if err != nil {
		return nil, err
	}
	legal := util.CheckPrefix(prefix, id)
	if !legal {
		return nil, fmt.Errorf("prefix:%s, and ID:%s ,not match", prefix, id)
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat:        "2006-01-02 15:04:05",
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		},
	}
	return &Worker{
		ID:               id,
		mu:               new(sync.Mutex),
		GroupPrefix:      prefix,
		ETCDClient:       cli,
		Logger:           logger,
		CloseOnce:        new(sync.Once),
		WatcherChan:      make(chan *Job),
		alterPrefixCh:    make(chan string),
		alterEndpointsCh: make(chan []string),
		OracleWriter:     ow,
	}, nil
}

func (worker Worker) Run(port int) {
	ctx, cancel := context.WithCancel(context.Background())
	go worker.GetJobs(ctx)
	go worker.Work(ctx)
	go StartHttpServer(&worker, port)
	for {
		select {
		case newPrefix := <-worker.alterPrefixCh:
			worker.mu.Lock()
			oldPrefix := worker.GroupPrefix
			// no need to alter prefix
			if oldPrefix == newPrefix {
				worker.mu.Unlock()
				continue
			}
			worker.GroupPrefix = newPrefix
			worker.Logger.WithFields(logrus.Fields{
				"newPrefix": newPrefix,
				"oldPrefix": oldPrefix,
			}).Info("Prefix Changed")
			// call cancel() to stop GetJob(), Work()
			cancel()
			ctx, cancel = context.WithCancel(context.Background())
			go worker.GetJobs(ctx)
			go worker.Work(ctx)
			worker.mu.Unlock()
		case newEndPoints := <-worker.alterEndpointsCh:
			worker.mu.Lock()

			cancel()
			worker.Logger.WithFields(logrus.Fields{
				"newEndpoints": newEndPoints,
			}).Info("Endpoints Changed")
			worker.Close()

			var err error
			// TODO: need to check endpoints first
			worker.ETCDClient, err = clientv3.New(clientv3.Config{
				Endpoints:   newEndPoints,
				DialTimeout: 5 * time.Second,
			})
			if err != nil {
				log.Fatal("Bad ETCD Endpoints")
			}
			worker.WatcherChan = make(chan *Job)
			worker.CloseOnce = new(sync.Once)
			ctx, cancel = context.WithCancel(context.Background())
			go worker.GetJobs(ctx)
			go worker.Work(ctx)
			worker.mu.Unlock()
		}
	}
}

// GetJobs() wirtes watch result to WatcherChan
func (woker Worker) GetJobs(ctx context.Context) {
	ch := woker.ETCDClient.Watch(ctx, woker.GroupPrefix,
		clientv3.WithPrefix())

	for resp := range ch {
		for _, event := range resp.Events {
			switch event.Type {
			case mvccpb.PUT:
				modVersion := event.Kv.ModRevision
				// this is a delete action, other worker has done the job
				if modVersion == 0 {
					continue
				}
				cancelFunc, locked, err := woker.acquireLock(string(event.Kv.Key), 5)
				if !locked || err != nil {
					woker.Logger.WithFields(logrus.Fields{
						"JobID":      string(event.Kv.Key),
						"ErrMessage": err,
					}).Info("Acquire Lock Failed")
					continue
				}
				var jobVal JobVal
				err = json.Unmarshal(event.Kv.Value, &jobVal)
				if err != nil {
					woker.Logger.WithFields(logrus.Fields{
						"JobID":      string(event.Kv.Key),
						"ErrMessage": err,
					}).Info("Error unmarshalling JobVal")
				}
				woker.WatcherChan <- &Job{
					Cancel: cancelFunc,
					ID:     string(event.Kv.Key),
					JobVal: jobVal,
				}
			default:
				continue
			}
		}
	}
}

func (woker Worker) Work(ctx context.Context) {
	// use token to limit goroutine number
	token := make(chan struct{}, paraIndex)
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-woker.WatcherChan:
			if !ok {
				return
			}
			token <- struct{}{}
			go func(job *Job) {
				defer func() {
					<-token
				}()
				woker.Logger.WithFields(logrus.Fields{
					"JobID": job.ID,
				}).Info("Get Job")
				// whether success or not, release the Lock
				defer woker.releaseLock(job.Cancel, job.ID)
				data, err := job.Scrap()
				if err != nil {
					woker.Logger.WithFields(logrus.Fields{
						"JobID":      job.ID,
						"ErrMessage": err,
					}).Warning("Error Scraping")
					return
				}
				time.Sleep(6 * time.Second) // this is for testing
				// if success, delete job
				defer woker.ETCDClient.Delete(context.Background(), job.ID)
				// 这里做了修改，写回数据的时候，需要将job.JobFrom字段传回去，智能合约需要
				woker.OracleWriter.WriteData(data)
			}(job)
		}
	}
}

func (woker Worker) acquireLock(key string, ttl int) (context.CancelFunc, bool, error) {
	resp, err := woker.ETCDClient.Grant(context.Background(), int64(ttl))
	if err != nil {
		return nil, false, err
	}
	leaseID := resp.ID

	ctx, cancel := context.WithCancel(context.Background())
	ch, keepErr := woker.ETCDClient.KeepAlive(ctx, leaseID)
	if keepErr != nil {
		woker.Logger.WithFields(logrus.Fields{
			"JobID":      string(key),
			"ErrMessage": keepErr,
		}).Warning("Error Keep Alive")
		cancel()
		return nil, false, keepErr
	}

	go func(key string) {
		for {
			select {
			case ka := <-ch:
				woker.Logger.WithFields(logrus.Fields{
					"JobID":        string(key),
					"KeepAliveTTL": ka.TTL,
				}).Info("KeepAlive")
			case <-ctx.Done():
				woker.Logger.WithFields(logrus.Fields{
					"JobID": string(key),
				}).Info("Release Lease")
				return
			}
		}
	}(key)

	txn := woker.ETCDClient.Txn(context.Background())
	txn.If(clientv3.Compare(clientv3.CreateRevision("lock/"+key), "=", 0)).
		Then(clientv3.OpPut("lock/"+key, "JOB", clientv3.WithLease(leaseID)))

	txnResp, err := txn.Commit()
	if err != nil {
		woker.Logger.WithFields(logrus.Fields{
			"JobID":      string(key),
			"ErrMessage": err,
		}).Warning("Error Excuting Acquire Lock TXN")
		cancel()
		return nil, false, fmt.Errorf("acquire lock txn failed, err: %s", err)
	}
	if !txnResp.Succeeded { // Dose not get the lock
		cancel()
		return nil, false, fmt.Errorf("didnt get lock for JobID: %s", key)
	}
	return cancel, true, nil
}

func (woker Worker) releaseLock(cancel context.CancelFunc, key string) {
	cancel()
	txn := woker.ETCDClient.Txn(context.Background())
	txn.If(clientv3.Compare(clientv3.Value("lock/"+key), "=", "JOB")).
		Then(clientv3.OpDelete("lock/" + key))

	txnResp, err := txn.Commit()
	if err != nil {
		woker.Logger.WithFields(logrus.Fields{
			"JobID":      string(key),
			"ErrMessage": err,
		}).Warning("Error Excuting Release Lock TXN")
	}
	// don't need more action, cause called cancel() in the first place
	if !txnResp.Succeeded {
		woker.Logger.WithFields(logrus.Fields{
			"JobID": string(key),
		}).Warning("Failed to Releasing Lock")
		return
	}
}

func (worker Worker) AlterPrefix(newPrefix string) {
	worker.alterPrefixCh <- newPrefix
}

func (worker Worker) AlterEndpoints(newEndpoints []string) {
	worker.alterEndpointsCh <- newEndpoints
}

func (woker Worker) Close() {
	woker.CloseOnce.Do(func() {
		woker.ETCDClient.Close()
		close(woker.WatcherChan)
	})
}

func StartHttpServer(worker *Worker, port int) {
	// 更新配置的路由
	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var config WorkerConfig
		err := json.NewDecoder(r.Body).Decode(&config)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		worker.AlterPrefix(config.Prefix)
		worker.AlterEndpoints(config.Endpoint)

		/*worker.Logger.WithFields(logrus.Fields{
			"newPrefix":   config.Prefix,
			"newEndpoint": config.Endpoint,
		}).Info("Update config success")*/

		w.WriteHeader(http.StatusNoContent)
	})

	// 尝试监听端口
	for ; port < 65535; port++ {
		worker.Logger.WithFields(logrus.Fields{
			"port": port,
		}).Info("Start http server")
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			worker.Logger.WithFields(logrus.Fields{
				"ServerPort": port,
				"ErrMsg":     err,
			}).Info("Start http server failed")
			continue
		}
	}
}

type WorkerConfig struct {
	Prefix   string   `json:"prefix"`
	Endpoint []string `json:"endpoint"`
}
