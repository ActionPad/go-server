package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func generateRandomStr(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	str := fmt.Sprintf("%X", bytes)
	return str
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		if runtime.GOOS == "darwin" {
			return "Mac"
		} else if runtime.GOOS == "windows" {
			return "Windows PC"
		}
		return "Computer"
	}

	if runtime.GOOS == "darwin" {
		hostname = strings.TrimSuffix(hostname, ".local")
	}

	return hostname
}
