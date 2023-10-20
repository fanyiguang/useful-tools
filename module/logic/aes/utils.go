package aes

import (
	"encoding/hex"
	"strings"
)

func proTemplate() string {
	return ""

}

func strToByte(s string) ([]byte, error) {
	hexStr := strings.TrimPrefix(s, "0x")
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
