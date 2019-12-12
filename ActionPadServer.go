package main

import (
	"fmt"
	"log"
	"github.com/faiface/mainthread"
)

func run() {
	fmt.Println("ActionPad Server")
	server := Server{}
	err := server.runOnDeviceIP(2960)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	mainthread.Run(run)
}
