package util

import (
	"fmt"
	"testing"
)

const legalPrefix = "99"
const legalHexID = "99c82bb73505a3c0b453f9fa0e881d6e5a32a0c1"
const illegalHexID = "xxc82bb73505a3c3b453f9fa0e881d6e5a32a0c1"

func TestLegalHexID(t *testing.T) {
	b, err := DecodeHex(legalHexID)
	if err != nil {
		t.Fatal(err)
	}
	h, err := DecodeBinary(b)
	if err != nil {
		t.Fatal(err)
	}
	if h != legalHexID {
		t.Fatal(h)
	}
}

func TestIllegalHexID(t *testing.T) {
	b, err := DecodeHex(illegalHexID)
	if err == nil {
		t.Fatal("No err")
	}
	fmt.Println(b)
}

func TestCheckPrefix(t *testing.T) {
	cool := CheckPrefix(legalPrefix, legalHexID)
	if !cool {
		t.Fatal("Want true, got false")
	}
}
