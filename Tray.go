package main

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
)

func (instanceManager *ActionPadInstanceManager) runInterface() {
	systray.RunWithAppWindow("ActionPad Server", 800, 600, instanceManager.onReady, instanceManager.onExit)
}

func (instanceManager *ActionPadInstanceManager) onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("ActionPad Server")
	mConfig := systray.AddMenuItem("Server Settings", "Server Settings")
	mConnect := systray.AddMenuItem("Connect Devices", "Connect Devices")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	systray.ShowAppWindow("https://actionpad.co/p/hoto")

	go func() {
		for true {
			<-mConnect.ClickedCh
			systray.ShowAppWindow("https://actionpad.co")
		}
	}()

	go func() {
		for true {
			<-mQuit.ClickedCh
			systray.Quit()
		}
	}()

	go func() {
		for true {
			<-mConfig.ClickedCh
			instanceManager.configurator = spawnConfigurator()
		}
	}()
}

func (instanceManager *ActionPadInstanceManager) onExit() {
	// clean up here
	if instanceManager.configurator != nil {
		fmt.Print("Killing configurator on PID ")
		fmt.Println(instanceManager.configurator.Pid)
		instanceManager.configurator.Kill()
	}
	if instanceManager.engine != nil {
		fmt.Print("Killing engine on PID ")
		fmt.Println(instanceManager.engine.Pid)
		instanceManager.engine.Kill()
	}

	os.Exit(0)
}
