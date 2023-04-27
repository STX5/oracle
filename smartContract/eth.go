// Package smartcontract 和eth智能合约操作相关的
package smartcontract

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"time"
)

// 表示eth客户端结构
type ethClient struct {
	client  *ethclient.Client
	url     string
	timeout time.Duration
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

// 获取eth客户端的单例对象
func getEthClientInstance(url string, timeout time.Duration) *ethClient {
	once.Do(func() {
		cli, err := ethclient.Dial(url)
		if err != nil {
			log.Fatal("eth客户端连接失败")
		}
		// 创建单例对象
		ethClientInstance = new(ethClient)
		// 设置eth单例对象的属性
		ethClientInstance.url = url
		ethClientInstance.timeout = timeout
		ethClientInstance.client = cli
	})
	return ethClientInstance
}

// 注册oracle请求合约的监听器
func (e *ethClient) registerRequestContractMonitor(addr string,
	failure func(err error), handleLog func(logData types.Log)) {

	if failure == nil {
		// 默认的错误处理策略
		failure = func(err error) {
			fmt.Println("监听Oracle请求合约的过程中发生了错误:", err)
		}
	}

	// 将用户传入的hex格式的地址，转换为Address对象
	contractAddress := common.HexToAddress(addr)
	// 创建查询过滤器
	queryFilter := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	// 创建日志通道
	logChannel := make(chan types.Log)
	// 订阅智能合约的日志事件
	subscription, err := e.client.SubscribeFilterLogs(context.Background(), queryFilter, logChannel)
	if err != nil {
		// 如果发生了错误，那么直接返回该错误
		failure(err)
		return
	}

	// 开启一个协程额外的处理事件到来的逻辑
	go func() {
		for {
			fmt.Println("等待事件触发")
			select {
			case err := <-subscription.Err():
				// 如果失败了，那么调用外界传入的失败处理器
				failure(err)
			case logData := <-logChannel:
				handleLog(logData)
			}
		}
	}()
}
