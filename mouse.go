package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

type MousePos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func mouseEventSequence(commands []string) {
	robotgo.Sleep(3)
	for _, command := range commands {
		if command == "lclick" {
			robotgo.MouseClick("left")
			robotgo.MilliSleep(125)
		}
		if command == "rclick" {
			robotgo.MouseClick("right")
			robotgo.MilliSleep(125)
		}
		if strings.Contains(command, "scroll") {
			mouseScrollParse(command)
		}
		if strings.Contains(command, "pointer") {
			mousePointerParse(command)
		}
	}
}

func mouseScrollParse(command string) {
	components := strings.Split(command, "::")
	// check scroll direction
	direction := "down"
	if strings.Contains(command, "up") {
		direction = "up"
	}
	// get magnitude from command string
	if magnitude, err := strconv.Atoi(components[1]); err == nil &&
		len(components) == 2 {
		robotgo.ScrollMouse(magnitude, direction)
		robotgo.MilliSleep(125)
	} else {
		fmt.Printf("Command <%s> not formatted properly.\n", command)
	}
}

func mousePointerParse(command string) {
	components := strings.Split(command, "::")
	posX, err := strconv.Atoi(components[1])
	posY, err := strconv.Atoi(components[2])
	if err != nil || len(components) < 3 {
		fmt.Printf("Command <%s> not formatted properly.\n", command)
		return
	}
	robotgo.MoveMouse(posX, posY)
	robotgo.MilliSleep(125)
}

func mouseHold(key string) {
	robotgo.KeyToggle(key, "down")
}

func mouseRelease(key string) {
	robotgo.KeyToggle(key, "up")
	robotgo.MilliSleep(125)
}
