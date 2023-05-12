package main

import (
	jobdeamon "oracle/src/jobDeamon"
	"os"
	"strconv"
)

var endpoints = []string{"localhost:2379"}

func main() {
	var port int

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	} else if len(os.Args) == 1 {
		port = 8080
	}

	jd, _ := jobdeamon.NewJobDeamon(endpoints)
	defer jd.Close()

	jd.Run(port)
}
