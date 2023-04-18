package smartcontract

import (
	"fmt"
	"go.etcd.io/etcd/client/v3"
)

// OracleWriter interface OracleWriter defines the methods to interact with smart contract
// e.g. WriteData() writes job result into oracle contract
// there might be more methods to be added
type OracleWriter interface {
	WriteData(data string) (bool, error)
}

// Oracle 预言机的实现
type Oracle struct {
	// etcd客户端
	etcdClient *clientv3.Client
}

// oracle对象
var oracle Oracle

// WriteData 将数据写入指定的智能合约
func (o Oracle) WriteData(data string) (bool, error) {
	return false, nil
}

func GetOracleWriter() OracleWriter {
	once.Do(func() {
		oracle = Oracle{
			// 初始化etcdClient
			etcdClient: GetEtcdClient([]string{"192.168.31.229:2379"}),
		}
		// todo 创建完oracle对象之后，需要注册智能合约的监听事件
		// todo
	})
	// 返回oracle
	return oracle
}

type TestWriter struct{}

func (TestWriter) WriteData(data string) (bool, error) {
	fmt.Println(data)
	return true, nil
}
