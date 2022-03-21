# Useful-tools v1.0.0
**useful-tools** 是方便平常开发，测试和运维的小工具集合。

## 下载
[Useful-tools v1.0.0](https://github.com/fanyiguang/useful-tools/releases/download/v1.0.0/useful-tools.zip)

## 功能介绍

### 代理检测

#### 普通模式
![image](https://github.com/fanyiguang/useful-tools/blob/master/resource/proxy-normal.png)  
支持SOCKS5，SSL，SSH，HTTP，HTTPS五种代理的检测，输入对应的参数点击检测结果就会输出到右侧的文本框中。

#### 解析模式
![image](https://github.com/fanyiguang/useful-tools/blob/master/resource/dns-parser.png)  
设置视图 > 解析模式 可以设置输入参数为解析的模式。在解析模式下输入参数支持多种格式：
1. {"[.\*]type[.\*]": "socks5", "[.\*]ip[.\*]": "127.0.0.1", "[.\*]port[.\*]": "8888", "[.\*]name[.\*]": "admin", "[.\*]pass[.\*]": "use-ful-tools"}
2. {"[.\*]type[.\*]": "socks5", "[.\*]addr[.\*]": "127.0.0.1", "[.\*]port[.\*]": "8888", "[.\*]name[.\*]": "admin", "[.\*]pass[.\*]": "use-ful-tools"}
2. type name pass ip port
3. type:name:pass:ip:port
4. type:::ip:port

如果没有账号和密码的话留空即可，[.\*]代表任意字符。

### 端口检测

#### 普通模式
![image](https://github.com/fanyiguang/useful-tools/blob/master/resource/tcp-udp-normal.png)  
支持使用自选或随机本地网卡发送TCP/UDP请求到对应的IP端口，结果会输出到右侧的文本框中。

#### 解析模式
![image](https://github.com/fanyiguang/useful-tools/blob/master/resource/tcp-udp-parser.png)  
支持解析的格式：
1. {"[.\*]network[.\*]": "tcp", "[.\*]iFace[.\*]": "127.0.0.1", "[.\*]ip[.\*]": "127.0.0.1", "[.\*]port[.\*]": "8888"}
2. {"[.\*]network[.\*]": "tcp", "[.\*]interface[.\*]": "127.0.0.1", "[.\*]ip[.\*]": "127.0.0.1", "[.\*]port[.\*]": "8888"}
3. {"[.\*]network[.\*]": "tcp", "[.\*]iFace[.\*]": "127.0.0.1", "[.\*]addr[.\*]": "127.0.0.1", "[.\*]port[.\*]": "8888"}
4. {"[.\*]network[.\*]": "tcp", "[.\*]interface[.\*]": "127.0.0.1", "[.\*]addr[.\*]": "127.0.0.1", "[.\*]port[.\*]": "8888"}
5. network interface ip port
6. network:interface:ip:port
7. network:ip:port
8. network:interface

network和interface可以忽略不写默认为TCP请求网卡随机，也可以单独忽略interface，[.\*]代表任意字符。

### DNS检测

#### 普通模式
![image](https://github.com/fanyiguang/useful-tools/blob/master/resource/dns-normal.png)  
默认使用本地配置的dns Server，也可以输入自定义的dns Server。

#### 解析模式
![image](https://github.com/fanyiguang/useful-tools/blob/master/resource/dns-parser.png)  
支持的格式：
1. {"[.\*]domain[.\*]": "www.github.com", "[.\*]dns[.\*]": "8.8.8.8"}
2. {"[.\*]domain[.\*]": "www.github.com", "[.\*]server[.\*]": "8.8.8.8"}
3. {"[.\*]url[.\*]": "www.github.com", "[.\*]server[.\*]": "8.8.8.8"}
4. {"[.\*]url[.\*]": "www.github.com", "[.\*]dns[.\*]": "8.8.8.8"}
5. url dns
6. url:dns
7. url

dns客户忽略不写默认使用本机配置的dns Server，[.\*]代表任意字符。
