package main

import (
	jobdeamon "oracle/src/jobDeamon"
	"time"
)

func main() {
	var endpoints = []string{"localhost:2379"}
	jd, _ := jobdeamon.NewJobDeamon(endpoints)

	// to test alterEndpoints
	go func(t time.Duration) {
		time.Sleep(t)
		jd.AlterEndpoints(endpoints)
	}(5 * time.Second)

	jd.Run()
}
