package main

import (
	"fmt"
)

func main() {
	fmt.Println("ActionPad Server")
	server := Server{}
	server.runOnDeviceIP(2960)
}
