package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
	keytype()
	server := Server{}
	server.runOnDeviceIP(2960)
}
