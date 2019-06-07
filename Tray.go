package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/skratchdot/open-golang/open"
)

func showTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("ActionPad Server")
	mConfig := systray.AddMenuItem("Config", "Config")

	// Sets the icon of a menu item. Only available on Mac.
	mConfig.SetIcon(icon.Data)
	go func() {
		<-mConfig.ClickedCh
		open.Run("https://actionpad.co/")
	}()
}

func onExit() {
	// clean up here
}
