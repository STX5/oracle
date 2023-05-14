package smartcontract

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"log"
	"testing"
)

import (
	"bytes"
	"net/http"
)

type ethRequest struct {
	general
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type general struct {
	JsonRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

// 生成ethRequest，Marshal 成[]byte ，传入do函数即可操作ethereum 节点
func do(jsonParam []byte) ([]byte, error) {
	reader := bytes.NewReader(jsonParam)
	url := "http://127.0.0.1:8545"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return []byte(""), err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return []byte(""), err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}
	return respBytes, nil
}

func TestJsonRpc(t *testing.T) {
	//unlockRequest := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_unlockAccount\",\"params\":[\"0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040\", \"cs237239\", 30],\"id\":1}"
	//unlockRequest := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_unlockAccount\",\"params\":[\"0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040\", \"cs237239\", 30]}"
	//lockRequest := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_lockAccount\",\"params\":[\"0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040\"],\"id\":1}"
	listAccounts := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_listWallets\",\"params\":[],\"id\":1}"
	//isLockFormat := "{\"jsonrpc\":\"2.0\",\"method\":\"personal_listWallets\",\"params\":[]}"
	i, err := do([]byte(listAccounts))
	if err != nil {
		fmt.Println(err)
	}
	result := new(listWalletsResult)
	bytes, err := do([]byte(listAccounts))
	err = json.Unmarshal(bytes, result)
	fmt.Println(err)
	fmt.Println(result)
	privateKey, err := crypto.HexToECDSA("9e082d9aa6240a7af133616ac9f740aa67ba72b0372a366ac9890bb6e9740321")
	if err != nil {
		fmt.Println(err)
	}
	status := result.Result[0].Status
	fmt.Println(status)
	// 获得当前私钥对应的公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Println(address == common.HexToAddress("0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040"))
	fmt.Println(string(i))
}

// 测试Oracle的使用
func TestOracleWriter(t *testing.T) {
	// 创建oracle对象
	_ = NewOracle()
	select {}
}

func TestWrite(t *testing.T) {
	oracle := NewOracle()
	taskMap.jobMap["0xca5b322c08d498403051304c4c9ec9092bdeaf381b0804ba1e6c6227c55721cc"] = "0x8Ff1EBd1639dF7ca2FCF67eAcDC7488d567ec040"
	ok, err := oracle.WriteData("0xca5b322c08d498403051304c4c9ec9092bdeaf381b0804ba1e6c6227c55721cc", "hello,oracle")
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
