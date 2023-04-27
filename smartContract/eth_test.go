package smartcontract

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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

type TestMonitor struct{}

func (t *TestMonitor) handleError(err error) {
	fmt.Println(err)
}

func (t *TestMonitor) handleLogData(logData types.Log) {
	fmt.Println(logData.Data)
}

func (t *TestMonitor) getMonitorAddr() string {
	return "0x12dD89a5285Bda38548B3A915757A2DD3CB52992"
}

type TestInvoker struct{}

// 调用智能合约
func (t *TestInvoker) invoke(opts *bind.TransactOpts, toAddress common.Address) error {
	fmt.Println("模拟调用")
	return nil
}

// 获取当前调用者的私钥
func (t *TestInvoker) getPrivateKey() string {
	return "模拟私钥"
}

// 获取当前调用的智能合约的地址
func (t *TestInvoker) getContractAddr() string {
	return "模拟智能合约地址"
}

func (t *TestInvoker) getData() string {
	return "模拟数据"
}

// 测试请求合约的监听
func TestRegisterRequestContractMonitor(t *testing.T) {
	ethClient := getEthClientInstance("ws://192.168.31.229:8546", 10)
	if ethClient == nil {
		log.Fatal("客户端连接失败")
	}
	// 注册监听事件
	ethClient.registerContractMonitor(new(TestMonitor))
	select {}
}

// 测试智能合约调用的对象
func TestContractInvoke(t *testing.T) {
	ethClient := getEthClientInstance("ws://192.168.31.229:8546", 10)
	if ethClient == nil {
		log.Fatal("客户端连接失败")
	}
	err := ethClient.writeDataToContract(new(TestInvoker))
	if err != nil {
		fmt.Println("调用出错")
	}
}
