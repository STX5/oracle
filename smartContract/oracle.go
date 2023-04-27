package smartcontract

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
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

// OracleRequestContractMonitor 默认的Oracle请求智能合约监听器
// 监听请求的智能合约
type OracleRequestContractMonitor struct {
	contractAddr string
}

// OracleResponseContractInvoker  默认的Oracle响应智能合约调用者
// 调用响应的智能合约
type OracleResponseContractInvoker struct {
	// 调用合约的时候使用的私钥
	privateKey string
	// 调用的合约的地址
	contractAddr string
	data         string
}

// 定义仅本包内可见的数据
var (
	// 默认的Oracle配置
	defaultOracleConfig *OracleConfig
	// oracle对象
	oracle *Oracle
	// 用于实现单例模式的工具对象
	once sync.Once
)

// OracleClient 该变量是暴露给外界使用的对象
// 主要用于向Oracle合约写入数据
var OracleClient OracleWriter

// 进行初始化
// 初始化的过程中
func init() {
	// 生成Oracle的默认配置
	config, err := GetOracleConfigBuilder().
		SetEtcdUrls([]string{}).
		SetEthUrl("ws://192.168.31.229:8546").
		SetPrivateKey("123").
		SetConnectTimeout(10).
		SetRequestContractAddr("123").
		SetResponseContractAddr("123").
		Build()

	// 如果默认配置出错
	if err != nil {
		log.Fatal("初始化oracle默认配置失败")
	}

	// 初始化默认配置全局变量
	defaultOracleConfig = config

	// 获取OracleWriter对象
	// 获取OracleWriter对象的时候，就已经默认设置了RequestContract的
	// 监听事件
	writer, err := getOracleWriter(defaultOracleConfig)

	// 如果获取OracleWriter对象失败
	if err != nil {
		log.Fatal("创建OracleWriter失败")
	}

	// 初始化OracleClient对象
	// 该对象最终暴露给外界使用
	OracleClient = writer
}

// 获取OracleWriter接口对象
func getOracleWriter(config *OracleConfig) (OracleWriter, error) {
	if config == nil {
		return nil, fmt.Errorf("oracle的配置项不能为空")
	}

	once.Do(func() {
		// 创建oracle对象
		oracle = new(Oracle)
		// 设置oracle的配置项
		oracle.config = config
		// 设置oracle依赖的ethCli对象
		oracle.ethCli = getEthClientInstance(config.ethUrl, config.connectTimeout)
		// 设置oracle依赖的etcdCli对象
		oracle.etcdCli = getEtcdClientInstance(config.etcdUrls, config.connectTimeout)
		// 开始监听请求智能合约
		oracle.ethCli.registerContractMonitor(&OracleRequestContractMonitor{
			contractAddr: config.requestContractAddr,
		})
	})
	// 返回oracle对象
	return oracle, nil
}

// WriteData 将数据写入指定的智能合约
func (o *Oracle) WriteData(data string) (bool, error) {
	// 将数据写回智能合约
	err := o.ethCli.writeDataToContract(&OracleResponseContractInvoker{
		data:         data,
		privateKey:   o.config.privateKey,
		contractAddr: o.config.responseContractAddr,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

// 当监视的过程中出现了错误的时候进行处理的逻辑
// 默认的逻辑是打印
func (o *OracleRequestContractMonitor) handleError(err error) {
	fmt.Println(err)
}

// 当监听时间到来的时候，需要通过该方法，将解析后的logData数据写入到etcd中
func (o *OracleRequestContractMonitor) handleLogData(logData types.Log) {
	// TODO 解析logData中的数据
	// TODO 将解析后的数据写入etcd中
	// TODO 生成数据的key
	// TODO 生成数据的value
	// 获取etcdCli对象
	etcdCli := oracle.etcdCli
	key := ""
	value := ""
	err := etcdCli.put(key, value)
	if err != nil {
		log.Fatal("写入etcd失败")
	}
	panic("这里待实现，非常重要的逻辑，监听的智能合约事件被触发后应该做什么")
}

// 返回当前要监视的智能合约
func (o *OracleRequestContractMonitor) getMonitorAddr() common.Address {
	return common.HexToAddress(o.contractAddr)
}

// 这里面要写调用写入ResponseContract智能合约的逻辑
func (o *OracleResponseContractInvoker) invoke(opts *bind.TransactOpts) error {
	panic("这里待实现，非常重要的逻辑，当用户通过OracleClient调用WriteData的时候，怎么将数据写入智能合约")
}

// 返回当前调用者的私钥
func (o *OracleResponseContractInvoker) getPrivateKey() string {
	return o.privateKey
}

// 返回调用的智能合约的地址
func (o *OracleResponseContractInvoker) getContractAddr() common.Address {
	return common.HexToAddress(o.contractAddr)
}

// 返回写入到ResponseContract的数据
func (o *OracleResponseContractInvoker) getData() string {
	return o.data
}
