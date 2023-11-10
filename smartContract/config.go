package smartcontract

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// oracleConfig 定义和oracle相关的配置信息
type oracleConfig struct {
	EtcdUrls           []string      `yaml:"etcd-urls"`
	EtcdConnectTimeout time.Duration `yaml:"etcd-connect-timeout"`
	EthWsUrl           string        `yaml:"eth-ws-url"`
	EthHttpUrl         string        `yaml:"eth-http-url"`
	// oracle请求合约地址，预言机需要监听该地址
	RequestContractAddr string `yaml:"request-contract-addr"`
	// oracle响应合约地址，worker获取数据后写入该地址
	ResponseContractAddr string `yaml:"response-contract-addr"`
	// 写入oracle智能合约的私钥
	OraclePrivateKey    string `yaml:"oracle-private-key"`
	OracleAccountPasswd string `yaml:"oracle-account-passwd"`
}

// 从配置文件中加载配置数据
func (o *oracleConfig) load(path string) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(bytes, o)
	if err != nil {
		return err
	}
	return nil
}

type configLoader interface {
	load(path string) error
}
