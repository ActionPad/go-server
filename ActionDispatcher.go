package main

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
	}

	return nil
}
