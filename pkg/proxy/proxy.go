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
	wlog.Info("proxyInfo: %v", proxyInfo)
	if len(proxyInfo) != 5 {
		err = errors.New("proxyInfo params number neq 5")
		return
	}

	inputConfig := InputParams{
		Ip:       proxyInfo[3],
		Port:     proxyInfo[4],
		Username: proxyInfo[1],
		Password: proxyInfo[2],
		Type:     strings.ToLower(proxyInfo[0]),
	}
	switch inputConfig.Type {
	case SSH:
		var client *ssh.Client
		sshConfig := getSshConfig(inputConfig)
		client, err = ssh.Dial("tcp", net.JoinHostPort(inputConfig.Ip, inputConfig.Port), sshConfig)
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
		res, err = sendRequest(httpClient)

	case SSL:
		var dialer netProxy.Dialer
		ssl := mySsl.DialSsl{}
		if inputConfig.Username == "" || inputConfig.Password == "" {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(inputConfig.Ip, inputConfig.Port), nil, ssl)
		} else {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(inputConfig.Ip, inputConfig.Port), &netProxy.Auth{User: inputConfig.Username, Password: inputConfig.Password}, ssl)
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
		res, err = sendRequest(httpClient)

	case SOCKS5:
		var dialer netProxy.Dialer
		if inputConfig.Username == "" || inputConfig.Password == "" {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(inputConfig.Ip, inputConfig.Port), nil, netProxy.Direct)
		} else {
			dialer, err = netProxy.SOCKS5("tcp", net.JoinHostPort(inputConfig.Ip, inputConfig.Port), &netProxy.Auth{User: inputConfig.Username, Password: inputConfig.Password}, netProxy.Direct)
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
		res, err = sendRequest(httpClient)

	case HTTPS:
		var proxy *url.URL
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", inputConfig.Username, inputConfig.Password, inputConfig.Ip, inputConfig.Port)
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
		res, err = sendRequest(httpClient)

	case HTTP:
		var proxy *url.URL
		proxyURL := fmt.Sprintf("http://%s:%s@%s:%s", inputConfig.Username, inputConfig.Password, inputConfig.Ip, inputConfig.Port)
		proxy, err = url.Parse(proxyURL)
		if err != nil {
			err = errors.Wrap(err, "url.Parse(proxyURL)-2")
			return
		}

		httpTransport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		httpClient := &http.Client{Transport: httpTransport, Timeout: 30 * time.Second}
		res, err = sendRequest(httpClient)

	default:
		err = errors.New("this proxy type non-existent")
	}

	err = errors.WithMessagef(err, "proxy info: %v", inputConfig)
	return
}

/*??????ssh?????????*/
func getSshConfig(conf InputParams) (sshConfig *ssh.ClientConfig) {
	sshConfig = &ssh.ClientConfig{
		User:            conf.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * time.Duration(30),
	}

	/*????????????????????????*/
	sshConfig.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	return
}

func sendRequest(httpClient *http.Client) (res string, err error) {
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
			if resp != nil {
				defer resp.Body.Close()
			}
			if rErr != nil {
				wlog.Warm("httpClient.Do failed: %v", rErr)
				return
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
