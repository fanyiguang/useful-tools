package controller

func ProProxyCheckTemplate() string {
	return "{\r" +
		"    \"proxy\": {\r" +
		"        \"type\": \"\",\r" +
		"        \"host\": \"\",\r" +
		"        \"port\": \"\",\r" +
		"        \"username\": \"\",\r" +
		"        \"password\": \"\"\r" +
		"    },\r" +
		"    \"request\": {\r" +
		"        \"method\": \"GET\",\r" +
		"        \"urls\": [\r" +
		"            \"https://www.baidu.com\"\r" +
		"        ],\r" +
		"        \"header\": {\r" +
		"            \"Host\": [\r" +
		"                \"www.baidu.com\"\r" +
		"            ],\r" +
		"            \"Accept-Encoding\": [\r" +
		"                \"zh-CN,zh\",\r" +
		"                \"q=0.9,ko\"\r" +
		"            ]\r" +
		"        },\r" +
		"        \"body\": \"\"\r" +
		"    },\r" +
		"    \"timeout\": 10\r" +
		"}"
}

func ProPortCheckTemplate() string {
	return "{\r" +
		"    \"local_ip\": \"自动\",\r" +
		"    \"network\": \"tcp\",\r" +
		"    \"host\": \"\",\r" +
		"    \"port\": \"\"\r" +
		"}"
}
