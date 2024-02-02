package crypto

import (
	"testing"
)

var (
	Key = []byte{0x55, 0x4a, 0xba, 0xe4, 0xcc, 0x66, 0x2e, 0x2c, 0xb0, 0x73, 0x46, 0xdc, 0xf3, 0x9f, 0xe2, 0x3e}
	IV  = []byte{0xcc, 0x4b, 0x20, 0x66, 0xcf, 0xef, 0x89, 0xf2, 0x47, 0x5d, 0xe1, 0xd4, 0xda, 0x4b, 0x29, 0xc7}
)

func TestAESEncrypt(t *testing.T) {
	//key, iv := "g0MW7KcyAX6DaEBT", "7KcyAX6DaEBT5fsP"
	data := []byte("2024/01/25 22:23:03 [Info] [3850030891] proxy/shadowsocks: tunneling request to tcp:p214-fmf.icloud.com.cn:443 via TCP:38.143.18.57:37827")
	decrypt, err := AESEncrypt(Key, IV, data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(decrypt))
}

func TestAESDecrypt(t *testing.T) {
	//key, iv := "g0MW7KcyAX6DaEBT", "7KcyAX6DaEBT5fsP"
	data := []byte("EYTivX7KmN853IYTbqN9EzrDJG+VxECESQv06VfV0B70YccJ7qXlFI8fIMnfLbUO28HDhm4uNL3YjgwxgNXzd0P9Ff74f7WeLbuQCf9A7kykHfFYXZ3r6iYSYPBoIE398+gpsPB23fiTkE+sWOeiX5VP5JqlzAqaJJj9/1qGsTZhcUogP8VWLUaABveBhKPU")
	decrypt, err := AESDecrypt(Key, IV, data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(decrypt))
}
