package smartcontract

import (
	"crypto/sha256"
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

func TestHash(t *testing.T) {
	hash := sha256.New()
	hash.Write([]byte("123456"))
	sum := hash.Sum([]byte(""))
	fmt.Println(sum)
	fmt.Println(len(sum))
}
