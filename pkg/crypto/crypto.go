package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
)

func aesEncrypt(origData []byte, key []byte, iv []byte, paddingFunc func([]byte, int) []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = paddingFunc(origData, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func aesDecrypt(crypted, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	return bytes.TrimRight(origData, "\x00"), nil
}

func aesDecryptPkcs7(crypted, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)

	c := origData[len(origData)-1]
	n := int(c)

	return bytes.TrimRight(origData[:len(origData)-n], "\x00"), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// AESEncrypt AES加密，并使用base64输出
func AESEncrypt(key []byte, iv []byte, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	encryptData, err := aesEncrypt(data, key, iv, PKCS5Padding)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(encryptData)))
	base64.StdEncoding.Encode(buf, encryptData)

	return buf, nil
}

func AesPKCS7Encrypt(key []byte, iv []byte, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	encryptData, err := aesEncrypt(data, key, iv, PKCS7Padding)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, base64.StdEncoding.EncodedLen(len(encryptData)))
	base64.StdEncoding.Encode(buf, encryptData)

	return buf, nil
}

// AESEncryptPKCS7WithoutBase64 AES加密
func AESEncryptPKCS7WithoutBase64(key []byte, iv []byte, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	encryptData, err := aesEncrypt(data, key, iv, PKCS7Padding)
	if err != nil {
		return nil, err
	}

	return encryptData, nil
}

// AESEncryptIgnoreErr AES加密，并使用base64输出，忽略错误
func AESEncryptIgnoreErr(key []byte, iv []byte, data []byte) []byte {
	_data, _ := AESEncrypt(key, iv, data)
	return _data
}

// AESDecrypt 先base64解密， 后AES解密
func AESDecrypt(key, iv []byte, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	buf := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(buf, data)
	if err != nil {
		return nil, err
	}
	return aesDecrypt(buf[:n], key, iv)
}

// AESDecryptPKCS7WithoutBase64 AES解密
func AESDecryptPKCS7WithoutBase64(key, iv []byte, data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	return aesDecryptPkcs7(data, key, iv)
}

// AESDecryptIgnoreErr 先base64解密， 后AES解密， 忽略错误
func AESDecryptIgnoreErr(key, iv []byte, data []byte) []byte {
	_data, _ := AESDecrypt(key, iv, data)
	return _data
}

func MD5(data string) string {
	h := md5.New()
	io.WriteString(h, data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// DefaultEncryptWithMergeIV 默认加密合并iv与数据
// aes + cbc + pkcs7
// 先加密data数据，后合并iv+data，再base64
func DefaultEncryptWithMergeIV(data []byte, key []byte) ([]byte, error) {
	ik := key[:16]
	iv := RandomBytes(16)
	encryptData, err := AESEncryptPKCS7WithoutBase64(ik, iv, data)
	if err != nil {
		return nil, err
	}
	mergeData := make([]byte, len(iv)+len(encryptData))
	copy(mergeData[:16], iv)
	copy(mergeData[16:], encryptData)

	result := make([]byte, base64.StdEncoding.EncodedLen(len(mergeData)))
	base64.StdEncoding.Encode(result, mergeData)

	return result, nil
}

// DefaultDecryptWithMergeIV 默认解密合并iv与数据
// aes + cbc + pkcs7
// 先base64解密，后切分iv与data
func DefaultDecryptWithMergeIV(data []byte, key []byte) ([]byte, error) {
	ik := key[:16]
	encryptData := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	n, err := base64.StdEncoding.Decode(encryptData, data)
	if err != nil {
		return nil, err
	}
	encryptData = encryptData[:n]
	iv := encryptData[:16]
	return AESDecryptPKCS7WithoutBase64(ik, iv, encryptData[16:])
}

func AesEncryptCFB(src io.Reader, dst io.Writer, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	writer := cipher.StreamWriter{
		S: stream,
		W: dst,
	}
	_, err = io.Copy(writer, src)
	return err
}

func AesDecryptCFB(src io.Reader, dst io.Writer, key, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	stream := cipher.NewCFBDecrypter(block, iv)
	reader := cipher.StreamReader{
		S: stream,
		R: src,
	}
	_, err = io.Copy(dst, reader)
	return err
}

func RandomBytes(n int) []byte {
	token := make([]byte, n, n)
	rand.Read(token)
	return token
}
