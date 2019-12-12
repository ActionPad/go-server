package main

import (
	"crypto/rand"
	"fmt"
	"runtime"
	"os/exec"
	"github.com/sqweek/dialog"
)

func generateRandomStr(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	str := fmt.Sprintf("%X", bytes)
	return str
}

func browseFile() (string, error) {
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("osascript", "-e", "set apFile to choose file\nPOSIX path of apFile").Output()
		if err != nil {
			fmt.Printf("cmd.Run() failed with %s\n", err)
			return "", err
		}
		fmt.Printf("combined out:\n%s\n", string(out))

		filename := string(out[:len(out)-1])
		return filename, nil
	} else {
		filename, err := dialog.File().Title("ActionPad Server").Load()
		if err != nil {
			return "", err
		}

		return filename, nil
	}
}