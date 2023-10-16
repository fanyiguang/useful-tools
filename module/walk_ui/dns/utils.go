package dns

import (
	"fmt"
	"time"
)

func logFormat(server, host, msg string) string {
	return fmt.Sprintf("[%v] %v(%v) => %v\r\n\r\n", time.Now().Format(`06-01-02 15:04:05`), host, server, msg)
}

func getServer(server string) string {
	if server == "" {
		return "默认"
	}
	return server
}
