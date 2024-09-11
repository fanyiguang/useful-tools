package controller

func ProProxyCheckTemplate() string {
	return "{\n" +
		"    \"proxy\": {\n" +
		"        \"type\": \"\",\n" +
		"        \"host\": \"\",\n" +
		"        \"port\": \"\",\n" +
		"        \"username\": \"\",\n" +
		"        \"password\": \"\"\n" +
		"    },\n" +
		"    \"request\": {\n" +
		"        \"method\": \"GET\",\n" +
		"        \"urls\": [\n" +
		"            \"https://www.baidu.com\"\n" +
		"        ],\n" +
		"        \"header\": {\n" +
		"            \"Host\": [\n" +
		"                \"www.baidu.com\"\n" +
		"            ],\n" +
		"            \"Accept-Encoding\": [\n" +
		"                \"zh-CN,zh\",\n" +
		"                \"q=0.9,ko\"\n" +
		"            ]\n" +
		"        },\n" +
		"        \"body\": \"\"\n" +
		"    },\n" +
		"    \"timeout\": 10\n" +
		"}"
}

func ProPortCheckTemplate() string {
	return "{\n" +
		"    \"local_ip\": \"自动\",\n" +
		"    \"network\": \"tcp\",\n" +
		"    \"host\": \"\",\n" +
		"    \"port\": \"\"\n" +
		"}"
}

func ProDnsQueryTemplate() string {
	return "{\n" +
		"    \"server\": \"114.114.114.114\",\n" +
		"    \"domain\": \"www.baidu.com\",\n" +
		"    \"timeout\": 10,\n" +
		"    \"qtype\": \"ipv4\"\n" +
		"}"
}
