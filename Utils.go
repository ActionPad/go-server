package main

import (
	"crypto/rand"
	"fmt"
)

func generateRandomStr(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	str := fmt.Sprintf("%X", bytes)
	return str
}
