package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/spf13/viper"
)

type InputRequest struct {
	UUID        string `json:"uuid"`
	InputAction Action `json:"inputAction"`
}

type InputDispatcher struct {
	ExecuteTimer *time.Timer
	InputAction  Action
	Sustain      bool
	Running      bool
	MouseActive  map[string]bool
}

func (inputDispatcher *InputDispatcher) startKeyboardExecute() {
	for _, command := range inputDispatcher.InputAction.Commands {
		keyStr := convertShortPanelKeyStr(command)
		keyHold(keyStr)
	}
}

func (inputDispatcher *InputDispatcher) stopKeyboardExecute() {
	for _, command := range inputDispatcher.InputAction.Commands {
		keyStr := convertShortPanelKeyStr(command)
		keyRelease(keyStr)
	}
}

func (inputDispatcher *InputDispatcher) startMouseExecute() {
	for _, command := range inputDispatcher.InputAction.Commands {
		if command == "lclick" {
			if inputDispatcher.MouseActive["left"] != true {
				mouseHold("left")
				inputDispatcher.MouseActive["left"] = true
				robotgo.MilliSleep(viper.GetInt("mouseDelay"))
			}
		}
		if command == "rclick" {
			if inputDispatcher.MouseActive["right"] != true {
				mouseHold("right")
				inputDispatcher.MouseActive["right"] = true
				robotgo.MilliSleep(viper.GetInt("mouseDelay"))
			}
		}
		if strings.Contains(command, "scroll") {
			mouseScrollInputExecute(command)
		}
		if strings.Contains(command, "pointer") {
			mousePointerInputExecute(command)
		}
	}
}

func (inputDispatcher *InputDispatcher) stopMouseExecute() {
	for _, command := range inputDispatcher.InputAction.Commands {
		if command == "lclick" {
			mouseRelease("left")
			inputDispatcher.MouseActive["left"] = false
		}
		if command == "rclick" {
			mouseRelease("right")
			inputDispatcher.MouseActive["right"] = false
		}
	}
}

func mouseScrollInputExecute(command string) {
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
		robotgo.MilliSleep(viper.GetInt("mouseDelay"))
	} else {
		fmt.Printf("Command <%s> not formatted properly.\n", command)
	}

}

func mousePointerInputExecute(command string) {
	components := strings.Split(command, "::")
	direction := components[1]
	magnitude, err := strconv.Atoi(components[2])
	if err != nil || len(components) < 3 {
		fmt.Printf("Command <%s> not formatted properly.\n", command)
		return
	}
	if direction == "left" || direction == "up" {
		magnitude *= -1
	}

	var pos MousePos
	pos.X, pos.Y = robotgo.GetMousePos()

	if direction == "up" || direction == "down" {
		robotgo.MoveMouse(pos.X, pos.Y+magnitude)
	} else {
		robotgo.MoveMouse(pos.X+magnitude, pos.Y)
	}

	robotgo.MilliSleep(viper.GetInt("mouseDelay"))
}

func (inputDispatcher *InputDispatcher) inputTimeout() {
	if inputDispatcher.Sustain {
		inputDispatcher.ExecuteTimer.Reset(2 * time.Second)
		inputDispatcher.Sustain = false
		<-inputDispatcher.ExecuteTimer.C
		inputDispatcher.inputTimeout()
	} else {
		inputDispatcher.stopExecute()
	}
}

func (inputDispatcher *InputDispatcher) sustainExecute() {
	inputDispatcher.Sustain = true
}

func (inputDispatcher *InputDispatcher) startExecute() {
	if inputDispatcher.Running {
		return
	}
	inputDispatcher.ExecuteTimer = time.NewTimer(2 * time.Second) // 3 second timeout
	inputDispatcher.Sustain = false
	inputDispatcher.Running = true
	inputDispatcher.MouseActive = make(map[string]bool)
	go func(inputDispatcher *InputDispatcher) {
		for inputDispatcher.Running {
			if inputDispatcher.InputAction.Type == "keyboard" {
				inputDispatcher.startKeyboardExecute()
			} else {
				inputDispatcher.startMouseExecute()
			}
		}
	}(inputDispatcher)
	go func(inputDispatcher *InputDispatcher) {
		<-inputDispatcher.ExecuteTimer.C
		inputDispatcher.inputTimeout()
	}(inputDispatcher)
}

func (inputDispatcher *InputDispatcher) stopExecute() {
	inputDispatcher.Sustain = false
	inputDispatcher.Running = false
	inputDispatcher.ExecuteTimer.Stop()
	if inputDispatcher.InputAction.Type == "keyboard" {
		inputDispatcher.stopKeyboardExecute()
	} else {
		inputDispatcher.stopMouseExecute()
	}
}

func testInputDispatcher() {
	robotgo.Sleep(5)
	inputDispatcher := &InputDispatcher{}
	inputDispatcher.InputAction = Action{
		Type:     "keyboard",
		Commands: []string{"w", "space"}, //"rclick","pointer::down::3","pointer::left::3"},
	}
	inputDispatcher.startExecute()
}
