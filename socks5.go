package main

import (
	"context"
	"flag"
	"log"
	"net"

	socks5 "github.com/armon/go-socks5"
)

var socks5Conf = &socks5.Config{}
var socks5Server *socks5.Server

// 新增用户名密码配置
var socks5User string
var socks5Pass string

var socks5AuthFlagSet bool

func init() {
	setSocks5AuthFromFlag()
}

func newSocks5Server() *socks5.Server {
	authenticator := socks5.UserPassAuthenticator{
		Credentials: socks5.StaticCredentials{
			socks5User: socks5Pass,
		},
	}

	conf := &socks5.Config{
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {

			outgoingIP, err := generateRandomIPv6(cidr)
			if err != nil {
				log.Printf("Generate random IPv6 error: %v", err)
				return nil, err
			}
			outgoingIP = "[" + outgoingIP + "]"

			// 使用指定的出口 IP 地址创建连接
			localAddr, err := net.ResolveTCPAddr("tcp", outgoingIP+":0")
			if err != nil {
				log.Printf("[socks5] Resolve local address error: %v", err)
				return nil, err
			}
			dialer := net.Dialer{
				LocalAddr: localAddr,
			}
			// 强制使用tcp6，确保IPv6代理
			network = "tcp6"
			// 通过代理服务器建立到目标服务器的连接

			log.Println("[socks5]", addr, "via", outgoingIP)
			return dialer.DialContext(ctx, network, addr)
		},
		AuthMethods: []socks5.Authenticator{authenticator},
	}
	server, err := socks5.New(conf)
	if err != nil {
		log.Fatal(err)
	}
	return server
}

func setSocks5AuthFromFlag() {
	if socks5AuthFlagSet {
		return
	}
	flag.StringVar(&socks5User, "socks5-user", "user", "socks5 username")
	flag.StringVar(&socks5Pass, "socks5-pass", "pass", "socks5 password")
	socks5AuthFlagSet = true
}
