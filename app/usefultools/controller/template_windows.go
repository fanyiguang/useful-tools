package controller

import "useful-tools/app/usefultools/i18n"

func ProProxyCheckTemplate() string {
	return "{\r\n" +
		"    \"proxy\": {\r\n" +
		"        \"type\": \"\",\r\n" +
		"        \"host\": \"\",\r\n" +
		"        \"port\": \"\",\r\n" +
		"        \"username\": \"\",\r\n" +
		"        \"password\": \"\"\r\n" +
		"    },\r\n" +
		"    \"request\": {\r\n" +
		"        \"method\": \"GET\",\r\n" +
		"        \"urls\": [\r\n" +
		"            \"https://www.baidu.com\"\r\n" +
		"        ],\r\n" +
		"        \"header\": {\r\n" +
		"            \"Host\": [\r\n" +
		"                \"www.baidu.com\"\r\n" +
		"            ],\r\n" +
		"            \"Accept-Encoding\": [\r\n" +
		"                \"zh-CN,zh\",\r\n" +
		"                \"q=0.9,ko\"\r\n" +
		"            ]\r\n" +
		"        },\r\n" +
		"        \"body\": \"\"\r\n" +
		"    },\r\n" +
		"    \"timeout\": 10\r\n" +
		"}"
}

func ProPortCheckTemplate() string {
	return "{\r\n" +
		"    \"local_ip\": \"" + i18n.T(i18n.KeyAuto) + "\",\r\n" +
		"    \"network\": \"tcp\",\r\n" +
		"    \"host\": \"\",\r\n" +
		"    \"port\": \"\"\r\n" +
		"}"
}

func ProDnsQueryTemplate() string {
	return "{\r\n" +
		"    \"server\": \"114.114.114.114\",\r\n" +
		"    \"domain\": \"www.baidu.com\",\r\n" +
		"    \"timeout\": 10,\r\n" +
		"    \"qtype\": \"ipv4\"\r\n" +
		"}"
}
