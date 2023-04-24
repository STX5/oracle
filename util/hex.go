package util

import (
	"fmt"
	"math/big"
)

// decode hex string to binary string
// need to manually check return val's length
func DecodeHex(hID string) (string, error) {
	decimalNumber := new(big.Int)
	decimalNumber, success := decimalNumber.SetString(hID, 16)
	if !success {
		return "", fmt.Errorf("hex string illegal ")
	}
	binaryNumber := decimalNumber.Text(2)
	return binaryNumber, nil
}

// decode binary string to hex string
func DecodeBinary(bID string) (string, error) {
	decimalNumber := new(big.Int)
	decimalNumber, success := decimalNumber.SetString(bID, 2)
	if !success {
		return "", fmt.Errorf("binary string illegal ")
	}
	hexNumber := decimalNumber.Text(16)
	return hexNumber, nil
}

func CheckPrefix(prefix, id string) bool {
	pLen, idLen := len(prefix), len(id)
	if pLen > idLen {
		return false
	}
	for i := 0; i < pLen; i++ {
		if prefix[i] != id[i] {
			return false
		}
	}
	return true
}
