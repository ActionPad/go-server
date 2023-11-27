package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	log "github.com/sirupsen/logrus"

	"github.com/skratchdot/open-golang/open"
)

func (instanceManager *ActionPadInstanceManager) runInterface() {
	systray.Run(instanceManager.onReady, instanceManager.onExit)
}

func (instanceManager *ActionPadInstanceManager) onReady() {
	if runtime.GOOS == "windows" {
		systray.SetIcon(WinIcon)
	} else {
		systray.SetIcon(UnixIcon)
	}

	systray.SetTooltip("ActionPad Server")
	mTitle := systray.AddMenuItem("ActionPad Server "+CURRENT_VERSION+" (by Andrew Arpasi)", "ActionPad Server")
	mStatus := systray.AddMenuItem("Status: ", "Status")
	mTitle.Disable()
	mStatus.Disable()
	systray.AddSeparator()
	mPairing := systray.AddMenuItem("Connect Devices (Enable Pairing)", "Connect Devices")
	mShowInfo := systray.AddMenuItem("Show Server Information", "Show Server Information")
	systray.AddSeparator()
	// mSettings := systray.AddMenuItem("Change IP/Port", "Server Settings")
	mConfig := systray.AddMenuItem("Edit Server Config File", "Edit Config File")
	mUpdate := systray.AddMenuItem("Check For Server Update", "Check For Server Update")
	systray.AddSeparator()
	mRestart := systray.AddMenuItem("Restart", "Restart")
	mQuit := systray.AddMenuItem("Quit", "Quit")

	configLoad()

	mStatus.SetTitle("Status: " + instanceManager.statusMessage)
	mPairing.SetTitle(getPairingButtonMessage())

	if GetBool("pairingEnabled") {
		go func() {
			time.Sleep(time.Second)
			open.Start(qrPageURL())
		}()
	}

	go func() {
		for {
			select {
			case <-mPairing.ClickedCh:
				setPairingEnabled(!GetBool("pairingEnabled"))
				configLoad()
				mPairing.SetTitle(getPairingButtonMessage())
				if GetBool("pairingEnabled") {
					open.Start(qrPageURL())
				}
			case <-mShowInfo.ClickedCh:
				configLoad()
				open.Start(qrPageURL())
			case <-mQuit.ClickedCh:
				instanceManager.onExit()
				return
			case <-mConfig.ClickedCh:
				if runtime.GOOS == "darwin" {
					open.RunWith(ConfigFileUsed(), "/System/Applications/TextEdit.app/Contents/MacOS/TextEdit")
				} else if runtime.GOOS == "windows" {
					open.RunWith(ConfigFileUsed(), "notepad.exe")
				}
			case <-mUpdate.ClickedCh:
				open.Start("https://actionpad.co/update.html?version=" + CURRENT_VERSION)
			//
			// BUG:
			//
			// Restarting while pairing mode stops pairing mode from ever working again until actual force restart.
			//
			case <-mRestart.ClickedCh:
				log.Println("Engine:", instanceManager.engine)
				if instanceManager.engine != nil {
					fmt.Print("Killing engine on PID ")
					log.Println(instanceManager.engine.Pid)
					instanceManager.engine.Kill()
					instanceManager.spawnEngine()
					mStatus.SetTitle("Status: " + instanceManager.statusMessage)
				} else {
					instanceManager.spawnEngine()
					mStatus.SetTitle("Status: " + instanceManager.statusMessage)
				}
				go func() {
					time.Sleep(time.Second)
					configLoad()
					log.Println("Should update status:", instanceManager.statusMessage)
					mStatus.SetTitle("Status: " + instanceManager.statusMessage)
					mPairing.SetTitle(getPairingButtonMessage())
				}()
			}
		}
	}()
}

func getPairingButtonMessage() string {
	if GetBool("pairingEnabled") {
		return "Done Connecting Devices (Disable Pairing)"
	}
	return "Connect Devices (Enable Pairing)"
}

func (instanceManager *ActionPadInstanceManager) onExit() {
	// clean up here
	if instanceManager.engine != nil {
		fmt.Print("Killing engine on PID ")
		log.Println(instanceManager.engine.Pid)
		instanceManager.engine.Kill()
	}

	os.Exit(0)
}
