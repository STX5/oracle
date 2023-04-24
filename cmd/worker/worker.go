package main

import (
	smartcontract "oracle/smartContract"
	"oracle/src/worker"
	"time"
)

const legalPrefix = "1"
const legalHexID = "99c82bb73505a3c0b453f9fa0e881d6e5a32a0c1"

func main() {
	var endpoints = []string{"localhost:2379"}
	woker, _ := worker.NewWoker(legalHexID, legalPrefix, endpoints, smartcontract.TestWriter{})
	defer woker.Close()
	// to test alterPrefix
	go func(t time.Duration) {
		time.Sleep(t)
		woker.AlterEndpoints(endpoints)
		woker.AlterPrefix("000")
	}(5 * time.Second)
	woker.Run()
}
