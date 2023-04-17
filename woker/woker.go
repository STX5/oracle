package woker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	smartContract "oracle/smartContract"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3/concurrency"
)

type Job struct {
	// length : 160
	Cancel context.CancelFunc
	ID     string
	JobVal
}

type JobVal struct {
	URL     string
	Pattern string
	//SM OracleWiter
}

func (j Job) Scrap() (string, error) {
	res, err := http.Get(j.URL)
	if err != nil {
		return "", err
	}
	data, err := j.resolve(res)
	if err != nil {
		return "", err
	}
	return data, nil
}

// Not implemented
// TODO: add resolver
func (j Job) resolve(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// TODO: ADD RESOLVER
	return string(body), nil
}

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
	CloseChan   chan struct{}

	/*
		OracleWriter writes job result to Oracle smart contract
	*/
	// TODO: add ETH CLIENT
	smartContract.OracleWriter
}

func NewWoker(id string, prefix string, endpoints []string, ow smartContract.OracleWriter) (*Worker, error) {
	if len(id) != 160 {
		return nil, fmt.Errorf("%s", "id length != 160")
	}
	if len(prefix) >= len(id) {
		return nil, fmt.Errorf("%s", "prefix length can not be larger than id")
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
				if modVersion == 0 { // this is a delete action, other worker has done the job
					continue
				}
				cancelFunc, locked, err := woker.acquireLock(event.Kv.Key, 5)
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

func (woker Worker) acquireLock(key []byte, ttl int) (context.CancelFunc, bool, error) {
	resp, err := woker.ETCDClient.Grant(context.Background(), int64(ttl))
	if err != nil {
		return nil, false, err
	}
	leaseID := resp.ID
	session, err := concurrency.NewSession(woker.ETCDClient, concurrency.WithLease(leaseID))
	if err != nil {
		return nil, false, err
	}
	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Duration(ttl))
	mutex := concurrency.NewMutex(session, "/lock")
	err = mutex.Lock(timeout)
	if err != nil {
		defer cancelFunc()
		return nil, false, err
	}
	return cancelFunc, true, nil
}

func releaseLock(cancel context.CancelFunc) {
	cancel()
}

// create goroutines to deal with jobs
func (woker Worker) Work() {
	for {
		job := <-woker.WatcherChan
		go func(job *Job) {
			// release Lock
			defer releaseLock(job.Cancel)
			data, err := job.Scrap()
			if err != nil {
				return
			}
			woker.OracleWriter.WriteData(data)
		}(job)
	}
}

// TODO: use CloseChan to inform other goroutines
func (woker Worker) Close() {
	woker.ETCDClient.Close()
}
