package aes

import (
	"useful-tools/module/walk_ui/common"
)

func getConvertType(name string) int {
	for _, info := range convertType() {
		if info.Name == name {
			return info.Key
		}
	}
	return 1
}

func convertType() []*common.CompanyItem {
	return []*common.CompanyItem{
		{Key: 1, Name: "解密"},
		{Key: 2, Name: "加密"},
	}
}

func getText(_default, new string) string {
	if new == "" {
		return _default
	}
	return new
}
