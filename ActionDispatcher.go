package main

import (
	"errors"

	log "github.com/sirupsen/logrus"

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

	log.Println("Dispatching Action", action)

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
		path := action.Commands[0]
		log.Println("Opening: ", path)

		open.Start(path)
		break
	}

	return nil
}
