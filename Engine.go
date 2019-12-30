package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/viper"
)

func getExecPath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return executablePath, nil
}

func launchEngine() {
	fmt.Println("== ActionPad Server Engine ==")

	server := Server{}

	configInitialize()

	host := viper.GetString("ip")
	port := viper.GetInt("port")

	if len(host) > 0 {
		err := server.run(port, host)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := server.runOnDeviceIP(port)

	if err != nil {
		log.Fatal(err)
	}
}

func (instanceManager *ActionPadInstanceManager) spawnEngine() {
	execPath, err := getExecPath()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("* Spawning Server Engine")
	cmd := exec.Command(execPath, "-engine")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	proc := cmd.Process

	fmt.Print("Engine running on PID ")
	fmt.Println(proc.Pid)

	instanceManager.statusMessage = "Server running."

	instanceManager.engine = proc

	go func(instanceManager *ActionPadInstanceManager) {
		proc.Wait()
		fmt.Println("Engine has exited.")
		if instanceManager.engine.Pid == proc.Pid {
			instanceManager.engine = nil
		}
		instanceManager.statusMessage = "Server NOT running."
		configLoad()
		clearActiveServer()
	}(instanceManager)
}
