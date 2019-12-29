package main

import (
	"errors"

	"github.com/skratchdot/open-golang/open"
)

type Action struct {
	Type     string   `json:"type"`
	Commands []string `json:"commands"`
}

func (action Action) dispatch() error {
	if len(action.Commands) == 0 {
		return errors.New("No actions specified")
	}

	switch action.Type {
	case "keyboard":
		keyPressSequence(action.Commands)
		break
	case "mouse":
		mouseEventSequence(action.Commands)
		break
	case "text":
		typeText(action.Commands[0])
		break
	case "open":
		open.Run(action.Commands[0])
		break
	}

	return nil
}
