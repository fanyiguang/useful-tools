package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

type Code struct {
	key, iv []byte
}

func NewAesCode(key, iv []byte) *Code {
	obj := &Code{
		key, iv,
	}

	return obj
}

func (a *Code) pKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (a *Code) pKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// aesEncrypt 加密函数
func (a *Code) aesEncrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	plaintext = a.pKCS7Padding(plaintext, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, a.iv)
	crypted := make([]byte, len(plaintext))
	blockMode.CryptBlocks(crypted, plaintext)
	return crypted, nil
}

// aesDecrypt 解密函数
func (a *Code) aesDecrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, a.iv[:blockSize])
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	origData = a.pKCS7UnPadding(origData)
	return origData, nil
}

func (a *Code) AesDecrypt(plaintext []byte) ([]byte, error) {
	// key, _ := hex.DecodeString("6368616e676520746869732070617373")

	if a.iv == nil {
		c := make([]byte, aes.BlockSize+len(plaintext))
		a.iv = c[:aes.BlockSize]
	}

	// 加密
	return a.aesDecrypt(plaintext)
}

func (a *Code) AesEncrypt(plaintext []byte) ([]byte, error) {
	// key, _ := hex.DecodeString("6368616e676520746869732070617373")

	if a.iv == nil {
		c := make([]byte, aes.BlockSize+len(plaintext))
		a.iv = c[:aes.BlockSize]
	}

	// 加密
	return a.aesEncrypt(plaintext)
}

func AesDecodeWithKey(key, iv, msg []byte) (decrypt []byte, err error) {
	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("AesDecodeWithKey panic: %v", perr)
		}
	}()
	decrypt, err = NewAesCode(key, iv).AesDecrypt(msg)
	return
}

func AesEncodeWithKey(key, iv, msg []byte) (encode []byte, err error) {
	defer func() {
		if perr := recover(); perr != nil {
			err = fmt.Errorf("AesEncodeWithKey panic: %v", perr)
		}
	}()
	encode, err = NewAesCode(key, iv).AesEncrypt(msg)
	return
}
