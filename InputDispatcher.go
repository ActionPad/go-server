package main

import (
	"fmt"
	"time"
)

type InputAction struct {
}

func startInputListener(server Server) {
	ticker := time.NewTicker(3000 * time.Millisecond)
	go handleTick(ticker)
	server.inputTicker = ticker
}

func handleTick(ticker *time.Ticker) {
	for t := range ticker.C {
		fmt.Println("Tick at", t)
	}
}

func stopInputListener(server Server) {
	if server.inputTicker != nil {
		server.inputTicker.Stop()
	}
}
