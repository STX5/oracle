package smartcontract

import (
	"context"
	"fmt"
	"log"
	"testing"
)

// 测试Oracle的使用
func TestOracleWriter(t *testing.T) {
	// 创建oracle对象
	_ = NewOracle()
	select {}
}

func TestWrite(t *testing.T) {
	oracle := NewOracle()
	taskMap.jobMap["0xca5b322c08d498403051304c4c9ec9092bdeaf381b0804ba1e6c6227c55721cc"] = "0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040"
	ok, err := oracle.WriteData("0xca5b322c08d498403051304c4c9ec9092bdeaf381b0804ba1e6c6227c55721cc", "lab405_test")
	if err != nil {
		log.Fatal("写入错误", err)
	}
	if ok {
		fmt.Println("写入成功")
	} else {
		fmt.Println("写入失败")
	}
}

func TestQueryEtcd(t *testing.T) {
	client, err := getEtcdClient([]string{"127.0.0.1:2379"}, 100000)
	response, err := client.Get(context.Background(), "0xca5b322c08d498403051304c4c9ec9092bdeaf381b0804ba1e6c6227c55721cc")
	if err != nil {
		log.Fatal(err)
	}
	kvs := response.Kvs
	fmt.Println(kvs)
}
