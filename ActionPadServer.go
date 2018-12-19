package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("ActionPad Server")
	server := Server{}
	err := server.runOnDeviceIP(2960)
	if err != nil {
		log.Fatal(err)
	}
}
