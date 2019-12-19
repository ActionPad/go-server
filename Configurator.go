package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func spawnConfigurator() *os.Process {
	execPath, err := getExecPath()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(execPath, "-configurator")
	cmd.Start()
	proc := cmd.Process
	fmt.Print("Configurator running on PID ")
	fmt.Println(proc.Pid)
	return proc
}

func launchConfigurator(configPath string) {
	app := app.New()

	w := app.NewWindow("Hello")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))

	fmt.Println("About to show")
	w.ShowAndRun()
}
