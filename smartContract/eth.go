// Package smartcontract 和eth智能合约操作相关的
package smartcontract

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"sync"
	"time"
)

// 表示eth客户端结构
type ethClient struct {
	*ethclient.Client
	// eth连接的url
	url string
	// 连接的超时时间
	timeout time.Duration
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
		logger.Println("加载ethClient对象成功")
	})
	return ethClientInstance
}
