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
	ok, err := oracle.WriteData("任务id", "lab405_test")
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
	client := getEtcdClient([]string{"127.0.0.1:2379"}, 100000)
	response, err := client.Get(context.Background(), "0xf2441c45792b049da8907d5f6a5c94896bcac0b16a48bee4b6ea1b0c338dd309")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response)
}
