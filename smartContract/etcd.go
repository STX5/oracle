// Package smartcontract 和etcd操作相关的
package smartcontract

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// 定义etcd的客户端
type etcdClient struct {
	// 客户端对象
	client *clientv3.Client
	// 连接超时时间
	timeout time.Duration
	urls    []string
}

// etcd的单例客户端对象
var etcdClientInstance *etcdClient

// 获取etcd客户端的单例
func getEtcdClientInstance(urls []string, timeout time.Duration) *etcdClient {
	once.Do(func() {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   urls,
			DialTimeout: timeout * time.Second,
		})

		if err != nil {
			log.Fatal("etcd客户端连接错误")
		}

		etcdClientInstance = new(etcdClient)
		etcdClientInstance.client = cli
		etcdClientInstance.urls = urls
		etcdClientInstance.timeout = timeout
	})

	return etcdClientInstance
}

// 将新的键值对添加到etcd中
func (e *etcdClient) put(key string, value string) error {
	etcdCli := e.client
	_, err := etcdCli.Put(context.TODO(), key, value)
	if err != nil {
		return err
	}
	return nil
}
