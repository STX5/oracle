package woker

import (
	"context"
	"encoding/json"
	"fmt"
	smartContract "oracle/smartContract"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// TODO: add ETH CLIENT(Maybe OracleWriter will include it)
type Worker struct {
	// GroupPrefix indecates which group the worker belongs to,
	// which determines the key range the worker watches
	GroupPrefix string
	ID          string
	ETCDClient  *clientv3.Client
	Logger      *logrus.Logger
	// GetJobs() wirtes watch response to WatcherChan;
	// Work() read from WatcherChan to deal with jobs
	WatcherChan chan *Job
	CloseOnce   *sync.Once

	/*
		OracleWriter writes job result to Oracle smart contract
		TODO: add ETH CLIENT to OracleWriter
	*/
	smartContract.OracleWriter
}

func NewWoker(id string, prefix string, endpoints []string, ow smartContract.OracleWriter) (*Worker, error) {
	// if len(id) != 160 {
	// 	return nil, fmt.Errorf("%s", "id length != 160")
	// }
	// if len(prefix) >= 160 || len(prefix) < 1 {
	// 	return nil, fmt.Errorf("%s", "Illegal prefix")
	// }
	// // need to do more check, eg. prefix can't match with id
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	return &Worker{
		ID:           id,
		GroupPrefix:  prefix,
		ETCDClient:   cli,
		Logger:       logger,
		CloseOnce:    new(sync.Once),
		WatcherChan:  make(chan *Job),
		OracleWriter: ow,
	}, nil
}

// GetJobs() wirtes watch result to WatcherChan
func (woker Worker) GetJobs() {
	ch := woker.ETCDClient.Watch(context.Background(), woker.GroupPrefix,
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

// create goroutines to deal with jobs
// TODO: limit goroutine number
func (woker Worker) Work() {
	for {
		job := <-woker.WatcherChan
		go func(job *Job) {
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
			woker.OracleWriter.WriteData(data)
		}(job)
	}
}

func (woker Worker) Close() {
	woker.CloseOnce.Do(func() {
		woker.ETCDClient.Close()
	})
}
