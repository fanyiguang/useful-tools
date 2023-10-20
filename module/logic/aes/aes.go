package aes

import (
	"encoding/base64"
	"strings"
	"useful-tools/module/logic/base"
	"useful-tools/pkg/aes"
)

type Aes struct {
	base.Base
	key          string
	iv           string
	inputContent string
	viewContent  string
	convertType  string
}

func New() *Aes {
	p := new(Aes)
	p.SetProTemplate(proTemplate())
	return p
}

func (a *Aes) ConvertType() string {
	return a.convertType
}

func (a *Aes) SetConvertType(convertType string) {
	a.convertType = convertType
}

func (a *Aes) Key() string {
	return a.key
}

func (a *Aes) SetKey(key string) {
	a.key = key
}

func (a *Aes) Iv() string {
	return a.iv
}

func (a *Aes) SetIv(iv string) {
	a.iv = iv
}

func (a *Aes) InputContent() string {
	return a.inputContent
}

func (a *Aes) SetInputContent(inputContent string) {
	a.inputContent = inputContent
}

func (a *Aes) ViewContent() string {
	return a.viewContent
}

func (a *Aes) SetViewContent(viewContent string) {
	a.viewContent = viewContent
}

func (a *Aes) Encode(key, iv, content string) (string, error) {
	_key, err := a.buildSecretKey(key)
	if err != nil {
		return "", err
	}
	_iv, err := a.buildSecretKey(iv)
	if err != nil {
		return "", err
	}

	aesData, err := aes.AesEncodeWithKey(_key, _iv, []byte(content))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(aesData), nil
}

func (a *Aes) Decode(key, iv, content string) (string, error) {
	_key, err := a.buildSecretKey(key)
	if err != nil {
		return "", err
	}
	_iv, err := a.buildSecretKey(iv)
	if err != nil {
		return "", err
	}
	base64decodedData, err := base64.StdEncoding.DecodeString(content)
	aesData, err := aes.AesDecodeWithKey(_key, _iv, base64decodedData)
	if err != nil {
		return "", err
	}
	return string(aesData), nil
}

func (a *Aes) buildSecretKey(key string) (new []byte, err error) {
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
	for _, k := range key {
		b, err := strToByte(string(k))
		if err != nil {
			return nil, err
		}
		new = append(new, b...)
	}
	return
}
