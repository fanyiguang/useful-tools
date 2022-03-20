package proxy

import (
	"fmt"
	"testing"
)

func TestSendHttpRequestByProxy(t *testing.T) {
	proxy, err := SendHttpRequestByProxy("http", "", "", "", "8081")
	fmt.Println("wenjianjia", proxy, err)
}
