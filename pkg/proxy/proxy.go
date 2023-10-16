package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"useful-tools/common/config"
	mySsl "useful-tools/pkg/ssl"
	"useful-tools/pkg/wlog"

	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"

	netProxy "golang.org/x/net/proxy"
)

func SendHttpRequestByProxy(proxyInfo ...string) (res string, err error) {
	proxyConfig, err := buildConfig(proxyInfo)
	if err != nil {
		return "", err
	}

	switch proxyConfig.Type {
	case SSH:
		var client *ssh.Client
		sshConfig := getSshConfig(proxyConfig)
		client, err = ssh.Dial("tcp", net.JoinHostPort(proxyConfig.Ip, proxyConfig.Port), sshConfig)
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
		res, err = sendRequest(httpClient, proxyConfig.ReqUrls)

	case SSL:
		var dialer netProxy.Dialer
		ssl := mySsl.DialSsl{}
		if proxyConfig.Username == "" || proxyConfig.Password == "" {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(proxyConfig.Ip, proxyConfig.Port), nil, ssl)
		} else {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(proxyConfig.Ip, proxyConfig.Port), &netProxy.Auth{User: proxyConfig.Username, Password: proxyConfig.Password}, ssl)
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
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		res, err = sendRequest(httpClient, proxyConfig.ReqUrls)

	case SOCKS5:
		var dialer netProxy.Dialer
		if proxyConfig.Username == "" || proxyConfig.Password == "" {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(proxyConfig.Ip, proxyConfig.Port), nil, netProxy.Direct)
		} else {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(proxyConfig.Ip, proxyConfig.Port), &netProxy.Auth{User: proxyConfig.Username, Password: proxyConfig.Password}, netProxy.Direct)
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
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		res, err = sendRequest(httpClient, proxyConfig.ReqUrls)

	case HTTPS:
		var proxy *url.URL
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", proxyConfig.Username, proxyConfig.Password, proxyConfig.Ip, proxyConfig.Port)
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
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		res, err = sendRequest(httpClient, proxyConfig.ReqUrls)

	case HTTP:
		var proxy *url.URL
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", proxyConfig.Username, proxyConfig.Password, proxyConfig.Ip, proxyConfig.Port)
		proxy, err = url.Parse(proxyURL)
		if err != nil {
			err = errors.Wrap(err, "url.Parse(proxyURL)-2")
			return
		}

		httpTransport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		res, err = sendRequest(httpClient, proxyConfig.ReqUrls)

	default:
		err = errors.New("this proxy type non-existent")
	}

	err = errors.WithMessagef(err, "proxy info: %v", proxyConfig)
	return
}

func buildConfig(proxyInfo []string) (InputParams, error) {
	wlog.Info("proxyInfo: %v", proxyInfo)
	if len(proxyInfo) != 6 {
		return InputParams{}, errors.New("proxyInfo params number neq 5")
	}
	if proxyInfo[5] == "" {
		proxyInfo[5] = CheckIpUrls
	}
	return InputParams{
		Ip:       proxyInfo[3],
		Port:     proxyInfo[4],
		Username: proxyInfo[1],
		Password: proxyInfo[2],
		Type:     strings.ToLower(proxyInfo[0]),
		ReqUrls:  strings.Split(proxyInfo[5], ";"),
	}, nil
}

/*获取ssh的配置*/
func getSshConfig(conf InputParams) (sshConfig *ssh.ClientConfig) {
	sshConfig = &ssh.ClientConfig{
		User:            conf.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * time.Duration(30),
	}

	/*选择登录认证类型*/
	sshConfig.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	return
}

func sendRequest(httpClient *http.Client, CheckIpUrls []string) (res string, err error) {
	resCh := make(chan string)
	for _, _url := range CheckIpUrls {
		go func(url string) {
			var req *http.Request
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				wlog.Warm("http.NewRequest failed: %v", err)
				return
			}
			req.Header.Set("Connection", "Close")
			resp, rErr := httpClient.Do(req)
			if rErr != nil {
				wlog.Warm("httpClient.Do failed: %v", rErr)
				return
			}
			if resp != nil {
				defer resp.Body.Close()
			}
			if resp.StatusCode != http.StatusOK {
				wlog.Warm("resp.StatusCode: %v not StatusOK", resp.StatusCode)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			compile := regexp.MustCompile(config.GetRegIpRule())
			strBody := string(body)
			Ip := compile.FindString(strings.TrimSpace(strBody))
			if Ip != "" {
				select {
				case <-time.After(2 * time.Second):
				case resCh <- strBody:
				}
			}
		}(_url)
	}

	select {
	case res = <-resCh:
	case <-time.After(15 * time.Second):
		err = errors.Wrap(errors.New("http request timeout"), "")
	}
	return
}
