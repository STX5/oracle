package smartcontract

import (
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
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
	ethClient  *ETHClient
}

// oracle对象
var oracle Oracle

// GetOracleWriter 获取OracleWriter接口对象
func GetOracleWriter() OracleWriter {
	once.Do(func() {
		// todo 创建oracle对象，同时初始化该对象的成员
		oracle = Oracle{
			// 初始化etcdClient
			etcdClient: GetEtcdClient([]string{"192.168.31.229:2379"}),
			// 初始化eth客户端对象
			ethClient: GetEthClient("ws://192.168.31.229:8546"),
		}
		// todo 创建完oracle对象之后，需要注册智能合约的监听事件
		err := oracle.ethClient.RegisterEventListener("0x8920661F546cd2FABc538432b2f821E69A5558a7")
		if err != nil {
			log.Fatal("监听oracle合约失败")
		}
	})
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
