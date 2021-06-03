package utils

import (
	"fmt"
	"net"
	"strings"
)

var publicIPAddr string = ""

// PublicIPAddr 公网IP
func PublicIPAddr() string {
	if publicIPAddr != "" {
		return publicIPAddr
	}
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer conn.Close()
	publicIPAddr = strings.Split(conn.LocalAddr().String(), ":")[0]
	return publicIPAddr
}

var localIPAddr string = ""

// LocalIPAddr 本地IP
func LocalIPAddr() string {
	if localIPAddr != "" {
		return localIPAddr
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIPAddr = ipnet.IP.String()
				break
			}
		}
	}
	return localIPAddr
}
