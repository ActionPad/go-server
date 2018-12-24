package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-vgo/robotgo"
)

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

}

func getMousePos() {

}

func mouseHold(key string) {
	robotgo.KeyToggle(key, "down")
}

func mouseRelease(key string) {
	robotgo.KeyToggle(key, "up")
	robotgo.MilliSleep(125)
}
