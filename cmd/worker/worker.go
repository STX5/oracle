package main

import (
	smartcontract "oracle/smartContract"
	"oracle/src/worker"
)

func main() {
	var endpoints = []string{"localhost:2379"}
	prefix := "01"
	id := "011"
	woker, _ := worker.NewWoker(id, prefix, endpoints, smartcontract.TestWriter{})
	defer woker.Close()

	go woker.GetJobs()
	woker.Work()
}
