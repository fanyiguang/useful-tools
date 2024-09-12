package model

type DialInfo struct {
	LocalIp string `json:"local_ip"`
	Network string `json:"network"`
	Host    string `json:"host"`
	Port    string `json:"port"`
}
