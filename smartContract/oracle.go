package smartcontract

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
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
	etcdCli *etcdClient
	ethCli  *ethClient
	config  *OracleConfig
}

// oracle对象
var oracle *Oracle

// 用于实现单例模式的工具对象
var once sync.Once

// GetOracleWriter 获取OracleWriter接口对象
func GetOracleWriter(config *OracleConfig) (OracleWriter, error) {
	if config == nil {
		return nil, fmt.Errorf("oracle的配置项不能为空")
	}

	once.Do(func() {
		// 创建oracle对象
		oracle = new(Oracle)
		oracle.config = config
		oracle.ethCli = getEthClientInstance(config.ethUrl, config.connectTimeout)
		oracle.etcdCli = getEtcdClientInstance(config.etcdUrls, config.connectTimeout)
		// 开始监听请求智能合约
		oracle.monitorRequestContract()
	})
	// 返回oracle对象
	return oracle, nil
}

// WriteData 将数据写入指定的智能合约
func (o *Oracle) WriteData(data string) (bool, error) {
	// 将数据写回智能合约
	err := o.ethCli.writeDataToContract(data)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *Oracle) name() {

}

// 监听智能合约监听事件
func (o *Oracle) monitorRequestContract() {

	// 声明处理失败的函数
	handleFailure := func(err error) {
		// 如果失败了，这里默认是打印错误
		fmt.Println(err)
	}

	// 声明处理监听事件数据的函数
	handleLogData := func(logData types.Log) {
		// 这里需要对logData进行解析
		// 首先需要有abi
	}

	// 注册请求合约的监听事件
	o.ethCli.registerContractMonitor(o.config.requestContractAddr,
		handleFailure, handleLogData)
}
