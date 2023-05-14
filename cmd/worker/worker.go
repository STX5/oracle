package main

import (
	"fmt"
	"oracle/src/worker"
	"os"
	"strconv"
)

const (
	legalPrefix = "1"
	legalHexID  = "99c82bb73505a3c0b453f9fa0e881d6e5a32a0c1"
)

type TestWriter struct {
}

func (t TestWriter) WriteData(data string) (bool, error) {
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

	woker, _ := worker.NewWoker(legalHexID, legalPrefix, endpoints, &TestWriter{})
	defer woker.Close()
	woker.Run(port)
}
