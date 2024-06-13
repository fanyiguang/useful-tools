package aes

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"unicode"
)

func TestStrChange(t *testing.T) {
	str := "0x55,0x4a,0xba,0xe4,0xcc,0x66,0x2e,0x2c,0xb0,0x73,0x46,0xdc,0xf3,0x9f,0xe2,0x3e"
	arr := strings.Split(str, ",")
	for _, a := range arr {
		// 移除 "0x" 前缀
		hexStr := strings.TrimPrefix(a, "0x")
		bytes, err := hex.DecodeString(hexStr)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(bytes) // 输出: [65]
		}
	}
}

func TestName(t *testing.T) {
	str := `,"lg_chain_coefficient":"0","lg_disable_auto_switch_line_ips":"","lg_disable_auto_switch_line_regions":"","enable_encrypt_transfer":false,"hosts":{},"local_remote_dns_server":null}`
	result := strings.TrimRightFunc(str, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})
	t.Log(result)
	t.Log([]byte(str))
}
