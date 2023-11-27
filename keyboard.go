package main

import (
	"regexp"
	"runtime"
	"strings"

	"github.com/go-vgo/robotgo"

	log "github.com/sirupsen/logrus"
)

func keyIsModifier(key string) bool {
	return strings.Contains(key, "ctrl") || strings.Contains(key, "control") || strings.Contains(key, "cmd") || strings.Contains(key, "alt") || strings.Contains(key, "shift") || key == "super" || key == "command"
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
	if len(keys) == 0 && len(modifiers) != 0 {
		keys = modifiers
		modifiers = nil
	}
	for _, key := range keys {
		if key != "" {
			keyPress(key, modifiers)
		}
	}
}

func keyPress(key string, modifiers []string) {
	// mutex.Lock()
	log.Printf("Key Press: <%s> Modifiers: %v (%d)\n", key, modifiers, len(modifiers))
	// check if character is "special", can't be typed in a single keypress
	match, _ := regexp.MatchString("[~!@#$%^&*()_+{}|:\"<>?]", key)
	if len(modifiers) > 0 {
		robotgo.KeyTap(key, modifiers)
	} else if match == true && len(key) == 1 {
		robotgo.TypeStr(key)
	} else {
		robotgo.KeyTap(key)
	}
	robotgo.MilliSleep(GetInt("keyDelay"))
	// mutex.Unlock()
}

func typeText(text string) {
	if len(text) == 0 {
		return
	}

	robotgo.TypeStr(text)
}

func keyHold(key string) {
	robotgo.KeyToggle(key, "down")
}

func keyRelease(key string) {
	robotgo.KeyToggle(key, "up")
	robotgo.MilliSleep(GetInt("keyDelay"))
}

// Support shortened panel key strings
func convertShortPanelKeyStr(key string) string {
	switch key {
	case "super":
		if runtime.GOOS == "darwin" {
			return "command"
		} else {
			return "cmd" // actually Windows key
		}
	case "ctrl":
		if runtime.GOOS == "darwin" {
			return "control"
		} else {
			return "ctrl"
		}
	case "del":
		return "delete"
	case "esc":
		return "escape"
	case "ins":
		return "insert"
	case "caps":
		return "capslock"
	case "pgup":
		return "pageup"
	case "pgdn":
		return "pagedown"
	case "back":
		return "backspace"
	}

	return key
}

func testKeyboard() {
	robotgo.Sleep(10)
	keyHold("w")
	robotgo.Sleep(5)
	keyRelease("w")
}
