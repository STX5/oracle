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
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"math/big"
	"oracle/smartContract/contract/request"
	"oracle/smartContract/contract/response"
	"os"
	"strings"
	"sync"
	"time"
)

// OracleWriter interface OracleWriter defines the methods to interact with smart contract
// e.g. WriteData() writes job result into oracle contract
// there might be more methods to be added
type OracleWriter interface {
	// WriteData 将数据写入oracle响应合约
	// jobID: 本次任务的id，也就是worker接收任务时对应的etcd的key
	// data: 写入到oracle响应合约的数据
	WriteData(jobID string, data string) (bool, error)
}

// 定义eth客户端
type ethClient struct {
	*ethclient.Client
}

// 定义etcd客户端
type etcdClient struct {
	*clientv3.Client
}

// oracleClient 预言机的实现
type oracleClient struct {
	// etcd客户端
	*etcdClient
	// eth客户端匿名对象
	*ethClient
	// oracleConfig
	*oracleConfig
}

// JobMap 记录JobId和JobFrom的映射
type JobMap map[string]string

// 任务映射
type oracleTaskMap struct {
	sync.Mutex
	JobMap
}

// 定义仅本包内可见的数据
var (
	// oracle对象
	oracle *oracleClient
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
	// 任务映射
	taskMap *oracleTaskMap
)

// NewOracle 初始化预言机对象
func NewOracle() OracleWriter {
	oracleOnce.Do(func() {
		// 加载oracle的配置文件
		oracle = new(oracleClient)
		oracle.oracleConfig = new(oracleConfig)
		if err := oracle.oracleConfig.loadFromYaml("oracle.yaml"); err != nil {
			logger.Fatal("加载Oracle的配置文件失败", err)
		} else {
			logger.Println("加载Oracle配置文件oracle.yaml成功")
		}

		// 设置oracle依赖的ethCli对象
		oracle.ethClient = getEthClient(oracle.EthUrl)
		// 设置oracle依赖的etcdCli对象
		oracle.etcdClient = getEtcdClient(oracle.EtcdUrls, oracle.EtcdConnectTimeout)

		// 初始化任务映射记录结构
		taskMap = new(oracleTaskMap)

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
func (o *oracleClient) WriteData(jobID string, data string) (bool, error) {
	logger.Println("向Oracle的ResponseContract写入数据: ", data)
	// 获取私钥，该私钥是oracle的私钥
	privateKey, err := crypto.HexToECDSA(o.OraclePrivateKey)
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

	// 对taskMap进行加锁
	taskMap.Lock()
	var toAddr common.Address
	v, ok := taskMap.JobMap[jobID]
	if !ok {
		// 如果当前jobID不存在，那么说明是非法的jobID
		return false, fmt.Errorf("不存在的JobID: %s", jobID)
	}
	// 否则说明jobID是存在的，那么取出当前jobID对应的toAddr
	toAddr = common.HexToAddress(v)
	defer taskMap.Unlock()

	// 将worker传入的数据写入智能合约
	_, err = oracleResponseContract.SetValue(transactOpts, toAddr, data)
	if err != nil {
		return false, err
	}
	logger.Println("数据: {ToAddr: ", toAddr, ", Data: ", data, "}写入智能合约成功")
	return true, nil
}

// 注册oracle请求合约的监听事件
func (o *oracleClient) registerOracleRequestContractMonitor() error {
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
			Value struct {
				Pattern string `json:"pattern"`
				Url     string `json:"url"`
			} `json:"value"`
		}{}

		err = abiJson.UnpackIntoInterface(&eventInfo, "RequestEvent", data.Data)
		if err != nil {
			logger.Fatal("解析事件数据失败")
		}
		logger.Println("JobFrom: ", eventInfo.From)
		logger.Println("Pattern: ", eventInfo.Value.Pattern)
		logger.Println("EthUrl: ", eventInfo.Value.Url)
		logger.Println("BlockNumber: ", data.BlockNumber)
		// 根据From和BlockNumber计算Hash
		blockNumber := fmt.Sprintf("%d", data.BlockNumber)
		// 计算hash，生成jobID
		jobID := crypto.Keccak256Hash([]byte(eventInfo.From.Hex() + blockNumber))
		logger.Println("JobID: ", jobID)

		// 在这里加锁保证map操作的原子性
		taskMap.Lock()
		_, ok := taskMap.JobMap[jobID.Hex()]
		if !ok {
			// 说明存在了重复的任务
			return fmt.Errorf("出现了重复JobID: %s", jobID.Hex())
		} else {
			// 说明没有发现重复任务，那么需要将该任务的id和发起任务的发起人的公钥进行绑定
			taskMap.JobMap[jobID.Hex()] = eventInfo.From.Hex()
			logger.Println("记录{JobID: ", jobID, ", TaskFrom: ", eventInfo.From.Hex(), "}")
		}
		// 释放锁
		defer taskMap.Unlock()
		// 创建job值，传输给job的值，是符合worker的需要的
		workerData := struct {
			URL     string
			Pattern string
		}{URL: eventInfo.Value.Url, Pattern: eventInfo.Value.Pattern}

		// 	序列化jobVal
		workerDataBytes, err := json.Marshal(workerData)
		if err != nil {
			return err
		}

		logger.Println("生成任务{key: ", jobID.Hex(), ", value: ", string(workerDataBytes), "}")
		_, err = o.Put(context.Background(), jobID.Hex(), string(workerDataBytes))
		return err
	}

	go func() {
		for {
			logger.Println("等待智能合约事件触发")
			select {
			case err := <-subscription.Err():
				// 如果失败了，那么调用外界传入的失败处理器
				logger.Println("智能合约事件监听错误: ", err)
			case data := <-channel:
				logger.Println("收到一个新的智能合约事件", data)
				// 处理事件
				if err = handleEventFunc(data); err != nil {
					logger.Println("处理事件发生错误: ", err)
				} else {
					logger.Println("事件被成功解析并写入etcd，等到worker处理")
				}
			}
		}
	}()

	// 接收并处理事件
	return nil
}

// 获取etcd客户端的单例
func getEtcdClient(urls []string, timeout time.Duration) *etcdClient {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   urls,
		DialTimeout: timeout * time.Second,
	})

	if err != nil {
		logger.Fatal("加载etcd客户端对象失败")
	}

	etcdCli := new(etcdClient)
	etcdCli.Client = cli
	return etcdCli
}

// 获取eth客户端的单例对象
func getEthClient(url string) *ethClient {
	cli, err := ethclient.Dial(url)
	if err != nil {
		logger.Fatal("eth客户端连接失败")
	}
	// 创建单例对象
	ethCli := new(ethClient)
	// 设置eth单例对象的属性
	ethCli.Client = cli
	return ethCli
}
