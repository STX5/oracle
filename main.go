package main

import (
	"oracle/woker"
)

func main() {
	var endpoints = []string{"localhost:2379"}
	prefix := "01"
	id := "THIS IS AN INVALID ID"
	// nil OracleWriter
	woker, _ := woker.NewWoker(id, prefix, endpoints, nil)
	defer woker.Close()
	go woker.GetJobs()
	woker.Work()
}
