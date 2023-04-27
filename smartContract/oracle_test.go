package smartcontract

import (
	"fmt"
	"log"
	"testing"
)

func TestGetOracleWriter(t *testing.T) {
	// 初始化oracle的配置项
	config, err := GetOracleConfigBuilder().
		SetEtcdUrls([]string{}).
		SetEthUrl("").
		SetPrivateKey("123").
		SetConnectTimeout(10).
		SetRequestContractAddr("123").
		SetResponseContractAddr("123").
		Build()
	if err != nil {
		log.Fatal("创建配置对象出错")
	}

	oracleWriter, err := getOracleWriter(config)
	if err != nil {
		log.Fatal("创建OracleWriter失败")
	}
	ok, err := oracleWriter.WriteData("写入数据")
	if err != nil {
		log.Fatal("写入错误")
	}
	if ok {
		fmt.Println("写入成功")
	} else {
		log.Fatal("写入错误")
	}
}

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
