package main

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/spf13/viper"
)

func (instanceManager *ActionPadInstanceManager) runInterface() {
	systray.RunWithAppWindow("ActionPad Server", 600, 800, instanceManager.onReady, instanceManager.onExit)
}

func (instanceManager *ActionPadInstanceManager) showQRWindow() {
	// pageContent := url.PathEscape(assembleQRPage(viper.GetString("activeHost"), viper.GetInt("activePort")))
	// //systray.ShowAppWindow("data:text/html," + pageContent)
	systray.ShowAppWindow("http://" + viper.GetString("activeHost") + ":" + viper.GetString("activePort") + "/info?secret=" + viper.GetString("serverSecret"))
}

func (instanceManager *ActionPadInstanceManager) onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("ActionPad Server")
	mConfig := systray.AddMenuItem("Server Settings", "Server Settings")
	mConnect := systray.AddMenuItem("Connect Devices", "Connect Devices")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	configLoad()

	instanceManager.showQRWindow()

	go func() {
		for true {
			<-mConnect.ClickedCh
			instanceManager.showQRWindow()
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
