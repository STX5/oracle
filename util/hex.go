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
	bPre, err := DecodeHex(prefix)
	if err != nil {
		return false
	}
	bID, err := DecodeHex(id)
	if err != nil {
		return false
	}
	pLen, idLen := len(bPre), len(bID)
	if pLen > idLen {
		return false
	}
	for i := 0; i < pLen; i++ {
		if bPre[i] != bID[i] {
			return false
		}
	}
	return true
}
