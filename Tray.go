package main

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/viper"
)

func (instanceManager *ActionPadInstanceManager) runInterface() {
	systray.RunWithAppWindow("ActionPad Server", 600, 800, instanceManager.onReady, instanceManager.onExit)
}

func (instanceManager *ActionPadInstanceManager) showQRWindow() {
	// pageContent := url.PathEscape(assembleQRPage(viper.GetString("runningHost"), viper.GetInt("runningPort")))
	// //systray.ShowAppWindow("data:text/html," + pageContent)
	systray.ShowAppWindow("http://" + viper.GetString("runningHost") + ":" + viper.GetString("runningPort") + "/info?secret=" + viper.GetString("serverSecret"))
}

func (instanceManager *ActionPadInstanceManager) onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTooltip("ActionPad Server")
	mTitle := systray.AddMenuItem("ActionPad Server 2.0 (by Andrew Arpasi)", "ActionPad Server 2.0")
	mStatus := systray.AddMenuItem("Status: ", "Status")
	mTitle.Disable()
	mStatus.Disable()
	systray.AddSeparator()
	mConnect := systray.AddMenuItem("Connect Devices", "Connect Devices")
	mSettings := systray.AddMenuItem("Change IP/Port", "Server Settings")
	mConfig := systray.AddMenuItem("Edit Server Config File", "Edit Config File")
	systray.AddSeparator()
	mRestart := systray.AddMenuItem("Restart", "Restart")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	configLoad()

	mStatus.SetTitle("Status: " + instanceManager.statusMessage)

	instanceManager.showQRWindow()

	viper.WatchConfig()

	go func() {
		for {
			select {
			case <-mConnect.ClickedCh:
				instanceManager.showQRWindow()
				break
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			case <-mSettings.ClickedCh:
				instanceManager.spawnConfigurator()
				break
			case <-mConfig.ClickedCh:
				open.Run(viper.ConfigFileUsed())
				break
			case <-mRestart.ClickedCh:
				if instanceManager.engine != nil {
					fmt.Print("Killing engine on PID ")
					fmt.Println(instanceManager.engine.Pid)
					instanceManager.engine.Kill()
					instanceManager.spawnEngine()
					mStatus.SetTitle("Status: " + instanceManager.statusMessage)
				}
				break
			}
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
