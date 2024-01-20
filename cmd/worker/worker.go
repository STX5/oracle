package main

import (
	"fmt"
	smartcontract "oracle/smartContract"
	"oracle/src/worker"
	"os"
	"strconv"
)

const (
	legalPrefix = "1"
	legalHexID  = "99c82bb73505a3c0b453f9fa0e881d6e5a32a0c1"
)

type TestWriter struct{}

func (t TestWriter) WriteData(jobID string, data string) (bool, error) {
	fmt.Println(data)
	return true, nil
}

var endpoints = []string{"localhost:2379"}

func main() {
	var port int

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	} else if len(os.Args) == 1 {
		port = 8080
	}
	var oracle, err = smartcontract.NewOracle()
	if err != nil {
		panic(err)
	}
	woker, _ := worker.NewWoker(legalHexID, legalPrefix, endpoints, oracle)
	// ow := TestWriter{}
	// woker, _ := worker.NewWoker(legalHexID, legalPrefix, endpoints, ow)

	defer woker.Close()
	go oracle.Run()
	woker.Run(port)
}
