package jobdeamon

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// if a ETCD lock expires and its key
// still exist, means a worker failed.
// JobDeamon should update this key's version,
// so other wokers will receive watch event and do it again.
//
// JobDeamon should be deployed with every ETCD member node
type JobDeamon struct {
	ETCDClient *clientv3.Client
	// WatchLock() wirtes deleteEvent to WatcherChan;
	// ProtectJob() read from WatcherChan to check job status
	WatcherChan chan *clientv3.Event
	Logger      *logrus.Logger
}

func NewJobDeamon(endpoints []string) (*JobDeamon, error) {
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
	return &JobDeamon{
		ETCDClient:  cli,
		Logger:      logger,
		WatcherChan: make(chan *clientv3.Event),
	}, nil
}

func (jd JobDeamon) WatchLock() {
	ch := jd.ETCDClient.Watch(context.Background(), "lock/",
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

// TODO: limit goroutine number
func (jd JobDeamon) ProtectJob() {
	for {
		event := <-jd.WatcherChan
		go func(event *clientv3.Event) {
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
