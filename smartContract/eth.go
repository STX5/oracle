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
)

type ETHClient struct {
	*ethclient.Client
}

var client *ETHClient

// GetEthClient 获取eth客户端对象
// 在获取eth对象的时候会做如下事情：
// （1）如果是第一次获取，那么创建eth客户端对象，并注册监听的预言机合约的地址
// （2）如果不是第一次获取，那么直接返回eth客户端对象
func GetEthClient(url string) *ETHClient {
	once.Do(func() {
		cli, err := ethclient.Dial(url)
		if err != nil {
			log.Fatal("eth客户端获取错误", err)
		}
		client = new(ETHClient)
		client.Client = cli
		// 在这里注册监听的智能合约的地址
	})
	return client
}

// RegisterEventListener 注册智能合约监听事件
// address: 智能合约的地址
func (e *ETHClient) RegisterEventListener(address string) error {

	contractAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logsChan := make(chan types.Log)
	sub, err := e.SubscribeFilterLogs(context.Background(), query, logsChan)
	if err != nil {
		return err
	}
	for {
		fmt.Println("循环等待")
		select {
		case err := <-sub.Err():
			fmt.Println(err)
		case vlog := <-logsChan:
			fmt.Println(vlog)
		}
	}
}
