package main

import (
	"time"
	"github.com/go-vgo/robotgo"
)

type InputDispatcher struct {
	ExecuteTimer	*time.Timer
	InputAction		Action
	Sustain 		bool
}

func (inputDispatcher InputDispatcher) startKeyboardExecute() {
	for _, command := range inputDispatcher.InputAction.Commands {
		keyHold(command)
	}
}

func (inputDispatcher InputDispatcher) stopKeyboardExecute() {
	for _, command := range inputDispatcher.InputAction.Commands {
		keyRelease(command)
	}
}

func (inputDispatcher InputDispatcher) startMouseExecute() {

}

func (inputDispatcher InputDispatcher) stopMouseExecute() {
	
}

func (inputDispatcher InputDispatcher) inputTimeout() {
	if inputDispatcher.Sustain {
		inputDispatcher.ExecuteTimer.Reset(3 * time.Second)
	} else {
		if inputDispatcher.InputAction.Type == "keyboard" {
			inputDispatcher.stopKeyboardExecute()
		} else {
			inputDispatcher.stopMouseExecute()
		}
	}
}

func (inputDispatcher InputDispatcher) sustainExecute() {
	inputDispatcher.Sustain = true
}

func (inputDispatcher InputDispatcher) startExecute() {
	inputDispatcher.ExecuteTimer = time.NewTimer(3 * time.Second) // 3 second timeout
	inputDispatcher.Sustain = false
	go func(inputDispatcher InputDispatcher) {
		if inputDispatcher.InputAction.Type == "keyboard" {
			inputDispatcher.startKeyboardExecute()
		} else {
			inputDispatcher.startMouseExecute()
		}
        <-inputDispatcher.ExecuteTimer.C
        inputDispatcher.inputTimeout()
	}(inputDispatcher)
}

func (inputDispatcher InputDispatcher) stopExecute() {
	inputDispatcher.Sustain = false
	inputDispatcher.ExecuteTimer.Stop()
}

func testInputDispatcher() {
	robotgo.Sleep(10)
	inputDispatcher := &InputDispatcher{}
	inputDispatcher.InputAction = Action{
		Type: "keyboard",
		Commands: []string{"w","a"},
	}
	inputDispatcher.startExecute()
}