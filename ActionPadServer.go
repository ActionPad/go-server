package main

import (
	"fmt"
	"log"

	"github.com/faiface/mainthread"
)

func main() {
	fmt.Println("ActionPad Server")
	mainthread.Run(run)
}

func run() {
	server := Server{}
	err := server.runOnDeviceIP(2960)
	if err != nil {
		log.Fatal(err)
	}

	mainthread.CallNonBlock(func() {
		showTray()
	})
}
