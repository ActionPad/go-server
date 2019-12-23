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
		if strings.Contains(command, "offset") {
			mouseOffsetParse(command)
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

func mouseOffsetParse(command string) {
	components := strings.Split(command, "::")
	offX, err := strconv.Atoi(components[1])
	offY, err := strconv.Atoi(components[2])
	if err != nil || len(components) < 3 {
		fmt.Printf("Command <%s> not formatted properly.\n", command)
		return
	}
	posX, posY := robotgo.GetMousePos()
	posX += offX
	posY += offY
	robotgo.MoveMouse(posX, posY)
	robotgo.MilliSleep(125)
}

func mouseHold(button string) {
	robotgo.MouseToggle("down", button)
}

func mouseRelease(button string) {
	robotgo.MouseToggle("up", button)
	robotgo.MilliSleep(125)
}
