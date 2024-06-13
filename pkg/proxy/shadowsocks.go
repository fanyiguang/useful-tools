package proxy

import (
	"context"
	"errors"
	"github.com/sagernet/sing-shadowsocks2"
	"github.com/sagernet/sing/common/bufio"
	"github.com/sagernet/sing/common/metadata"
	"net"
	"strings"
)

type Options struct {
	Host     string
	Port     string
	Method   string
	Password string
}

type Shadowsocks struct {
	host   string
	port   string
	method shadowsocks.Method
}

func NewShadowsocks(options Options) (*Shadowsocks, error) {
	method, err := shadowsocks.CreateMethod(context.Background(), options.Method, shadowsocks.MethodOptions{
		Password: options.Password,
	})
	if err != nil {
		return nil, err
	}
	return &Shadowsocks{
		host:   options.Host,
		port:   options.Port,
		method: method,
	}, nil
}

func (s *Shadowsocks) DialContext(ctx context.Context, network, addr string) (conn net.Conn, err error) {
	var dial net.Dialer
	switch strings.ToLower(network) {
	case "tcp":
		conn, err = dial.DialContext(ctx, network, net.JoinHostPort(s.host, s.port))
		if err != nil {
			return nil, err
		}
		return s.method.DialConn(conn, metadata.ParseSocksaddr(addr))
	case "udp":
		conn, err = dial.DialContext(ctx, network, net.JoinHostPort(s.host, s.port))
		if err != nil {
			return nil, err
		}
		return bufio.NewBindPacketConn(s.method.DialPacketConn(conn), metadata.ParseSocksaddr(addr)), nil
	default:
		return nil, errors.New("network not found")
	}
}

func (s *Shadowsocks) ListenPacket(ctx context.Context, addr string) (conn net.PacketConn, err error) {
	var dial net.Dialer
	uConn, err := dial.DialContext(ctx, "udp", net.JoinHostPort(s.host, s.port))
	if err != nil {
		return nil, err
	}
	return s.method.DialPacketConn(uConn), nil
}
