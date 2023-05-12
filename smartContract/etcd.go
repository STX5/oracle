// Package smartcontract 和etcd操作相关的
package smartcontract

import (
	"context"
	"go.etcd.io/etcd/client/v3"
	"sync"
	"time"
)

// 定义etcd的客户端
type etcdClient struct {
	// 客户端对象
	*clientv3.Client
	// 连接超时时间
	timeout time.Duration
	urls    []string
}

// etcd的单例客户端对象
var etcdClientInstance *etcdClient

// 实现单例的保证
var etcdOnce sync.Once

// 获取etcd客户端的单例
func getEtcdClientInstance(urls []string, timeout time.Duration) *etcdClient {
	etcdOnce.Do(func() {
		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   urls,
			DialTimeout: timeout * time.Second,
		})

		if err != nil {
			logger.Fatal("加载etcd客户端对象失败")
		}

		etcdClientInstance = new(etcdClient)
		etcdClientInstance.Client = cli
		etcdClientInstance.urls = urls
		etcdClientInstance.timeout = timeout
		logger.Println("加载etcdClient对象成功")
	})

	return etcdClientInstance
}

// 将新的键值对添加到etcd中
func (e *etcdClient) put(key string, value string) error {
	_, err := e.Put(context.TODO(), key, value)
	if err != nil {
		return err
	}
	return nil
}
