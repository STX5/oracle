package woker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	smartContract "oracle/smartContract"
	"sync"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// TODO: add ETH CLIENT(Maybe OracleWriter will include it)
type Worker struct {
	// GroupPrefix indecates which group the worker belongs to,
	// which determines the key range the worker watchs
	GroupPrefix string
	ID          string
	ETCDClient  *clientv3.Client

	// GetJobs() wirtes watch response to this channel
	// Work() read from this channel to deal with jobs
	WatcherChan chan *Job
	CloseOnce   *sync.Once

	/*
		OracleWriter writes job result to Oracle smart contract
		TODO: add ETH CLIENT to OracleWriter
	*/
	smartContract.OracleWriter
}

func NewWoker(id string, prefix string, endpoints []string, ow smartContract.OracleWriter) (*Worker, error) {
	if len(id) != 160 {
		return nil, fmt.Errorf("%s", "id length != 160")
	}
	if len(prefix) >= 160 || len(prefix) < 1 {
		return nil, fmt.Errorf("%s", "Illegal prefix")
	}
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &Worker{
		ID:           id,
		GroupPrefix:  prefix,
		ETCDClient:   cli,
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
					continue
				}
				var jobVal JobVal
				err = json.Unmarshal(event.Kv.Value, &jobVal)
				if err != nil {
					log.Println("Error unmarshalling JobVal")
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
	defer cancel()
	defer woker.ETCDClient.Revoke(context.Background(), leaseID)
	ch, keepErr := woker.ETCDClient.KeepAlive(ctx, leaseID)
	if keepErr != nil {
		fmt.Println("keep alive failed, err:", keepErr)
		return nil, false, keepErr
	}

	go func(key string) {
		for {
			select {
			case ka := <-ch:
				fmt.Printf("lease ttl is: %d for key: %s:", ka.TTL, key)
			case <-ctx.Done():
				fmt.Printf("Job done for Job ID: %s:", key)
				return
			}
		}
	}(key)

	txn := woker.ETCDClient.Txn(context.Background())
	txn.If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).
		Then(clientv3.OpPut("lock/"+key, "JOB", clientv3.WithLease(leaseID)))

	txnResp, err := txn.Commit()
	if err != nil {
		return nil, false, fmt.Errorf("txn failed, err: %s", err)
	}
	if !txnResp.Succeeded { // Dose not get the lock
		return nil, false, fmt.Errorf("lock failed for JObID: %s", key)
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
		fmt.Printf("txn failed, err: %s", err)
	}
	if !txnResp.Succeeded {
		fmt.Printf("failed for JObID: %s", key)
		return
	}
}

// create goroutines to deal with jobs
func (woker Worker) Work() {
	for {
		job := <-woker.WatcherChan
		go func(job *Job) {
			// whether success or not, release the Lock
			defer woker.releaseLock(job.Cancel, job.ID)
			data, err := job.Scrap()
			if err != nil {
				return
			}
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
