package proxy

import (
	"fmt"
	"testing"
	"time"
)

func TestSendHttpRequestByProxy(t *testing.T) {
	tokenCh := make(chan int, 50)
	for i := 0; i < 1000; i++ {
		go func() {
			tokenCh <- 1
			defer func() {
				<-tokenCh
			}()
			proxy, err := SendHttpRequestByProxy(RequestInfo{
				Proxy: struct {
					Type     string `json:"type"`
					Host     string `json:"host"`
					Port     string `json:"port"`
					Username string `json:"username"`
					Password string `json:"password"`
				}{
					Type:     "socks5",
					Host:     "127.0.0.1",
					Port:     "20000",
					Username: "admin",
					Password: "123",
				},
				Request: struct {
					Method string              `json:"method"`
					Urls   []string            `json:"urls"`
					Header map[string][]string `json:"header"`
					Body   string              `json:"body"`
				}{},
				Timeout:    0,
				HiddenBody: false,
			})
			fmt.Println("wenjianjia", proxy, err, (i+1)*4)
		}()
	}

	time.Sleep(time.Minute * 10)

}
