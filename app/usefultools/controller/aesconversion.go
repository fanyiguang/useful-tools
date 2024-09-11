package controller

import (
	"errors"
	"strings"
	"useful-tools/pkg/crypto"
)

type AesConversion struct {
	Base
	conversion     string
	viewText       string
	conversionType string
	aesKey         string
	aesIv          string
	data           string
}

func (a *AesConversion) Data() string {
	return a.data
}

func (a *AesConversion) SetData(data string) {
	a.data = data
}

func (a *AesConversion) AesIv() string {
	return a.aesIv
}

func (a *AesConversion) SetAesIv(aesIv string) {
	a.aesIv = aesIv
}

func (a *AesConversion) AesKey() string {
	return a.aesKey
}

func (a *AesConversion) SetAesKey(aesKey string) {
	a.aesKey = aesKey
}

func (a *AesConversion) ConversionType() string {
	return a.conversionType
}

func (a *AesConversion) SetConversionType(conversionType string) {
	a.conversionType = conversionType
}

func (a *AesConversion) ViewText() string {
	return a.viewText
}

func (a *AesConversion) SetViewText(viewText string) {
	a.viewText = viewText
}

func NewAesConversion() *AesConversion {
	return &AesConversion{}
}

func (a *AesConversion) ConversionList() []string {
	return []string{"解密", "加密"}
}

func (a *AesConversion) DoConversion(mode, key, iv, content string) (string, error) {
	switch mode {
	case "解密":
		return a.Decode(key, iv, content)
	case "加密":
		return a.Encode(key, iv, content)
	}
	return "", errors.New("转换类型错误")
}

func (a *AesConversion) Encode(key, iv, content string) (string, error) {
	jng, err := a.buildSecretKey(key)
	if err != nil {
		return "", err
	}
	pkq, err := a.buildSecretKey(iv)
	if err != nil {
		return "", err
	}
	aesData, err := crypto.AesPKCS7Encrypt(jng, pkq, []byte(content))
	if err != nil {
		return "", err
	}
	return string(aesData), nil
}

func (a *AesConversion) Decode(key, iv, content string) (string, error) {
	jng, err := a.buildSecretKey(key)
	if err != nil {
		return "", err
	}
	pkq, err := a.buildSecretKey(iv)
	if err != nil {
		return "", err
	}
	aesData, err := crypto.AesPKCS7Decrypt(jng, pkq, []byte(content))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(aesData)), nil
}

func (a *AesConversion) buildSecretKey(key string) (new []byte, err error) {
	keys := strings.Split(key, ",")
	if len(keys) > 1 {
		for _, k := range keys {
			b, err := strToByte(k)
			if err != nil {
				return nil, err
			}
			new = append(new, b...)
		}
		return
	}
	new = []byte(key)
	return
}
