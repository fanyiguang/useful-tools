package controller

type AesConversion struct {
	Base
	conversion string
}

func NewAesConversion() *AesConversion {
	return &AesConversion{}
}

func (a *AesConversion) ConversionList() []string {
	return []string{"解密", "加密"}
}
