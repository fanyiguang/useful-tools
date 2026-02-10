package i18n

import "fmt"

const (
	LangZhCN = "zh-CN"
	LangEnUS = "en-US"
)

var language = LangZhCN

var translations = map[string]map[Key]string{
	LangZhCN: {
		KeyMenuFile:                "文件",
		KeyMenuClearCache:          "清除所有缓存",
		KeyDialogClearCacheTitle:   "清除所有缓存",
		KeyDialogClearCacheMessage: "请问您要清理所有缓存吗?",
		KeyDialogYes:               "是",
		KeyDialogNo:                "否",
		KeyMenuMode:                "模式",
		KeyMenuProMode:             "专业模式",
		KeyMenuView:                "视图",
		KeyMenuHideBody:            "隐藏响应体",
		KeyMenuStyle:               "样式",
		KeyStyleDefault:            "默认",
		KeyStyleLowSaturationGreen: "低饱和绿",
		KeyStyleWarmLuxury:         "轻奢暖调",
		KeyStyleNeutralMinimal:     "中性极简",
		KeyMenuAction:              "功能",
		KeyMenuSaveAesKey:          "保存AES密钥",
		KeyMenuCloseUpgrade:        "关闭自动更新",
		KeyMenuLanguage:            "语言",
		KeyLangChinese:             "中文",
		KeyLangEnglish:             "英语",
		KeyMenuHelp:                "帮助",
		KeyMenuFeedback:            "反馈",
		KeyMenuHelpHome:            "帮助",
		KeySystrayOpen:             "打开工具",

		KeyButtonClear:   "清空",
		KeyButtonCheck:   "检测",
		KeyButtonCopy:    "复制",
		KeyButtonProcess: "处理",

		KeyHintRequired:       "必填",
		KeyHintOptional:       "选填",
		KeyHintSelectRequired: "必选",

		KeyErrorPrefix: "错误: ",

		KeyAuto:    "自动",
		KeyDefault: "默认",

		KeyPageDraftTitle: "草稿搭子",
		KeyPageDraftIntro: "不是很正经的草稿纸",
		KeyPageProxyTitle: "代理检测",
		KeyPageProxyIntro: "多协议代理可用性检测",
		KeyPagePortTitle:  "端口检测",
		KeyPagePortIntro:  "TCP/UDP端口检测",
		KeyPageDnsTitle:   "DNS查询",
		KeyPageDnsIntro:   "一个简单的DNS查询",
		KeyPageAesTitle:   "AES转换",
		KeyPageAesIntro:   "加密，解密，解密，加密",
		KeyPageJsonTitle:  "JSON工具",
		KeyPageJsonIntro:  "JSON格式化、压缩和去除转义工具",

		KeyDraftTabTitle: "草稿%d",

		KeyAesResultPlaceholder:   "解析结果",
		KeyAesKeyGroupPlaceholder: "可以设置密钥名称后保存",
		KeyAesDataPlaceholder:     "参数",
		KeyAesConversionTypeLabel: "转换类型:",
		KeyAesKeyNameLabel:        "密钥名称:",
		KeyAesKeyLabel:            "AES KEY:",
		KeyAesIvLabel:             "AES IV:",
		KeyAesParamLabel:          "参数:",
		KeyAesDecrypt:             "解密",
		KeyAesEncrypt:             "加密",
		KeyAesParamRequiredError:  "参数不可为空",
		KeyAesConversionTypeError: "转换类型错误",

		KeyJsonResultPlaceholder:         "处理结果",
		KeyJsonInputPlaceholder:          "请输入JSON内容",
		KeyJsonOperationLabel:            "操作类型:",
		KeyJsonInputLabel:                "内容输入:",
		KeyJsonFormat:                    "格式化",
		KeyJsonMinify:                    "压缩",
		KeyJsonUnescape:                  "去除转义",
		KeyJsonEscape:                    "增加转义",
		KeyJsonInvalidError:              "无效的JSON格式",
		KeyJsonEmptyError:                "请输入JSON内容",
		KeyJsonUnsupportedOperationError: "不支持的操作类型",

		KeyProxyResultPlaceholder:     "检测结果",
		KeyProxyHostPlaceholder:       "代理地址",
		KeyProxyPortPlaceholder:       "代理端口",
		KeyProxyUsernamePlaceholder:   "代理账号",
		KeyProxyPasswordPlaceholder:   "代理密码",
		KeyProxyURLPlaceholder:        "代理URL",
		KeyProxyTestURLPlaceholder:    "检测地址",
		KeyProxyTypeLabel:             "代理类型:",
		KeyProxyHostLabel:             "代理地址:",
		KeyProxyPortLabel:             "代理端口:",
		KeyProxyUsernameLabel:         "代理账号:",
		KeyProxyPasswordLabel:         "代理密码:",
		KeyProxyTestURLLabel:          "检测地址:",
		KeyProxyCheckInProgress:       "正在检测中，请稍等",
		KeyProxyInvalidHostError:      "地址格式错误！",
		KeyProxyInvalidPortError:      "代理端口错误！",
		KeyProxyInvalidPortRangeError: "代理端口不在合法范围内！",

		KeyPortResultPlaceholder:     "检测结果",
		KeyPortHostPlaceholder:       "目标地址",
		KeyPortPortPlaceholder:       "目标端口",
		KeyPortNetworkLabel:          "协议类型:",
		KeyPortInterfaceLabel:        "本地网卡:",
		KeyPortHostLabel:             "目标地址:",
		KeyPortPortLabel:             "目标端口:",
		KeyPortInvalidHostError:      "地址格式错误！",
		KeyPortInvalidPortError:      "端口错误！",
		KeyPortInvalidPortRangeError: "端口不在合法范围内！",

		KeyDnsResultPlaceholder:  "检测结果",
		KeyDnsHostPlaceholder:    "DNS地址",
		KeyDnsDomainPlaceholder:  "解析域名",
		KeyDnsHostLabel:          "DNS地址:",
		KeyDnsDomainLabel:        "解析域名:",
		KeyDnsInvalidDomainError: "域名格式错误！",
	},
	LangEnUS: {
		KeyMenuFile:                "File",
		KeyMenuClearCache:          "Clear All Cache",
		KeyDialogClearCacheTitle:   "Clear All Cache",
		KeyDialogClearCacheMessage: "Clear all cache?",
		KeyDialogYes:               "Yes",
		KeyDialogNo:                "No",
		KeyMenuMode:                "Mode",
		KeyMenuProMode:             "Pro Mode",
		KeyMenuView:                "View",
		KeyMenuHideBody:            "Hide Response Body",
		KeyMenuStyle:               "Style",
		KeyStyleDefault:            "Default",
		KeyStyleLowSaturationGreen: "Low Saturation Green",
		KeyStyleWarmLuxury:         "Warm Luxury",
		KeyStyleNeutralMinimal:     "Neutral Minimal",
		KeyMenuAction:              "Actions",
		KeyMenuSaveAesKey:          "Save AES Key",
		KeyMenuCloseUpgrade:        "Disable Auto Update",
		KeyMenuLanguage:            "Language",
		KeyLangChinese:             "Chinese",
		KeyLangEnglish:             "English",
		KeyMenuHelp:                "Help",
		KeyMenuFeedback:            "Feedback",
		KeyMenuHelpHome:            "Help",
		KeySystrayOpen:             "Open Tools",

		KeyButtonClear:   "Clear",
		KeyButtonCheck:   "Check",
		KeyButtonCopy:    "Copy",
		KeyButtonProcess: "Process",

		KeyHintRequired:       "Required",
		KeyHintOptional:       "Optional",
		KeyHintSelectRequired: "Required",

		KeyErrorPrefix: "Error: ",

		KeyAuto:    "Auto",
		KeyDefault: "Default",

		KeyPageDraftTitle: "Draft Pad",
		KeyPageDraftIntro: "A not-so-serious scratchpad",
		KeyPageProxyTitle: "Proxy Check",
		KeyPageProxyIntro: "Multi-protocol proxy availability check",
		KeyPagePortTitle:  "Port Check",
		KeyPagePortIntro:  "TCP/UDP port check",
		KeyPageDnsTitle:   "DNS Query",
		KeyPageDnsIntro:   "A simple DNS query",
		KeyPageAesTitle:   "AES Conversion",
		KeyPageAesIntro:   "Encrypt, decrypt, decrypt, encrypt",
		KeyPageJsonTitle:  "JSON Tools",
		KeyPageJsonIntro:  "JSON format, minify, and unescape tools",

		KeyDraftTabTitle: "Draft %d",

		KeyAesResultPlaceholder:   "Result",
		KeyAesKeyGroupPlaceholder: "Set a key name to save",
		KeyAesDataPlaceholder:     "Input",
		KeyAesConversionTypeLabel: "Mode:",
		KeyAesKeyNameLabel:        "Key Name:",
		KeyAesKeyLabel:            "AES Key:",
		KeyAesIvLabel:             "AES IV:",
		KeyAesParamLabel:          "Input:",
		KeyAesDecrypt:             "Decrypt",
		KeyAesEncrypt:             "Encrypt",
		KeyAesParamRequiredError:  "Input cannot be empty",
		KeyAesConversionTypeError: "Invalid conversion mode",

		KeyJsonResultPlaceholder:         "Output",
		KeyJsonInputPlaceholder:          "Enter JSON",
		KeyJsonOperationLabel:            "Operation:",
		KeyJsonInputLabel:                "Input:",
		KeyJsonFormat:                    "Format",
		KeyJsonMinify:                    "Minify",
		KeyJsonUnescape:                  "Unescape",
		KeyJsonEscape:                    "Escape",
		KeyJsonInvalidError:              "Invalid JSON format",
		KeyJsonEmptyError:                "Please enter JSON",
		KeyJsonUnsupportedOperationError: "Unsupported operation",

		KeyProxyResultPlaceholder:     "Results",
		KeyProxyHostPlaceholder:       "Proxy Host",
		KeyProxyPortPlaceholder:       "Proxy Port",
		KeyProxyUsernamePlaceholder:   "Proxy Username",
		KeyProxyPasswordPlaceholder:   "Proxy Password",
		KeyProxyURLPlaceholder:        "Proxy URL",
		KeyProxyTestURLPlaceholder:    "Test URL",
		KeyProxyTypeLabel:             "Proxy Type:",
		KeyProxyHostLabel:             "Proxy Host:",
		KeyProxyPortLabel:             "Proxy Port:",
		KeyProxyUsernameLabel:         "Proxy Username:",
		KeyProxyPasswordLabel:         "Proxy Password:",
		KeyProxyTestURLLabel:          "Test URL:",
		KeyProxyCheckInProgress:       "Checking, please wait",
		KeyProxyInvalidHostError:      "Invalid host format.",
		KeyProxyInvalidPortError:      "Invalid proxy port.",
		KeyProxyInvalidPortRangeError: "Proxy port out of range.",

		KeyPortResultPlaceholder:     "Results",
		KeyPortHostPlaceholder:       "Target Host",
		KeyPortPortPlaceholder:       "Target Port",
		KeyPortNetworkLabel:          "Protocol:",
		KeyPortInterfaceLabel:        "Local Interface:",
		KeyPortHostLabel:             "Target Host:",
		KeyPortPortLabel:             "Target Port:",
		KeyPortInvalidHostError:      "Invalid host format.",
		KeyPortInvalidPortError:      "Invalid port.",
		KeyPortInvalidPortRangeError: "Port out of range.",

		KeyDnsResultPlaceholder:  "Results",
		KeyDnsHostPlaceholder:    "DNS Server",
		KeyDnsDomainPlaceholder:  "Domain",
		KeyDnsHostLabel:          "DNS Server:",
		KeyDnsDomainLabel:        "Domain:",
		KeyDnsInvalidDomainError: "Invalid domain format.",
	},
}

func SetLanguage(lang string) {
	if _, ok := translations[lang]; ok {
		language = lang
		return
	}
	language = LangZhCN
}

func Language() string {
	return language
}

func T(key Key) string {
	if langMap, ok := translations[language]; ok {
		if val, ok := langMap[key]; ok {
			return val
		}
	}
	if langMap, ok := translations[LangZhCN]; ok {
		if val, ok := langMap[key]; ok {
			return val
		}
	}
	return string(key)
}

func Tf(key Key, args ...any) string {
	return fmt.Sprintf(T(key), args...)
}

func All(key Key) []string {
	values := make([]string, 0, len(translations))
	for _, langMap := range translations {
		if val, ok := langMap[key]; ok {
			values = append(values, val)
		}
	}
	return values
}

func Matches(key Key, value string) bool {
	if value == "" {
		return false
	}
	for _, langMap := range translations {
		if val, ok := langMap[key]; ok && val == value {
			return true
		}
	}
	return false
}
