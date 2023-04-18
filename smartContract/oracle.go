package smartcontract

import "fmt"

// interface OracleWriter defines the methods to interact with smart contract
// eg. WriteData() writes job result into oracle contract
// there might be more methods to be added
type OracleWriter interface {
	WriteData(data string) (bool, error)
}

type TestWriter struct{}

func (TestWriter) WriteData(data string) (bool, error) {
	fmt.Println(data)
	return true, nil
}
