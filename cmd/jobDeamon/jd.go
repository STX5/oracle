package main

import (
	"encoding/json"
	"log"
	"net/http"
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

	go startHttpServer(port)
	jd.Run()
}

type JobConfig struct {
	Endpoint []string `json:"endpoint"`
}

func startHttpServer(port int) {
	// 更新配置的路由
	http.HandleFunc("/update", updateConfig)

	// 尝试监听端口
	for ; port < 65535; port++ {
		log.Printf("Start http server, port: %d", port)
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			log.Printf("Start http server failed, port: %d, err: %v", port, err)
			continue
		}
	}
}

func updateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config JobConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jd.AlterEndpoints(config.Endpoint)

	log.Printf("Update endpoints: %v", config.Endpoint)

	w.WriteHeader(http.StatusNoContent)
}
