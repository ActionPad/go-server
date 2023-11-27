package main

import (
	_ "embed"
	b64 "encoding/base64"
	"strconv"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

func qrCodeGenPNG(content string) []byte {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil
	}
	return png
}

func qrCodeBase64PNG(content string) string {
	qrCodeData := qrCodeGenPNG(content)
	base64Str := b64.StdEncoding.EncodeToString(qrCodeData)
	return base64Str
}

func qrServerInfoString(ip string) string {
	hostname := getHostname()
	if len(hostname) > 100 {
		hostname = hostname[:100]
	}
	return qrCodeBase64PNG(GetString("serverSecret") + "," + ip + "," + strconv.Itoa(GetInt("port")) + "," + hostname)
}

//go:embed html/dist/index.html
var qrPageContents string

func qrPageURL() string {
	host := "localhost"
	ipOverride := GetString("ipOverride")
	if len(ipOverride) > 0 {
		host = ipOverride
	}
	return "http://" + host + ":" + strconv.Itoa(GetInt("port")) + "/info"
}

func assembleQRPage(host string, port int, secret string) string {
	pageStr := qrPageContents
	portStr := strconv.Itoa(port)
	pageStr = strings.Replace(pageStr, "~~IP~~", host, 1)
	pageStr = strings.Replace(pageStr, "~~PORT~~", portStr, 1)
	return pageStr
}
