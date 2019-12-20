package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func spawnConfigurator() *os.Process {
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
	return proc
}

func renderConfigurator(w fyne.Window, app fyne.App) {
	deviceList := widget.NewVBox()

	devices := viper.GetStringMap("devices")

	fmt.Println("devices", devices)

	for UUID, name := range devices {
		fmt.Println("Add device to list ", UUID)
		deviceRow := widget.NewHBox(
			widget.NewLabel(name.(string)),
			widget.NewButton("Delete", func() {
				configUnsaveDevice(UUID)
				renderConfigurator(w, app)
			}),
		)
		deviceList.Append(deviceRow)
	}

	w.SetContent(widget.NewVBox(
		widget.NewLabel("ActionPad Server"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
		deviceList,
	))

}

func launchConfigurator() {
	configLoad()

	app := app.New()

	w := app.NewWindow("ActionPad Server Configurator")

	renderConfigurator(w, app)

	go func() {
		watchConfig(func(e fsnotify.Event) {
			renderConfigurator(w, app)
		})
	}()

	fmt.Println("About to show")
	w.ShowAndRun()
}
