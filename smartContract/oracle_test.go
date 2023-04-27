package smartcontract

import (
	"fmt"
	"log"
	"testing"
)

// 外界应该这样子使用该客户端
func TestOracleClient(t *testing.T) {
	ok, err := OracleClient.WriteData("写入的数据")
	if err != nil {
		log.Fatal("写入错误")
	}
	if ok {
		fmt.Println("写入成功")
	} else {
		fmt.Println("写入失败")
	}
}
