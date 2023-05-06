package smartcontract

import (
	"fmt"
	"log"
	"testing"
)

// 外界应该这样子使用该客户端
func TestOracleClient(t *testing.T) {
	ok, err := OracleClient.WriteData("1")
	if err != nil {
		log.Fatal("写入错误", err)
	}
	if ok {
		fmt.Println("写入成功")
	} else {
		fmt.Println("写入失败")
	}
	select {}
}

func TestGetEthClient(t *testing.T) {
	ethCli := getEthClientInstance("ws://192.168.31.229:8546", 10)
	fmt.Println(ethCli)
}

func TestGetEtcdClient(t *testing.T) {
	ethCli := getEtcdClientInstance([]string{"192.168.31.229:2379"}, 10)
	fmt.Println(ethCli)
}
