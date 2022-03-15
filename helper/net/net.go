package net

import "net"

func DialTcpByIFace(ip, port string, iFaceIP string) (isSuccess bool, err error) {
	rAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(ip, port))
	if err != nil {
		return false, err
	}
	lAddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(iFaceIP, "0"))
	if err != nil {
		return false, err
	}

	conn, err := net.DialTCP("tcp", lAddr, rAddr)
	if err != nil {
		return false, err
	}

	_ = conn.Close()
	return true, nil
}

func DialUdpByIFace(ip, port string, iFaceIP string) (isSuccess bool, err error) {
	rAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(ip, port))
	if err != nil {
		return false, err
	}
	lAddr, err := net.ResolveUDPAddr("udp", net.JoinHostPort(iFaceIP, "0"))
	if err != nil {
		return false, err
	}

	conn, err := net.DialUDP("udp", lAddr, rAddr)
	if err != nil {
		return false, err
	}

	_ = conn.Close()
	return true, nil
}
