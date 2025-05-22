package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	mrand "math/rand"
	"time"
)

var cidr string
var port int

func main() {
	setSocks5AuthFromFlag()
	flag.IntVar(&port, "port", 52122, "server port")
	flag.StringVar(&cidr, "cidr", "", "ipv6 cidr")
	flag.Parse()

	// 获取本机已分配的IPv6地址池
	addrs, err := getLocalIPv6Addrs()
	if err != nil || len(addrs) == 0 {
		log.Fatalf("No available IPv6 addresses found: %v", err)
	}
	log.Println("Available IPv6 addresses:")
	for _, addr := range addrs {
		log.Println("  ", addr)
	}

	httpPort := port
	socks5Port := port + 1

	if socks5Port > 65535 {
		log.Fatal("port too large")
	}

	socks5Server = newSocks5Server()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err := socks5Server.ListenAndServe("tcp", fmt.Sprintf("0.0.0.0:%d", socks5Port))
		if err != nil {
			log.Fatal("socks5 Server err:", err)
		}

	}()
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", httpPort), httpProxy)
		if err != nil {
			log.Fatal("http Server err", err)
		}
	}()

	log.Println("server running ...")
	log.Printf("http running on 0.0.0.0:%d", httpPort)
	log.Printf("socks5 running on 0.0.0.0:%d", socks5Port)
	log.Printf("ipv6 cidr:[%s]", cidr)
	wg.Wait()

}

// 获取本机已分配的IPv6地址池
func getLocalIPv6Addrs() ([]string, error) {
	var addrs []string
	netIfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range netIfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // 跳过未启用的网卡
		}
		addrList, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrList {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP == nil || ipNet.IP.To16() == nil || ipNet.IP.To4() != nil {
				continue // 跳过非IPv6
			}
			if ipNet.IP.IsLinkLocalUnicast() || ipNet.IP.IsLoopback() {
				continue // 跳过本地链路和回环
			}
			addrs = append(addrs, ipNet.IP.String())
		}
	}
	return addrs, nil
}

// 随机获取本机已分配的IPv6地址
func getRandomLocalIPv6() (string, error) {
	addrs, err := getLocalIPv6Addrs()
	if err != nil || len(addrs) == 0 {
		return "", err
	}
	mrand.Seed(time.Now().UnixNano())
	return addrs[mrand.Intn(len(addrs))], nil
}

// 只允许从本机已分配的IPv6地址池中选取，不再自动生成新地址
func generateRandomIPv6(_ string) (string, error) {
	return getRandomLocalIPv6()
}
