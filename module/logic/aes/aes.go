package aes

import (
	"strings"
	"useful-tools/module/logic/base"
	"useful-tools/pkg/crypto"
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

func (a *Aes) Decode(key, iv, content string) (string, error) {
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
	new = []byte(key)
	return
}
