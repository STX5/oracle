package woker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"oracle/smartcontract"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Job struct {
	// length : 160
	ID string
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
	return string(body), nil
}

// TODO: add ETH CLIENT
type Worker struct {
	// length : 160
	ID string
	// GroupPrefix indecates which group the worker belongs to,
	// which determines the key range the worker watchs
	GroupPrefix string
	// // OracleWriter writes job result to Oracle smart contract
	// OracleWriter OracleWiter
	// ETCD client, Watch in specific key range
	ETCDClient *clientv3.Client
	// GetJobs() wirtes watch response to this channel
	// Work() read from this channel to deal with jobs
	WatcherChan chan *Job
	// TODO: add ETH CLIENT
	smartcontract.OracleWriter
}

func NewWoker(id string, prefix string, endpoints []string) (*Worker, error) {
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
		OracleWriter: nil,
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
				// TODO: get ETCD distributed lock
				var jobVal JobVal
				err := json.Unmarshal(event.Kv.Value, &jobVal)
				if err != nil {
					log.Println("Error unmarshalling JobVal")
				}
				woker.WatcherChan <- &Job{
					ID:     string(event.Kv.Key),
					JobVal: jobVal,
				}
			}
		}
	}
}

// create goroutines to deal with jobs
func (woker Worker) Work() {
	for {
		job := <-woker.WatcherChan
		go func(job *Job) {
			data, err := job.Scrap()
			if err != nil {
				return
			}
			woker.OracleWriter.WriteData(data)
		}(job)
	}
}
