package smartcontract

import (
	"fmt"
	"testing"
)

func TestGetOracleWriterFailure(t *testing.T) {
	oracleWriter, err := GetOracleWriter(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(oracleWriter)
}
