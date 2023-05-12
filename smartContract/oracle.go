package smartcontract

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
	"oracle/smartContract/contract"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
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
	*etcdClient
	*ethClient
	*oracleConfig
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
	// 写入合约的数据
	data string
}

// 定义仅本包内可见的数据
var (
	// oracle对象
	oracle *Oracle
	// 用于实现单例模式的工具对象
	oracleOnce sync.Once
	// 日志对象
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat:        "2006-01-02 15:04:05",
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		},
	}
)

// OracleClient 该变量是暴露给外界使用的对象 主要用于向Oracle合约写入数据
var OracleClient OracleWriter

// 初始化代码
// func init() {
// 	// 初始化oracle对象
// 	oracle = new(Oracle)
// 	err := oracle.initOracle()
// 	if err != nil {
// 		logger.Fatal("初始化oracle对象失败")
// 	}
// 	// 将预言机对象暴露出去
// 	OracleClient = oracle
// }

func (o *Oracle) initOracle() error {
	oracleOnce.Do(func() {
		// 加载oracle的配置文件
		o.oracleConfig = new(oracleConfig)
		if err := oracle.oracleConfig.loadFromYaml("oracle.yaml"); err != nil {
			logger.Fatal("加载oracle的配置文件失败", err)
		}
		// 设置oracle依赖的ethCli对象
		o.ethClient = getEthClientInstance(o.EthUrl, o.ConnectTimeout)
		// 设置oracle依赖的etcdCli对象
		o.etcdClient = getEtcdClientInstance(o.EtcdUrls, o.ConnectTimeout)
		// 开始监听请求智能合约
		logger.Println("设置Oracle请求智能合约监听事件")
		err := o.registerContractMonitor(&OracleRequestContractMonitor{
			contractAddr: o.RequestContractAddr,
		})

		if err != nil {
			logger.Fatal("监听请求智能合约失败")
		}
		logger.Println("oracle对象初始化成功: ", o.oracleConfig)
	})
	return nil
}

// WriteData 将数据写入指定的智能合约
func (o *Oracle) WriteData(data string) (bool, error) {
	logger.Println("向Oracle的ResponseContract写入数据: ", data)
	// 将数据写回智能合约
	err := o.writeDataToContract(&OracleResponseContractInvoker{
		data:         data,
		privateKey:   o.PrivateKey,
		contractAddr: o.ResponseContractAddr,
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
	logger.Println("开始进行事件日志解析")
	abiJson, err := abi.JSON(strings.NewReader(contract.ContractMetaData.ABI))
	event := struct {
		Key   [32]byte
		Value [32]byte
	}{}
	err = abiJson.UnpackIntoInterface(&event, "ItemSet", logData.Data)
	if err != nil {
		logger.Fatal("读取abi文件失败")
	}
	logger.Println("读取事件信息")
	fmt.Println(event)
	fmt.Println(string(event.Key[:]))
	fmt.Println(string(event.Value[:]))
	logger.Println("读取地址信息")
	logger.Println(logData.Address)
	// 获取合约地址
	hash := sha256.New()
	hash.Write(logData.Address.Bytes())
	sum := hash.Sum([]byte(""))
	fmt.Println(sum)
}

// 返回当前要监视的智能合约
func (o *OracleRequestContractMonitor) getMonitorAddr() common.Address {
	return common.HexToAddress(o.contractAddr)
}

// 这里面要写调用写入ResponseContract智能合约的逻辑
func (o *OracleResponseContractInvoker) invoke(opts *bind.TransactOpts) error {
	// todo 创建合约实例
	instance, err := contract.NewContract(o.getContractAddr(), oracle.ethClient.Client)
	if err != nil {
		log.Fatal(err)
	}

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], "foo")
	copy(value[:], "bar")

	num, err := strconv.Atoi(o.getData())
	if err != nil {
		return err
	}
	// todo 调用合约，调用合约的时候，传入opts参数
	tx, err := instance.Store(opts, big.NewInt(int64(num)))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex())

	result, err := instance.Retrieve(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)

	return nil
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
