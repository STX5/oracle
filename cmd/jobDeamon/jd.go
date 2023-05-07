package main

import (
	jobdeamon "oracle/src/jobDeamon"
	"os"
	"strconv"
)

var (
	endpoints = []string{"localhost:2379"}
	jd        *jobdeamon.JobDeamon
)

func main() {
	var port int

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	} else if len(os.Args) == 1 {
		port = 8080
	}

	jd, _ = jobdeamon.NewJobDeamon(endpoints)
	defer jd.Close()

	// to test alterEndpoints
	/*go func(t time.Duration) {
		time.Sleep(t)
		jd.AlterEndpoints(endpoints)
	}(5 * time.Second)*/

	jd.Run(port)
}
