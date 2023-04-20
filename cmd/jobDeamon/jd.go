package main

import jobdeamon "oracle/src/jobDeamon"

func main() {
	var endpoints = []string{"localhost:2379"}
	jd, _ := jobdeamon.NewJobDeamon(endpoints)
	go jd.WatchLock()
	jd.ProtectJob()
}
