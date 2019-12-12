package main

import (
	"strconv"
	"github.com/skratchdot/open-golang/open"
	"github.com/go-vgo/robotgo"
)

type Action struct {
	Type     string   `json:"type"`
	Commands []string `json:"commands"`
}

func (action Action) dispatch() error {
	switch action.Type {
	case "keyboard":
		keyPressSequence(action.Commands)
		break
	case "mouse":
		mouseEventSequence(action.Commands)
		break
	case "open":
		open.Run(action.Commands[0])
		break
	case "delay":
		delayStr := action.Commands[0]
		duration, err := strconv.ParseFloat(delayStr, 64)
		durationMs := int(duration * 1000)
		robotgo.MilliSleep(durationMs)
		if err != nil {
			return err
		}
		break
	}

	return nil
}
