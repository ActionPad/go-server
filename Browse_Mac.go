//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"os/exec"
)

func browseFile() (string, error) {
	out, err := exec.Command("osascript", "-e", "set apFile to choose file\nPOSIX path of apFile").Output()
	if err != nil {
		fmt.Printf("cmd.Run() failed with %s\n", err)
		return "", err
	}
	fmt.Printf("combined out:\n%s\n", string(out))

	filename := string(out[:len(out)-1])
	return filename, nil
}
