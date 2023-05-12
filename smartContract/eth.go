// Package smartcontract 和eth智能合约操作相关的
package smartcontract

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 表示eth客户端结构
type ethClient struct {
	*ethclient.Client
	// eth连接的url
	url string
	// 连接的超时时间
	timeout time.Duration
}

// 定义智能合约invoker对象的接口
type ethContractInvoker interface {
	invoke(opts *bind.TransactOpts) error
	getPrivateKey() string
	getContractAddr() common.Address
	getData() string
}

// 定义eth智能合约监听者接口
type ethContractMonitor interface {
	// 当监听到失败的消息的时候，执行该操作
	handleError(err error)
	// 当监听到具体的logData的时候，执行该操作
	handleLogData(logData types.Log)
	// 得到监听的地址
	getMonitorAddr() common.Address
}

// 用于封装监听的事件数据
type eventData struct {
	// worker请求的地址
	url string
	// 解析模式
	pattern string
	// 当前调用的用户信息
	user string
	// 附加的额外信息
	appending string
}

// eth客户端的单例对象
var ethClientInstance *ethClient

// 实现单例的保证
var ethOnce sync.Once

// 获取eth客户端的单例对象
func getEthClientInstance(url string, timeout time.Duration) *ethClient {
	ethOnce.Do(func() {
		cli, err := ethclient.Dial(url)
		if err != nil {
			logger.Fatal("eth客户端连接失败")
		}
		// 创建单例对象
		ethClientInstance = new(ethClient)
		// 设置eth单例对象的属性
		ethClientInstance.url = url
		ethClientInstance.timeout = timeout
		ethClientInstance.Client = cli
		logger.Println("加载eth客户端对象成功")
	})
	return ethClientInstance
}

// 写入数据到智能合约
func (e *ethClient) writeDataToContract(invoker ethContractInvoker) error {
	privateKey, err := crypto.HexToECDSA(invoker.getPrivateKey())
	if err != nil {
		return err
	}

	// 获得当前私钥对应的公钥
	publicKey := privateKey.Public()
	// 获取公钥的ECDSA形式
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("无法将公钥转换为ECDSA形式")
	}

	// 获取eth的客户端对象
	// 当前智能合约的调用地址
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := e.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	// 获得gas费用
	gasPrice, err := e.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("获取gas费用失败")
	}

	chainID, err := e.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("chainID获取失败")
	}
	transactOps, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("NewKeyedTransactorWithChainID调用失败")
	}
	transactOps.Nonce = big.NewInt(int64(nonce))
	transactOps.Value = big.NewInt(0)     // in wei
	transactOps.GasLimit = uint64(300000) // in units
	transactOps.GasPrice = gasPrice

	logger.Println("生成TransactOps对象: ", transactOps)

	// 真实的调用逻辑，因为智能合约调用取决于
	// 具体的智能合约的abi文件是怎么实现的
	// 因此这里需要传入对应的abi文件的invoker对象
	return invoker.invoke(transactOps)
}

// 注册oracle请求合约的监听器
func (e *ethClient) registerContractMonitor(monitor ethContractMonitor) error {
	// 将用户传入的hex格式的地址，转换为Address对象
	contractAddress := monitor.getMonitorAddr()
	// 创建查询过滤器
	queryFilter := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// 创建日志通道
	logChannel := make(chan types.Log)
	// 订阅智能合约的日志事件
	subscription, err := e.SubscribeFilterLogs(context.Background(), queryFilter, logChannel)
	if err != nil {
		// 如果发生了错误，那么直接返回该错误
		return err
	}

	// 开启一个协程额外的处理事件到来的逻辑
	go func() {
		for {
			logger.Println("等待智能合约事件触发")
			select {
			case err := <-subscription.Err():
				// 如果失败了，那么调用外界传入的失败处理器
				logger.Println("智能合约事件监听错误")
				monitor.handleError(err)
			case logData := <-logChannel:
				logger.Println("收到一个新的智能合约事件")
				monitor.handleLogData(logData)
			}
		}
	}()

	return nil
}
