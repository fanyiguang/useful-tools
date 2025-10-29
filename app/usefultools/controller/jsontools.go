package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"useful-tools/app/usefultools/adapter"

	"github.com/sirupsen/logrus"
)

var _ adapter.Controller = (*JsonTools)(nil)

type JsonTools struct {
	Base
	viewText       string
	data           string
	conversionType string
}

func NewJsonTools() *JsonTools {
	return &JsonTools{}
}

func (j *JsonTools) Data() string {
	return j.data
}

func (j *JsonTools) SetData(data string) {
	j.data = data
}

func (j *JsonTools) ViewText() string {
	return j.viewText
}

func (j *JsonTools) SetViewText(viewText string) {
	j.viewText = viewText
}

// FormatJson 格式化JSON字符串
func (j *JsonTools) FormatJson(data string) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(data), "", "    ")
	if err != nil {
		logrus.Warnf("json format error: %v", err)
		return data
	}
	return buf.String()
}

// MinifyJson 压缩JSON字符串，移除所有不必要的空白
func (j *JsonTools) MinifyJson(data string) (string, error) {
	// 首先验证输入是否为有效的JSON
	var temp interface{}
	err := json.Unmarshal([]byte(data), &temp)
	if err != nil {
		return "", errors.New("无效的JSON格式")
	}

	// 压缩JSON
	var buf bytes.Buffer
	err = json.Compact(&buf, []byte(data))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RemoveEscapes 去除JSON中的转义字符
func (j *JsonTools) RemoveEscapes(data string) (string, error) {
	// 首先尝试将字符串作为JSON字符串字面量解析
	var parsedStr string
	err := json.Unmarshal([]byte(`"`+data+`"`), &parsedStr)
	if err != nil {
		// 如果失败，尝试作为完整的JSON对象解析
		var temp interface{}
		err = json.Unmarshal([]byte(data), &temp)
		if err != nil {
			return "", errors.New("无法解析为有效的JSON")
		}
		// 如果是有效的JSON对象，重新编码并去除转义
		bytes, err := json.Marshal(temp)
		if err != nil {
			return "", err
		}
		// 再次解析以去除转义
		err = json.Unmarshal(bytes, &temp)
		if err != nil {
			return "", err
		}
		// 最后重新格式化为字符串
		formatted, err := json.MarshalIndent(temp, "", "    ")
		if err != nil {
			return "", err
		}
		return string(formatted), nil
	}
	return parsedStr, nil
}

// ProcessJson 根据操作类型处理JSON
func (j *JsonTools) ProcessJson(operation, content string) (string, error) {
	logrus.Infof("json tools operation: %s", operation)

	// 移除前后空白
	content = strings.TrimSpace(content)
	if content == "" {
		return "", errors.New("请输入JSON内容")
	}

	switch operation {
	case "格式化":
		return j.FormatJson(content), nil
	case "压缩":
		return j.MinifyJson(content)
	case "去除转义":
		return j.RemoveEscapes(content)
	case "增加转义":
		return j.AddEscapes(content)
	default:
		return "", errors.New("不支持的操作类型")
	}
}

// AddEscapes 为JSON字符串添加转义字符
func (j *JsonTools) AddEscapes(data string) (string, error) {
	// 将字符串作为JSON字符串字面量进行编码，自动添加转义
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	// 去除外层的引号
	return string(encoded[1 : len(encoded)-1]), nil
}

// GetOperations 获取支持的操作列表
func (j *JsonTools) GetOperations() []string {
	return []string{"格式化", "压缩", "去除转义", "增加转义"}
}

func (a *JsonTools) ConversionType() string {
	return a.conversionType
}

func (a *JsonTools) SetConversionType(conversionType string) {
	a.conversionType = conversionType
}

// ClearCache 清除缓存数据
func (j *JsonTools) ClearCache() {
	j.data = ""
	j.viewText = ""
}
