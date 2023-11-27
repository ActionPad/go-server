//go:build windows
// +build windows

package main

import (
	"github.com/sqweek/dialog"
)

func browseFile() (string, error) {
	filename, err := dialog.File().Title("ActionPad Server").Load()
	if err != nil {
		return "", err
	}

	return filename, nil
}
