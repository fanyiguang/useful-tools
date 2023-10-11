package proxy

import (
	"fmt"
	"testing"
	"time"
)

func TestSendHttpRequestByProxy(t *testing.T) {
	tokenCh := make(chan int, 15)
	for i := 0; i < 1000; i++ {
		go func() {
			tokenCh <- 1
			defer func() {
				<-tokenCh
			}()
			proxy, err := SendHttpRequestByProxy("socks5", "admin", "123", "127.0.0.1", "7777")
			fmt.Println("wenjianjia", proxy, err, (i+1)*4)
		}()
	}

	time.Sleep(time.Minute * 10)

}
