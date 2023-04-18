package smartcontract

import (
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"testing"
)

func whenSuccess(log types.Log) {
	data := log.Data
	fmt.Println(string(data))
}

func whenFailure(err error) {
	fmt.Println(err)
}

func TestGetEthClient(t *testing.T) {
	ethClient := GetEthClient("ws://192.168.31.229:8546")
	if ethClient == nil {
		log.Fatal("客户端连接失败")
	}
	fmt.Println("客户端连接成功")
	err := ethClient.RegisterEventListener(
		"0x8920661F546cd2FABc538432b2f821E69A5558a7")
	if err != nil {
		log.Fatal("监听智能合约事件失败", err)
	}
	fmt.Println("监听智能合约成功")
}
