package main

import (
	"oracle/woker"
)

func main() {
	var endpoints = []string{"localhost:2379"}
	prefix := "01"
	id := "THIS IS AN INVALID ID"
	woker, _ := woker.NewWoker(id, prefix, endpoints)
	woker.GetJobs()
	woker.Work()

	// defer woker.Close()
}
