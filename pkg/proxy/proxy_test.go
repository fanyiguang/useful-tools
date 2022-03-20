package proxy

import (
	"fmt"
	"testing"
)

func TestSendHttpRequestByProxy(t *testing.T) {
	proxy, err := SendHttpRequestByProxy("http", "", "", "39.102.37.208", "8081")
	fmt.Println("wenjianjia", proxy, err)
}
