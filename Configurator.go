package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/spf13/viper"
)

func (instanceManager *ActionPadInstanceManager) spawnConfigurator() {
	if instanceManager.configurator != nil {
		return
	}

	execPath, err := getExecPath()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(execPath, "-configurator")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	proc := cmd.Process

	fmt.Print("Configurator running on PID ")
	fmt.Println(proc.Pid)

	instanceManager.configurator = proc

	go func(instanceManager *ActionPadInstanceManager) {
		proc.Wait()
		fmt.Println("Configurator has exited.")
		instanceManager.configurator = nil
	}(instanceManager)
}

func renderConfigurator(w fyne.Window, app fyne.App) {
	host := viper.GetString("runningHost")
	port := viper.GetString("runningPort")

	ipField := widget.NewEntry()
	portField := widget.NewEntry()

	ipField.SetText(host)
	portField.SetText(port)

	w.SetContent(widget.NewVBox(
		widget.NewLabel("IP"),
		ipField,
		widget.NewLabel("Port"),
		portField,
		widget.NewButton("Save", func() {
			portVal, err := strconv.Atoi(portField.Text)
			if err != nil {
				fmt.Println("Error converting", portField.Text)
				portVal = 2960
				portField.SetText("2960")
			}
			setDesiredServer(ipField.Text, portVal)
			fmt.Println("Saved settings")
		}),
		widget.NewButton("Set to default", func() {
			clearDesiredServer()
			ipField.SetText("")
			portField.SetText("2960")
			fmt.Println("Set to default")
		}),
		widget.NewLabel("Once the settings are saved, click Restart in the ActionPad system tray menu."),
		widget.NewLabel("For more advanced settings, edit the config file manually."),
	))

}

func launchConfigurator() {
	configLoad()

	app := app.New()

	w := app.NewWindow("ActionPad Server - Set IP/Port")

	renderConfigurator(w, app)

	os.Setenv("FYNE_SCALE", "1")

	w.SetFixedSize(true)

	w.ShowAndRun()
}
