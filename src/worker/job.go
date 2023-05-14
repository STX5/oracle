package worker

import (
	"context"
	"io"

	//"io"
	"log"
	"net/http"
)

type Job struct {
	Cancel context.CancelFunc
	// ETCD key, length : 160
	ID string
	JobVal
}

// JobVal 修改了JobVal的字段值，增加了字段JobFrom
type JobVal struct {
	URL     string `json:"url"`
	Pattern string `json:"pattern"`
	// SM OracleWiter related, not sure yet
}

// TODO: add timeout
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
