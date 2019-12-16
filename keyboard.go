package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-vgo/robotgo"
)

func keyIsModifier(key string) bool {
	return strings.Contains(key, "ctrl") || strings.Contains(key, "cmd") || strings.Contains(key, "alt") || strings.Contains(key, "shift") || key == "super" || key == "command"
}

func keyPressSequence(commands []string) {
	modifiers := make([]string, 0)
	keys := make([]string, 0)
	// build modifiers array and clean commands
	for _, command := range commands {
		key := command
		// ~ suffix is remnant from older panels no longer needed
		if len(command) > 1 && strings.Contains(command, "~") {
			key = strings.TrimSuffix(command, "~")
		}

		key = convertShortPanelKeyStr(key)

		if keyIsModifier(key) == true {
			modifiers = append(modifiers, key)
		} else {
			keys = append(keys, key)
		}
	}
	for _, key := range keys {
		if key != "" {
			keyPress(key, modifiers)
		}
	}
}

func keyPress(key string, modifiers []string) {
	// mutex.Lock()
	fmt.Printf("Key Press: <%s> Modifiers: %v (%d)\n", key, modifiers, len(modifiers))
	// check if character is "special", can't be typed in a single keypress
	match, _ := regexp.MatchString("[~!@#$%^&*()_+{}|:\"<>?]", key)
	if len(modifiers) > 0 {
		robotgo.KeyTap(key, modifiers)
	} else if match == true {
		robotgo.TypeString(key)
	} else {
		robotgo.KeyTap(key)
	}
	robotgo.MilliSleep(125)
	// mutex.Unlock()
}

func keyHold(key string) {
	robotgo.KeyToggle(key, "down")
}

func keyRelease(key string) {
	robotgo.KeyToggle(key, "up")
	robotgo.MilliSleep(125)
}

func typeString(str string, cpm float64) {
	robotgo.TypeStr(str, cpm)
}

// Support shortened panel key strings
func convertShortPanelKeyStr(key string) string {
	switch key {
	case "super":
		return "cmd"
	case "del":
		return "delete"
	case "ins":
		return "insert"
	case "caps":
		return "capslock"
	case "pgup":
		return "pageup"
	case "pgdn":
		return "pagedown"
	}

	return key
}

func testKeyboard() {
	robotgo.Sleep(10)
	keyHold("w")
	robotgo.Sleep(5)
	keyRelease("w")
}