package smartcontract

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	config := new(oracleConfig)
	err := config.load("oracle.yaml")
	if err != nil {
		log.Fatal(err)
	}
	bytes, _ := json.Marshal(config)
	fmt.Println(string(bytes))
}
