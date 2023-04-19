package main

import (
	smartcontract "oracle/smartContract"
	"oracle/woker"
)

func main() {
	var endpoints = []string{"localhost:2379"}
	prefix := "01"
	id := "011"
	// jd, _ := jobdeamon.NewJobDeamon(endpoints)
	// go jd.WatchLock()
	// go jd.ProtectJob()

	woker, _ := woker.NewWoker(id, prefix, endpoints, smartcontract.TestWriter{})
	defer woker.Close()
	go woker.GetJobs()
	woker.Work()
}
