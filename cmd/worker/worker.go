package main

import (
	"encoding/json"
	"log"
	"net/http"
	smartcontract "oracle/smartContract"
	"oracle/src/worker"
	"os"
	"strconv"
)

const (
	legalPrefix = "1"
	legalHexID  = "99c82bb73505a3c0b453f9fa0e881d6e5a32a0c1"
)

var (
	endpoints = []string{"localhost:2379"}
	woker     *worker.Worker
)

func main() {
	var port int

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	} else if len(os.Args) == 1 {
		port = 8080
	}

	woker, _ = worker.NewWoker(legalHexID, legalPrefix, endpoints, smartcontract.TestWriter{})
	defer woker.Close()

	// 测试修改endpoints
	// to test alterPrefix
	/*go func(t time.Duration) {
		time.Sleep(t)
		woker.AlterEndpoints(endpoints)
		woker.AlterPrefix("000")
	}(5 * time.Second)*/

	// 开启http服务
	go startHttpServer(port)
	woker.Run()
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

type WorkerConfig struct {
	Prefix   string   `json:"prefix"`
	Endpoint []string `json:"endpoint"`
}

func updateConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config WorkerConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	woker.AlterPrefix(config.Prefix)
	woker.AlterEndpoints(config.Endpoint)

	log.Printf("Update config success, now prefix: %s, endpoint: %v", config.Prefix, config.Endpoint)

	w.WriteHeader(http.StatusNoContent)
}
