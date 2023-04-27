package smartcontract

import (
	"fmt"
	"testing"
)

func TestConfigRight(t *testing.T) {
	config, err := GetOracleConfigBuilder().
		SetEtcdUrls([]string{"localhost:8848"}).
		SetEthUrl("localhost:8849").
		SetRequestContractAddr("fdasfdsfds").
		SetResponseContractAddr("fdasfdsafd").
		SetPrivateKey("fdasfdsaf").
		Build()
	if err != nil {
		fmt.Println("配置生成错误")
		return
	}
	fmt.Println("生成的配置文件为: ", config)
}

func TestConfigWrong(t *testing.T) {
	config, err := GetOracleConfigBuilder().
		SetEtcdUrls([]string{"localhost:8848"}).
		SetEthUrl("localhost:8849").
		SetRequestContractAddr("").
		SetResponseContractAddr("fdasfdsafd").
		Build()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("生成的配置文件为: ", config)
}
