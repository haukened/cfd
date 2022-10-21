package main

import (
	"net"

	externalip "github.com/glendc/go-external-ip"
)

func GetPublicIPv4() (ip string, err error) {
	addr, err := getIP(4)
	if err != nil {
		return
	}
	ip = addr.String()
	return
}

func getIP(protocol uint) (ip net.IP, err error) {
	consensus := externalip.DefaultConsensus(nil, nil)
	consensus.UseIPProtocol(protocol)
	ip, err = consensus.ExternalIP()
	return
}
