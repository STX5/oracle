package smartcontract

import (
	"fmt"
	"time"
)

// OracleConfig 定义和oracle相关的配置信息
type OracleConfig struct {
	// 访问etcd的地址
	etcdUrls []string
	// 访问eth的地址
	ethUrl string
	// 连接的超时时间
	connectTimeout time.Duration
	// oracle请求合约地址，预言机需要监听该地址
	requestContractAddr string
	// oracle响应合约地址，worker获取数据后写入该地址
	responseContractAddr string
	// 写入oracle智能合约的私钥
	privateKey string
}

// OracleConfigBuilder 创建OracleConfig对象的建造者
type OracleConfigBuilder struct {
	// Oracle的配置项目
	config OracleConfig
	// 记录生成配置项过程中发生的错误
	err error
}

func GetOracleConfigBuilder() *OracleConfigBuilder {
	return &OracleConfigBuilder{}
}

// SetEtcdUrls 设置oracle依赖的etcd的地址
func (o *OracleConfigBuilder) SetEtcdUrls(urls []string) *OracleConfigBuilder {
	if len(urls) == 0 {
		o.err = fmt.Errorf("etcd的url列表不能为空")
	}
	o.config.etcdUrls = urls
	return o
}

// SetEthUrl 设置oracle依赖的eth的地址
func (o *OracleConfigBuilder) SetEthUrl(url string) *OracleConfigBuilder {
	if len(url) == 0 {
		o.err = fmt.Errorf("eth的url不能为空")
	}
	o.config.ethUrl = url
	return o
}

// SetConnectTimeout 设置访问eth和etcd的超时时间
func (o *OracleConfigBuilder) SetConnectTimeout(timeout time.Duration) *OracleConfigBuilder {
	if timeout <= 0 {
		o.err = fmt.Errorf("超时时间<=0")
	}
	o.config.connectTimeout = timeout
	return o
}

// SetRequestContractAddr 添加oracle的请求的智能合约
func (o *OracleConfigBuilder) SetRequestContractAddr(addr string) *OracleConfigBuilder {
	if len(addr) == 0 {
		o.err = fmt.Errorf("请求合约的地址不能为空")
	}
	o.config.requestContractAddr = addr
	return o
}

// SetResponseContractAddr 添加oracle响应的智能合约
func (o *OracleConfigBuilder) SetResponseContractAddr(addr string) *OracleConfigBuilder {
	if len(addr) == 0 {
		o.err = fmt.Errorf("响应合约的地址不能为空")
	}
	o.config.responseContractAddr = addr
	return o
}

func (o *OracleConfigBuilder) SetPrivateKey(privateKey string) *OracleConfigBuilder {
	if len(privateKey) == 0 {
		o.err = fmt.Errorf("privateKey不能为空")
	}
	o.config.privateKey = privateKey
	return o
}

// Build 返回创建好的Oracle配置对象
func (o *OracleConfigBuilder) Build() (*OracleConfig, error) {
	return &o.config, o.err
}

// Default 仅在测试的时候使用
func (o *OracleConfig) Default() *OracleConfig {
	return &OracleConfig{
		etcdUrls:             []string{"192.168.31.229:2379"},
		ethUrl:               "ws://192.168.31.229:8546",
		connectTimeout:       10,
		requestContractAddr:  "",
		responseContractAddr: "",
		privateKey:           "",
	}
}
