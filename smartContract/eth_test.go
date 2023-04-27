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

// 测试eth客户端的连接
func TestGetEthClientConnection(t *testing.T) {
	ethClient := getEthClientInstance("ws://192.168.31.229:8546", 10)
	if ethClient == nil {
		log.Fatal("客户端连接失败")
	}
	fmt.Println("客户端连接成功: ", ethClient)
}

// 测试请求合约的监听
func TestRegisterRequestContractMonitor(t *testing.T) {
	ethClient := getEthClientInstance("ws://192.168.31.229:8546", 10)
	if ethClient == nil {
		log.Fatal("客户端连接失败")
	}
	// 注册监听事件
	ethClient.registerContractMonitor("0x12dD89a5285Bda38548B3A915757A2DD3CB52992",
		func(err error) {
			fmt.Println("出现了错误", err)
		}, func(logData types.Log) {
			fmt.Println("有数据到来: ", string(logData.Data))
		})
	fmt.Println("监听事件注册成功")
	select {}
}
