package smartcontract

import (
	"fmt"
	"log"
	"testing"
)

// 测试Oracle的使用
func TestOracleWriter(t *testing.T) {
	// 创建oracle对象
	oracle := NewOracle()
	// 调用WriteData接口，向ResponseContract写入数据
	ok, err := oracle.WriteData("当前任务发起者的eth公钥地址", "Worker根据任务信息获取到的数据")
	if err != nil {
		log.Fatal("写入错误", err)
	}
	if ok {
		fmt.Println("写入成功")
	} else {
		fmt.Println("写入失败")
	}
}
