package deamon

import (
	"context"
	"log"
	"os"
	"runtime"
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

// if a ETCD lock expires and its key
// still exist, means a worker failed.
// JobDeamon should update this key's version,
// so other wokers will receive watch event and do it again.
//
// JobDeamon should be deployed with every ETCD member node
type JobDeamon struct {
	mu         *sync.Mutex
	ETCDClient *clientv3.Client
	// WatchLock() wirtes deleteEvent to WatcherChan;
	// ProtectJob() read from WatcherChan to check job status
	WatcherChan      chan *clientv3.Event
	Logger           *logrus.Logger
	alterEndpointsCh chan []string
	CloseOnce        *sync.Once
}

func NewJobDeamon(endpoints []string) (*JobDeamon, error) {
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
	return &JobDeamon{
		ETCDClient:       cli,
		Logger:           logger,
		WatcherChan:      make(chan *clientv3.Event),
		alterEndpointsCh: make(chan []string),
		mu:               new(sync.Mutex),
		CloseOnce:        new(sync.Once),
	}, nil
}

func (jd *JobDeamon) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	go jd.WatchLock(ctx)
	go jd.ProtectJob(ctx)
	for newEndPoints := range jd.alterEndpointsCh {
		jd.mu.Lock()
		cancel()
		jd.Logger.WithFields(logrus.Fields{
			"newEndpoints": newEndPoints,
		}).Info("Endpoints Changed")
		jd.Close()

		var err error
		// TODO: need to check endpoints first
		jd.ETCDClient, err = clientv3.New(clientv3.Config{
			Endpoints:   newEndPoints,
			DialTimeout: 5 * time.Second,
		})
		if err != nil {
			log.Fatal("Bad ETCD Endpoints")
		}
		jd.WatcherChan = make(chan *clientv3.Event)
		jd.CloseOnce = new(sync.Once)
		ctx, cancel = context.WithCancel(context.Background())
		go jd.WatchLock(ctx)
		go jd.ProtectJob(ctx)
		jd.mu.Unlock()
	}
	cancel()
}

func (jd JobDeamon) WatchLock(ctx context.Context) {
	ch := jd.ETCDClient.Watch(ctx, "lock/",
		clientv3.WithPrefix())
	for resp := range ch {
		for _, event := range resp.Events {
			switch event.Type {
			case mvccpb.DELETE:
				jd.WatcherChan <- event
			default:
				continue
			}
		}
	}
}

func (jd JobDeamon) ProtectJob(ctx context.Context) {
	// use token to limit goroutine number
	token := make(chan struct{}, paraIndex)
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-jd.WatcherChan:
			if !ok {
				return
			}
			token <- struct{}{}
			go func(event *clientv3.Event) {
				defer func() {
					<-token
				}()

				jobID := string(event.Kv.Key)
				jobID = jobID[5:]
				jd.Logger.WithFields(logrus.Fields{
					"JobID": jobID,
				}).Info("Detected Lock Delete Event")

				getResp, err := jd.ETCDClient.Get(context.Background(), jobID)
				if err != nil {
					jd.Logger.WithFields(logrus.Fields{
						"JobID":      jobID,
						"ErrMessage": err,
					}).Warning("Error Get")
				}
				if len(getResp.Kvs) == 0 { // no need to re-put job
					jd.Logger.WithFields(logrus.Fields{
						"JobID": jobID,
					}).Info("Job Successfully Done")
					return
				}
				// a job can be tried for at most 5 times
				if getResp.Kvs[0].Version >= 5 {
					jd.Logger.WithFields(logrus.Fields{
						"JobID": jobID,
					}).Warn("This is An Impossible Job")
					return
				}

				jd.Logger.WithFields(logrus.Fields{
					"JobID": jobID,
				}).Info("Get JobID")
				jobVal := string(getResp.Kvs[0].Value)
				txn := jd.ETCDClient.Txn(context.Background())
				txn.If(clientv3.Compare(clientv3.Value(jobID), "=", jobVal)).
					Then(clientv3.OpPut(jobID, string(jobVal)))

				_, err = txn.Commit()
				if err != nil {
					jd.Logger.WithFields(logrus.Fields{
						"JobID":      jobID,
						"ErrMessage": err,
					}).Warning("Error Excuting ProtectJob TXN")
					return
				}
			}(event)
		}
	}
}

func (jd JobDeamon) AlterEndpoints(newEndpoints []string) {
	jd.alterEndpointsCh <- newEndpoints
}

func (jd JobDeamon) Close() {
	jd.CloseOnce.Do(func() {
		jd.ETCDClient.Close()
		close(jd.WatcherChan)
	})
}
