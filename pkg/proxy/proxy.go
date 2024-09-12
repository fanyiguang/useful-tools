package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	netProxy "golang.org/x/net/proxy"
)

func SendHttpRequestByProxy(reqInfo RequestInfo) (res string, err error) {
	proxyConfig := reqInfo.Proxy
	request := reqInfo.Request
	for i, _url := range reqInfo.Request.Urls {
		if strings.Index(_url, "http") != 0 {
			reqInfo.Request.Urls[i] = fmt.Sprintf("http://%v", _url)
		}
	}
	switch strings.ToLower(reqInfo.Proxy.Type) {
	case SSH:
		var client *ssh.Client
		sshConfig := getSshConfig(proxyConfig.Username, proxyConfig.Password, reqInfo.Timeout)
		client, err = ssh.Dial("tcp", net.JoinHostPort(proxyConfig.Host, proxyConfig.Port), sshConfig)
		if err != nil {
			err = errors.Wrap(err, "ssh.Dial")
			return
		}
		defer client.Close()

		httpTransport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return client.Dial(network, addr)
			},
		}
		httpClient := &http.Client{Transport: httpTransport}
		res, err = sendRequest(httpClient, request.Urls, request.Method, request.Header, request.Body, reqInfo.Timeout, reqInfo.HiddenBody)

	case SSL:
		var dialer netProxy.Dialer
		ssl := DialSsl{}
		if proxyConfig.Username == "" || proxyConfig.Password == "" {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(proxyConfig.Host, proxyConfig.Port), nil, ssl)
		} else {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(proxyConfig.Host, proxyConfig.Port), &netProxy.Auth{User: proxyConfig.Username, Password: proxyConfig.Password}, ssl)
		}

		if err != nil {
			err = errors.Wrap(err, "netProxy.SOCKS5-1")
			return
		}

		httpTransport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return dialer.Dial(network, addr)
			},
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: time.Duration(reqInfo.Timeout) * time.Second}
		res, err = sendRequest(httpClient, request.Urls, request.Method, request.Header, request.Body, reqInfo.Timeout, reqInfo.HiddenBody)

	case SOCKS5:
		var dialer netProxy.Dialer
		if proxyConfig.Username == "" || proxyConfig.Password == "" {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(proxyConfig.Host, proxyConfig.Port), nil, netProxy.Direct)
		} else {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(proxyConfig.Host, proxyConfig.Port), &netProxy.Auth{User: proxyConfig.Username, Password: proxyConfig.Password}, netProxy.Direct)
		}

		if err != nil {
			err = errors.Wrap(err, "netProxy.SOCKS5-2")
			return
		}

		httpTransport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return dialer.Dial(network, addr)
			},
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: time.Duration(reqInfo.Timeout) * time.Second}
		res, err = sendRequest(httpClient, request.Urls, request.Method, request.Header, request.Body, reqInfo.Timeout, reqInfo.HiddenBody)

	case HTTPS:
		var proxy *url.URL
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", proxyConfig.Username, proxyConfig.Password, proxyConfig.Host, proxyConfig.Port)
		proxy, err = url.Parse(proxyURL)
		if err != nil {
			err = errors.Wrap(err, "url.Parse(proxyURL)-1")
			return
		}

		pool := x509.NewCertPool()
		dialer := &net.Dialer{
			Timeout: 30 * time.Second,
		}
		httpTransport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return tls.DialWithDialer(dialer, network, addr, &tls.Config{
					RootCAs:            pool,
					InsecureSkipVerify: true,
				})
			},
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: time.Duration(reqInfo.Timeout) * time.Second}
		res, err = sendRequest(httpClient, request.Urls, request.Method, request.Header, request.Body, reqInfo.Timeout, reqInfo.HiddenBody)

	case HTTP:
		var proxy *url.URL
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", proxyConfig.Username, proxyConfig.Password, proxyConfig.Host, proxyConfig.Port)
		proxy, err = url.Parse(proxyURL)
		if err != nil {
			err = errors.Wrap(err, "url.Parse(proxyURL)-2")
			return
		}

		httpTransport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: time.Duration(reqInfo.Timeout) * time.Second}
		res, err = sendRequest(httpClient, request.Urls, request.Method, request.Header, request.Body, reqInfo.Timeout, reqInfo.HiddenBody)

	case SHADOWSOCKS, "ss":
		httpTransport := &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				shadowsocks, err := NewShadowsocks(Options{
					Host:     proxyConfig.Host,
					Port:     proxyConfig.Port,
					Method:   proxyConfig.Username,
					Password: proxyConfig.Password,
				})
				if err != nil {
					return nil, err
				}
				return shadowsocks.DialContext(ctx, network, addr)
			},
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: time.Duration(reqInfo.Timeout) * time.Second}
		res, err = sendRequest(httpClient, request.Urls, request.Method, request.Header, request.Body, reqInfo.Timeout, reqInfo.HiddenBody)

	default:
		err = errors.New("this proxy type non-existent")
	}

	err = errors.WithMessagef(err, "proxy info: %v %v", proxyConfig.Host, proxyConfig.Port)
	return
}

/*获取ssh的配置*/
func getSshConfig(username, password string, timeout int) (sshConfig *ssh.ClientConfig) {
	sshConfig = &ssh.ClientConfig{
		User:            username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * time.Duration(timeout),
	}

	/*选择登录认证类型*/
	sshConfig.Auth = []ssh.AuthMethod{ssh.Password(password)}
	return
}

func sendRequest(httpClient *http.Client, CheckIpUrls []string, method string, header map[string][]string, body string, timeout int, hiddenBody bool) (res string, err error) {
	resCh := make(chan string)
	for _, _url := range CheckIpUrls {
		go func(url string) {
			var req *http.Request
			req, err := http.NewRequest(method, url, strings.NewReader(body))
			if err != nil {
				logrus.Warnf("http.NewRequest failed: %v", err)
				return
			}
			req.Header = header
			_, ok := header["Connection"]
			_, okk := header["connection"]
			if !ok && !okk {
				req.Header.Set("Connection", "Close")
			}
			resp, rErr := httpClient.Do(req)
			if rErr != nil {
				logrus.Warnf("httpClient.Do failed: %v", rErr)
				return
			}
			if resp != nil {
				defer resp.Body.Close()
			}
			if resp.StatusCode != http.StatusOK {
				logrus.Warnf("resp.StatusCode: %v not StatusOK", resp.StatusCode)
				return
			}

			var response []byte
			response, err = httputil.DumpResponse(resp, !hiddenBody)
			if err != nil {
				logrus.Warnf("DumpResponse error: %v", resp.StatusCode)
				return
			}
			if len(response) > 10240 {
				response = response[:10240]
			}

			select {
			case <-time.After(2 * time.Second):
			case resCh <- string(response):
			}
		}(_url)
	}

	select {
	case res = <-resCh:
	case <-time.After(time.Duration(timeout) * time.Second):
		err = errors.Wrap(errors.New("http request timeout"), "")
	}
	return
}
