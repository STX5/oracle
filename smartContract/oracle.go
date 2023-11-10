package smartcontract

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"oracle/smartContract/contract/request"
	"oracle/smartContract/contract/response"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
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

// oracleClient 预言机的实现
type oracle struct {
	// oracleConfig
	*oracleConfig
	// etcd客户端
	ETCDClient *clientv3.Client
	// eth客户端匿名对象
	ETHClient *ethclient.Client
	// oracleOnce sync.Once
	Logger *logrus.Logger
	// 任务映射
	taskMap       *oracleTaskMap
	alterConfigCh chan oracleConfig
}

// 任务映射
type oracleTaskMap struct {
	// 保证oracleTaskMap更新的原子性
	sync.Mutex
	// jobMap 记录JobId和JobFrom的映射
	jobMap map[string]string
}

// lockAndUnlockResult 结果
type lockAndUnlockResult struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  bool   `json:"result"`
}

func NewOracle() (*oracle, error) {
	logger := &logrus.Logger{
		Out:   os.Stderr,
		Level: logrus.DebugLevel,
		Formatter: &logrus.TextFormatter{
			TimestampFormat:        "2006-01-02 15:04:05",
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		},
	}
	o := &oracle{
		oracleConfig: new(oracleConfig),
		Logger:       logger,
		taskMap:      new(oracleTaskMap),
	}
	// 初始化任务映射记录结构
	o.taskMap.jobMap = make(map[string]string)
	// 加载oracle的配置文件
	// TODO: more methods
	// eg. read config from database/network
	if err := o.oracleConfig.load("oracle.yaml"); err != nil {
		return &oracle{}, fmt.Errorf("error loading oralce config", err)
	}
	if cli, err := buildEthClient(o.EthWsUrl); err != nil {
		return &oracle{}, fmt.Errorf("error building ETH Client", err)
	} else {
		o.ETHClient = cli
	}
	if cli, err := buildEtcdClient(o.EtcdUrls, o.EtcdConnectTimeout); err != nil {
		return &oracle{}, fmt.Errorf("error building ETCD Client", err)
	} else {
		o.ETCDClient = cli
	}

	o.Logger.Info("oracle init success", o.oracleConfig)
	// Println("oracle对象初始化成功: ", oracle.oracleConfig)

	return o, nil
}

func (o *oracle) Run() {
	ctx, cancel := context.WithCancel(context.Background())

	// 开始监听请求智能合约
	go func() {
		err := o.registerOracleRequestContractMonitor(ctx)
		if err != nil {
			o.Logger.Fatal("监听OracleRequestContract失败")
		} else {
			o.Logger.Info("开始监听OracleRequestContract合约事件")
		}
	}()
	// TODO: usa another goroutine to listen for config change
	for newConfig := range o.alterConfigCh {
		// TODO: add a lock and deal with it!!!
		cancel()
		o.oracleConfig = &newConfig
		o.Run()
	}
}

// WriteData 将数据写入指定的智能合约
func (o *oracle) WriteData(jobID string, data string) (bool, error) {
	// 写入数据之前，先尝试解锁账户
	err := o.tryUnlockAccount()
	if err != nil {
		o.Logger.Warnf("尝试解锁账户失败")
		return false, err
	}
	// 写入数据成功后，尝试重新锁定账户
	defer func(o *oracle) {
		err := o.tryLockAccount()
		if err != nil {
			o.Logger.Infoln("尝试重新锁定账户失败")
		}
	}(o)

	// 获取私钥，该私钥是oracle的私钥
	privateKey, err := crypto.HexToECDSA(o.OraclePrivateKey)
	if err != nil {
		return false, err
	}
	err, fromAddress := getPublicKeyAddress(privateKey)
	nonce, err := o.ETHClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return false, err
	}

	// 获得gas费用
	gas, err := o.ETHClient.SuggestGasPrice(context.Background())
	if err != nil {
		return false, fmt.Errorf("获取gas费用失败")
	}

	chainID, err := o.ETHClient.ChainID(context.Background())
	if err != nil {
		return false, fmt.Errorf("chainID获取失败")
	}
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return false, fmt.Errorf("NewKeyedTransactorWithChainID调用失败")
	}
	transactOpts.Nonce = big.NewInt(int64(nonce))
	transactOpts.Value = big.NewInt(0)
	transactOpts.GasLimit = uint64(300000)
	transactOpts.GasPrice = gas

	oracleResponseContract, err := response.NewResponse(common.HexToAddress(o.ResponseContractAddr), o.ETHClient)
	if err != nil {
		log.Fatal(err)
	}

	toAddr, err := o.taskMap.get(jobID)
	if err != nil {
		return false, err
	}

	_, err = oracleResponseContract.SetValue(transactOpts, common.HexToAddress(toAddr), data)
	if err != nil {
		return false, err
	}

	// 如果写入成功，那么需要删除当前jobID和任务发起者的映射
	o.taskMap.remove(jobID)
	o.Logger.WithFields(logrus.Fields{
		"ToAddr": toAddr,
		"Data":   data,
	}).Info("写入智能合约成功")
	return true, nil
}

// 注册oracle请求合约的监听事件
func (o *oracle) registerOracleRequestContractMonitor(ctx context.Context) error {
	// 将用户传入的hex格式的地址，转换为Address对象
	// 创建查询过滤器
	queryFilter := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(o.RequestContractAddr)},
	}

	// 创建日志通道
	channel := make(chan types.Log)
	// 订阅智能合约的日志事件
	subscription, err := o.ETHClient.SubscribeFilterLogs(context.Background(), queryFilter, channel)
	if err != nil {
		// 如果发生了错误，那么直接返回该错误
		return err
	}

	// 定义事件处理函数
	handleEventFunc := func(data types.Log) error {
		o.Logger.Println("开始进行事件日志解析")
		abiJson, err := abi.JSON(strings.NewReader(request.RequestMetaData.ABI))

		eventInfo := struct {
			// 表示当前事件的触发人
			From common.Address `json:"from"`
			// 当前事件的值
			Value string `json:"value"`
		}{}

		err = abiJson.UnpackIntoInterface(&eventInfo, "RequestEvent", data.Data)
		if err != nil {
			o.Logger.Fatal("解析事件数据失败")
		}
		o.Logger.Println("JobFrom: ", eventInfo.From)
		//o.Logger.Println("Pattern: ", eventInfo.Value.Pattern)
		//o.Logger.Println("EthWsUrl: ", eventInfo.Value.Url)
		o.Logger.Println("BlockNumber: ", data.BlockNumber)
		// 根据From和BlockNumber计算Hash
		blockNumber := fmt.Sprintf("%d", data.BlockNumber)
		// 计算hash，生成jobID
		jobID := crypto.Keccak256Hash([]byte(eventInfo.From.Hex() + blockNumber))
		o.Logger.Println("JobID: ", jobID)

		// 将当前jobID和job发起者的地址映射关系存放起来
		if err = o.taskMap.put(jobID.Hex(), eventInfo.From.Hex()); err != nil {
			return err
		}
		// 创建job值，传输给job的值，是符合worker的需要的
		workerData := struct {
			URL     string `json:"url"`
			Pattern string `json:"pattern"`
		}{}

		err = json.Unmarshal([]byte(eventInfo.Value), &workerData)
		if err != nil {
			return err
		}

		// 	序列化jobVal
		workerDataBytes, err := json.Marshal(workerData)
		if err != nil {
			return err
		}

		o.Logger.Println("生成任务{key: ", jobID.Hex(), ", value: ", string(workerDataBytes), "}")
		_, err = o.ETCDClient.Put(context.Background(), jobID.Hex(), string(workerDataBytes))
		return err
	}

	go func() {
		for {
			o.Logger.Println("等待智能合约事件触发")
			select {
			case err := <-subscription.Err():
				// 如果失败了，那么调用外界传入的失败处理器
				o.Logger.Println("智能合约事件监听错误: ", err)
			case data := <-channel:
				o.Logger.Println("收到一个新的智能合约事件", data)
				// 处理事件
				if err = handleEventFunc(data); err != nil {
					o.Logger.Println("处理事件发生错误: ", err)
				} else {
					o.Logger.Println("事件被成功解析并写入etcd，等到worker处理")
				}
			}
		}
	}()

	// 接收并处理事件
	return nil
}

func getPublicKeyAddress(privateKey *ecdsa.PrivateKey) (error, common.Address) {
	// 获得当前私钥对应的公钥
	publicKey := privateKey.Public()
	// 获取公钥的ECDSA形式
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("无法将公钥转换为ECDSA形式"), common.Address{}
	}

	// 获取eth的客户端对象
	// 当前智能合约的调用地址
	publicAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return nil, publicAddress
}

// 获取etcd客户端
func buildEtcdClient(urls []string, timeout time.Duration) (*clientv3.Client, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   urls,
		DialTimeout: timeout * time.Second,
	})

	if err != nil {
		return nil, err
	}
	return cli, nil
}

// 获取eth客户端
func buildEthClient(url string) (*ethclient.Client, error) {
	cli, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// 向taskMap中添加一个新的映射关系
func (o *oracleTaskMap) put(jobID string, jobFrom string) error {
	o.Lock()
	defer o.Unlock()
	_, ok := o.jobMap[jobID]
	if ok {
		// 说明存在了重复的任务
		return fmt.Errorf("出现了重复JobID: %s", jobID)
	} else {
		// 说明没有发现重复任务，那么需要将该任务的id和发起任务的发起人的公钥进行绑定
		o.jobMap[jobID] = jobFrom
		log.Println("记录{JobID: ", jobID, ", TaskFrom: ", jobFrom, "}")
	}
	return nil
}

// 从taskMap中查询jobID对应的job发起者地址
func (o *oracleTaskMap) get(jobID string) (string, error) {
	o.Lock()
	defer o.Unlock()
	v, ok := o.jobMap[jobID]
	if !ok {
		// 如果当前jobID不存在，那么说明是非法的jobID
		return "", fmt.Errorf("不存在的JobID: %s", jobID)
	}
	// 将jobID对应的任务发起者地址返回
	return v, nil
}

// 移除jobID对应的键值对
func (o *oracleTaskMap) remove(jobID string) {
	o.Lock()
	defer o.Unlock()
	delete(o.jobMap, jobID)
}

// 解锁预言机的账户，不然无法进行合约的执行
func (o *oracle) tryUnlockAccount() error {
	// 获取私钥，该私钥是oracle的私钥
	privateKey, err := crypto.HexToECDSA(o.OraclePrivateKey)
	if err != nil {
		return err
	}
	err, publicKeyAddress := getPublicKeyAddress(privateKey)
	if err != nil {
		return err
	}

	// 如果被锁住了，那么执行如下解锁账户的操作
	unlockRequest := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_unlockAccount\",\"params\":[\"%s\", \"%s\", 30],\"id\":1}"
	param := fmt.Sprintf(unlockRequest, publicKeyAddress.Hex(), o.OracleAccountPasswd)
	unlockResultBytes, err := invokeJsonRpc(o.EthHttpUrl, []byte(param))

	result := new(lockAndUnlockResult)
	err = json.Unmarshal(unlockResultBytes, result)
	if err != nil {
		return err
	}

	if result.Result {
		return nil
	}

	return fmt.Errorf("解锁账户失败")
}

// 尝试锁定账户
func (o *oracle) tryLockAccount() error {
	privateKey, err := crypto.HexToECDSA(o.OraclePrivateKey)
	if err != nil {
		return err
	}
	err, publicKeyAddress := getPublicKeyAddress(privateKey)
	if err != nil {
		return err
	}
	lockRequest := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_lockAccount\",\"params\":[\"%s\"],\"id\":1}"
	param := fmt.Sprintf(lockRequest, publicKeyAddress.Hex())
	lockResultBytes, err := invokeJsonRpc(o.EthHttpUrl, []byte(param))
	if err != nil {
		return err
	}

	result := new(lockAndUnlockResult)
	err = json.Unmarshal(lockResultBytes, result)
	if err != nil {
		return err
	}

	if result.Result {
		return nil
	}

	return fmt.Errorf("锁定账户失败")
}

// 生成ethRequest，Marshal 成[]byte ，传入do函数即可操作ethereum 节点
func invokeJsonRpc(url string, param []byte) ([]byte, error) {
	reader := bytes.NewReader(param)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	respBytes, err := io.ReadAll(resp.Body)
	fmt.Println(string(respBytes))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return respBytes, nil
}
