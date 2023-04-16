package smartcontract

// interface OracleWriter defines the methods to interact with smart contract
// eg. WriteData() writes job result into oracle contract
// there might be more methods to be added
type OracleWriter interface {
	WriteData(data string) (bool, error)
}
