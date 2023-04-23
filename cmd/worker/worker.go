package main

import (
	smartcontract "oracle/smartContract"
	"oracle/src/worker"
)

const legalPrefix = "99"
const legalHexID = "99c82bb73505a3c0b453f9fa0e881d6e5a32a0c1"

func main() {
	var endpoints = []string{"localhost:2379"}
	woker, _ := worker.NewWoker(legalHexID, legalPrefix, endpoints, smartcontract.TestWriter{})
	defer woker.Close()

	go woker.GetJobs()
	woker.Work()
}
