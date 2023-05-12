package smartcontract

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"log"
	"math/big"
	"oracle/smartContract/contract/request"
	"oracle/smartContract/contract/response"
	"os"
	"strings"
	"sync"
)

// OracleWriter interface OracleWriter defines the methods to interact with smart contract
// e.g. WriteData() writes job result into oracle contract
// there might be more methods to be added
type OracleWriter interface {
	// WriteData oracle被智能合约事件触发后，会向etcd中写入如下格式的数据：
	// {taskHash: xxx, taskFrom: xxx, taskInfo: xxx}
	// taskHash: 该值是oracle根据事件数据计算出来的etcd的任务key，这里选择将其作为值的一部分重复写入，方便worker需要的时候
	// 			 取用
	// taskFrom: 该值是当前任务发起人的eth账户公钥，后续worker根据任务信息取完数据后，该值需要被worker传回，即WriteData的to参数
	// taskInfo: 该值是任务的描述，其形式为{url:xxx, pattern:xxx}
	WriteData(to string, data string) (bool, error)
}

// Oracle 预言机的实现
type Oracle struct {
	// etcd客户端
	*etcdClient
	*ethClient
	*oracleConfig
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

// NewOracle 初始化预言机对象
func NewOracle() OracleWriter {
	oracleOnce.Do(func() {
		// 加载oracle的配置文件
		oracle = new(Oracle)
		oracle.oracleConfig = new(oracleConfig)
		if err := oracle.oracleConfig.loadFromYaml("oracle.yaml"); err != nil {
			logger.Fatal("加载Oracle的配置文件失败", err)
		} else {
			logger.Println("加载Oracle配置文件oracle.yaml成功")
		}

		// 设置oracle依赖的ethCli对象
		oracle.ethClient = getEthClientInstance(oracle.EthUrl, oracle.ConnectTimeout)
		// 设置oracle依赖的etcdCli对象
		oracle.etcdClient = getEtcdClientInstance(oracle.EtcdUrls, oracle.ConnectTimeout)

		// 开始监听请求智能合约
		err := oracle.registerOracleRequestContractMonitor()
		if err != nil {
			logger.Fatal("监听OracleRequestContract失败")
		} else {
			logger.Println("开始监听OracleRequestContract合约事件")
		}
		// 处理监听事件

		logger.Println("oracle对象初始化成功: ", oracle.oracleConfig)
	})
	return oracle
}

// WriteData 将数据写入指定的智能合约
func (o *Oracle) WriteData(to string, data string) (bool, error) {
	logger.Println("向Oracle的ResponseContract写入数据: ", data)
	// 获取私钥，该私钥是oracle的私钥
	privateKey, err := crypto.HexToECDSA(o.PrivateKey)
	if err != nil {
		return false, err
	}

	// 获得当前私钥对应的公钥
	publicKey := privateKey.Public()
	// 获取公钥的ECDSA形式
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return false, fmt.Errorf("无法将公钥转换为ECDSA形式")
	}

	// 获取eth的客户端对象
	// 当前智能合约的调用地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := o.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return false, err
	}

	// 获得gas费用
	gas, err := o.SuggestGasPrice(context.Background())
	if err != nil {
		return false, fmt.Errorf("获取gas费用失败")
	}

	chainID, err := o.ChainID(context.Background())
	if err != nil {
		return false, fmt.Errorf("chainID获取失败")
	}
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return false, fmt.Errorf("NewKeyedTransactorWithChainID调用失败")
	}
	transactOpts.Nonce = big.NewInt(int64(nonce))
	transactOpts.Value = big.NewInt(0)     // in wei
	transactOpts.GasLimit = uint64(300000) // in units
	transactOpts.GasPrice = gas

	oracleResponseContract, err := response.NewResponse(common.HexToAddress(o.ResponseContractAddr),
		oracle.ethClient.Client)
	if err != nil {
		log.Fatal(err)
	}

	_, err = oracleResponseContract.SetValue(transactOpts, common.HexToAddress(to), data)
	if err != nil {
		return false, err
	}
	logger.Println("写入数据成功")
	return true, nil
}

// 注册oracle请求合约的监听事件
func (o *Oracle) registerOracleRequestContractMonitor() error {
	// 将用户传入的hex格式的地址，转换为Address对象
	// 创建查询过滤器
	queryFilter := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(o.RequestContractAddr)},
	}

	// 创建日志通道
	channel := make(chan types.Log)
	// 订阅智能合约的日志事件
	subscription, err := o.SubscribeFilterLogs(context.Background(), queryFilter, channel)
	if err != nil {
		// 如果发生了错误，那么直接返回该错误
		return err
	}

	// 定义事件处理函数
	handleEventFunc := func(data types.Log) error {
		logger.Println("开始进行事件日志解析")
		abiJson, err := abi.JSON(strings.NewReader(request.RequestMetaData.ABI))

		eventInfo := struct {
			// 表示当前事件的触发人
			From common.Address `json:"from"`
			// 当前事件的值
			Value string `json:"value"`
		}{}

		err = abiJson.UnpackIntoInterface(&eventInfo, "RequestEvent", data.Data)
		if err != nil {
			logger.Fatal("解析事件数据失败")
		}
		logger.Println("sender: ", eventInfo.From)
		logger.Println("taskId: ", eventInfo.Value)
		logger.Println("BlockNumber: ", data.BlockNumber)
		// 根据From和BlockNumber计算Hash
		blockNumber := fmt.Sprintf("%d", data.BlockNumber)
		// 计算hash
		hash := crypto.Keccak256Hash([]byte(eventInfo.From.Hex() + blockNumber))
		// 将该hash值和value写入etcd
		logger.Printf("将{%s,%s}写入etcd", hash.Hex(), eventInfo.Value)

		workerData := struct {
			TaskHash string `json:"taskHash"`
			TaskFrom string `json:"taskFrom"`
			TaskInfo string `json:"taskInfo"`
		}{}

		workerData.TaskHash = hash.Hex()
		workerData.TaskFrom = eventInfo.From.Hex()
		workerData.TaskInfo = eventInfo.Value

		bytes, err := json.Marshal(workerData)
		if err != nil {
			return err
		}

		logger.Println("生成任务", string(bytes))
		return o.put(hash.Hex(), string(bytes))
	}

	go func() {
		for {
			logger.Println("等待智能合约事件触发")
			select {
			case err := <-subscription.Err():
				// 如果失败了，那么调用外界传入的失败处理器
				logger.Println("智能合约事件监听错误", err)
			case data := <-channel:
				logger.Println("收到一个新的智能合约事件", data)
				// 处理事件
				if err = handleEventFunc(data); err != nil {
					logger.Println("处理事件发生错误")
				} else {
					logger.Println("事件被成功解析并写入etcd，等到worker处理")
				}
			}
		}
	}()

	// 接收并处理事件
	return nil
}
