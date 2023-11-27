package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
)

func getExecPath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return executablePath, nil
}

func launchEngine() {
	log.Println("== ActionPad Server Engine ==")

	server := Server{}

	port := GetInt("port")

	err := server.run(port)

	if err != nil {
		robotgo.ShowAlert("ActionPad Server", "The server has failed to start. You can try manually setting an IP override, or changing the server port in the config file. To try running again, restart the server in the ActionPad tray menu.", "Ok")
		log.Fatal(err)
	}
}

func (instanceManager *ActionPadInstanceManager) spawnEngine() {
	execPath, err := getExecPath()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("* Spawning Server Engine")
	cmd := exec.Command(execPath, "-engine")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	proc := cmd.Process

	fmt.Print("Engine running on PID ")
	log.Println(proc.Pid)

	instanceManager.statusMessage = "Server running."

	instanceManager.engine = proc

	go func(instanceManager *ActionPadInstanceManager) {
		proc.Wait()
		log.Println("Engine has exited.")
		if instanceManager.engine.Pid == proc.Pid {
			instanceManager.engine = nil
			instanceManager.statusMessage = "Server NOT running."
		}
	}(instanceManager)
}
