package worker

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Job struct {
	Cancel context.CancelFunc
	// ETCD key, length : 160
	ID string
	JobVal
}

type JobVal struct {
	URL     string `json:"url"`
	Pattern string `json:"pattern"`
	// SM OracleWiter related, not sure yet
}

// TODO: add timeout
func timeout() {
	n := time.Duration(3)
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*n)
	defer cancel()
	ch := make(chan struct{}, 0)
	go func() {
		Scrap()
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}

func Scrap() {
	panic("unimplemented")
}

func (j Job) Scrap() (string, error) {
	log.Println("start scraping")
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
