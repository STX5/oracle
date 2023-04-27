// Package smartcontract 和etcd操作相关的
package smartcontract

import (
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

// etcd的单例客户端对象
var clientInstance *clientv3.Client

var once sync.Once

// GetEtcdClient 根据传入的etcd的地址返回etcd的客户端连接
// 返回的是客户端的单例
func GetEtcdClient(endpoints []string) *clientv3.Client {
	once.Do(func() {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
		})

		if err != nil {
			panic("创建etcd客户端连接错误")
		}

		clientInstance = cli
	})
	return clientInstance
}
