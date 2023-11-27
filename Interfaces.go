package main

import (
	"net"
	"strings"
)

type InterfaceInfo struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	QrData string `json:"qr"`
}

func getInterfaceInfo() ([]InterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var interfaceInfos []InterfaceInfo
	ipOverride := GetString("ipOverride")

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if strings.Contains(v.IP.String(), "169.254") || v.IP.To4() == nil || v.IP.String() == "127.0.0.1" {
					continue
				}
				if len(ipOverride) > 0 && v.IP.String() != ipOverride {
					continue
				}

				interfaceInfo := InterfaceInfo{
					Name:   iface.Name,
					IP:     v.IP.String(),
					QrData: qrServerInfoString(v.IP.String()),
				}
				interfaceInfos = append(interfaceInfos, interfaceInfo)
			}
		}
	}

	return interfaceInfos, nil
}
