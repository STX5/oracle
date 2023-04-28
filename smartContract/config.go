package smartcontract

import (
	"fmt"
	"time"
)

// oracleConfig 定义和oracle相关的配置信息
type oracleConfig struct {
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

// oracleConfigBuilder 创建OracleConfig对象的建造者
type oracleConfigBuilder struct {
	// Oracle的配置项目
	config oracleConfig
	// 记录生成配置项过程中发生的错误
	err error
}

func getOracleConfigBuilder() *oracleConfigBuilder {
	return &oracleConfigBuilder{}
}

// setEtcdUrls 设置oracle依赖的etcd的地址
func (o *oracleConfigBuilder) setEtcdUrls(urls []string) *oracleConfigBuilder {
	if len(urls) == 0 {
		o.err = fmt.Errorf("etcd的url列表不能为空")
	}
	o.config.etcdUrls = urls
	return o
}

// setEthUrl 设置oracle依赖的eth的地址
func (o *oracleConfigBuilder) setEthUrl(url string) *oracleConfigBuilder {
	if len(url) == 0 {
		o.err = fmt.Errorf("eth的url不能为空")
	}
	o.config.ethUrl = url
	return o
}

// setConnectTimeout 设置访问eth和etcd的超时时间
func (o *oracleConfigBuilder) setConnectTimeout(timeout time.Duration) *oracleConfigBuilder {
	if timeout <= 0 {
		o.err = fmt.Errorf("超时时间<=0")
	}
	o.config.connectTimeout = timeout
	return o
}

// setRequestContractAddr 添加oracle的请求的智能合约
func (o *oracleConfigBuilder) setRequestContractAddr(addr string) *oracleConfigBuilder {
	if len(addr) == 0 {
		o.err = fmt.Errorf("请求合约的地址不能为空")
	}
	o.config.requestContractAddr = addr
	return o
}

// setResponseContractAddr 添加oracle响应的智能合约
func (o *oracleConfigBuilder) setResponseContractAddr(addr string) *oracleConfigBuilder {
	if len(addr) == 0 {
		o.err = fmt.Errorf("响应合约的地址不能为空")
	}
	o.config.responseContractAddr = addr
	return o
}

func (o *oracleConfigBuilder) setPrivateKey(privateKey string) *oracleConfigBuilder {
	if len(privateKey) == 0 {
		o.err = fmt.Errorf("privateKey不能为空")
	}
	o.config.privateKey = privateKey
	return o
}

// build 返回创建好的Oracle配置对象
func (o *oracleConfigBuilder) build() (*oracleConfig, error) {
	logger.Println("当前Oracle配置项为: ", o)
	return &o.config, o.err
}
