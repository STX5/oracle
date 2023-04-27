package smartcontract

import (
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"sync"
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
	etcdCli *clientv3.Client
	ethCli  *ethClient
}

// oracle对象
var oracle Oracle

// 用于实现单例模式的工具对象
var once sync.Once

// GetOracleWriter 获取OracleWriter接口对象
func GetOracleWriter(config *OracleConfig) OracleWriter {
	// 返回oracle对象
	return oracle
}

// WriteData 将数据写入指定的智能合约
func (o Oracle) WriteData(data string) (bool, error) {
	return false, nil
}

type TestWriter struct{}

func (TestWriter) WriteData(data string) (bool, error) {
	fmt.Println(data)
	return true, nil
}
