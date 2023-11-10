package smartcontract

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func TestLockAndUnlockAccount(t *testing.T) {
	unlockRequest := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_unlockAccount\",\"params\":[\"0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040\", \"cs237239\", 30],\"id\":1}"
	lockRequest := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_lockAccount\",\"params\":[\"0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040\"],\"id\":1}"
	listWalletRequest := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_listWallets\",\"params\":[],\"id\":1}"
	result, err := invokeJsonRpc("http://127.0.0.1:8545", []byte(unlockRequest))
	if err != nil {
		logger.Fatal("解锁账户失败")
	}
	logger.Println(string(result))

	result, err = invokeJsonRpc("http://127.0.0.1:8545", []byte(lockRequest))
	if err != nil {
		logger.Fatal("锁定账户失败")
	}
	logger.Println(string(result))

	result, err = invokeJsonRpc("http://127.0.0.1:8545", []byte(listWalletRequest))
	if err != nil {
		logger.Fatal("查询钱包数据失败")
	}
	logger.Println(string(result))
}

// 测试Oracle的使用
func TestOracleWriter(t *testing.T) {
	// 创建oracle对象
	_, _ = NewOracle()
	select {}
}

func TestWrite(t *testing.T) {
	oracle, _ := NewOracle()
	oracle.taskMap.jobMap["0xca5b322c08d498403051304c4c9ec9092bdeaf381b0804ba1e6c6227c55721cc"] = "0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040"
	ok, err := oracle.WriteData("0xca5b322c08d498403051304c4c9ec9092bdeaf381b0804ba1e6c6227c55721cc", "测试数据")
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

func TestQueryEtcd(t *testing.T) {
	client, err := buildEtcdClient([]string{"127.0.0.1:2379"}, 100000)
	if err != nil {
		log.Fatalln("query ETCD failed")
	}
	response, err := client.Get(context.Background(), "0xca5b322c08d498403051304c4c9ec9092bdeaf381b0804ba1e6c6227c55721cc")
	if err != nil {
		log.Fatal(err)
	}
	kvs := response.Kvs
	fmt.Println(kvs)
}
