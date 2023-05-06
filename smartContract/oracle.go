package smartcontract

import (
	"crypto/sha256"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/sirupsen/logrus"
	"log"
	"math/big"
	"oracle/smartContract/contract"
	"os"
	"strconv"
	"strings"
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
	data         string
}

// 定义仅本包内可见的数据
var (
	// 默认的Oracle配置
	defaultOracleConfig *oracleConfig
	// oracle对象
	oracle *Oracle
	// 用于实现单例模式的工具对象
	oracleOnce sync.Once
)

// 日志工具
var logger *logrus.Logger

// OracleClient 该变量是暴露给外界使用的对象 主要用于向Oracle合约写入数据
var OracleClient OracleWriter

// 初始化代码
func init() {
	// 初始化logger对象，该对象负责整个smartContract包内的日志打印工作
	logger = &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat:        "2006-01-02 15:04:05",
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		},
	}

	// 生成Oracle的默认配置
	config, err := getOracleConfigBuilder().
		setEtcdUrls([]string{"192.168.31.229:2379"}).
		setEthUrl("ws://192.168.31.229:8546").
		setPrivateKey("9d8d022bc819c303592b7f384377769f9877b46d087c371c35f5e5049c1e28cf").
		setConnectTimeout(10).
		setRequestContractAddr("0x06f4f74252bf1E82651Dc8b141b26D443a3c7e11").
		setResponseContractAddr("0x06f4f74252bf1E82651Dc8b141b26D443a3c7e11").
		build()
	if err != nil {
		logger.Fatal("初始化oracle默认配置失败")
	}
	defaultOracleConfig = config

	// 根据默认配置获取OracleWriter对象
	// 获取OracleWriter对象的时候，就已经默认设置了RequestContract的监听事件
	writer, err := getOracleWriter(defaultOracleConfig)
	if err != nil {
		logger.Fatal("创建OracleWriter失败")
	}

	// 初始化OracleClient对象，该对象是整个smartContract包中提供给外界的
	// 接口对象之一，另外一个是OracleConfig
	OracleClient = writer
}

// 获取OracleWriter接口对象
func getOracleWriter(config *oracleConfig) (OracleWriter, error) {
	if config == nil {
		return nil, fmt.Errorf("oracle的配置项不能为空")
	}

	oracleOnce.Do(func() {
		logger.Println("创建Oracle对象")
		// 创建oracle对象
		oracle = new(Oracle)
		// 设置oracle的配置项
		oracle.oracleConfig = config
		// 设置oracle依赖的ethCli对象
		oracle.ethClient = getEthClientInstance(config.ethUrl, config.connectTimeout)
		// 设置oracle依赖的etcdCli对象
		oracle.etcdClient = getEtcdClientInstance(config.etcdUrls, config.connectTimeout)
		// 开始监听请求智能合约
		logger.Println("设置Oracle请求智能合约监听事件")
		err := oracle.registerContractMonitor(&OracleRequestContractMonitor{
			contractAddr: config.requestContractAddr,
		})

		if err != nil {
			logger.Fatal("监听请求智能合约失败")
		}
	})
	// 返回oracle对象
	return oracle, nil
}

// WriteData 将数据写入指定的智能合约
func (o *Oracle) WriteData(data string) (bool, error) {
	logger.Println("向Oracle的ResponseContract写入数据: ", data)
	// 将数据写回智能合约
	err := o.writeDataToContract(&OracleResponseContractInvoker{
		data:         data,
		privateKey:   o.privateKey,
		contractAddr: o.responseContractAddr,
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
