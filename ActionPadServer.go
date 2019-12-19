package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type ActionPadInstanceManager struct {
	configurator *os.Process
	engine       *os.Process
}

func launchEngine(host string, port int) {
	fmt.Println("== ActionPad Server Engine ==")

	server := Server{}

	configInitialize()

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

func getExecPath() (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return executablePath, nil
}

func main() {
	engine := flag.Bool("engine", false, "Start core server engine")
	configurator := flag.Bool("configurator", false, "Start configurator UI")
	port := flag.Int("port", 2960, "Custom server port")
	ip := flag.String("ip", "", "Custom server IP")

	flag.Parse()

	// if err {
	// 	fmt.Println("Error opening UI process.")
	// 	os.Exit(1)
	// }

	if *engine {
		launchEngine(*ip, *port)
	} else if *configurator {
		fmt.Println("* Spawning configurator")
		launchConfigurator("")
	} else {
		fmt.Println("====== ActionPad Server ======")

		execPath, err := getExecPath()
		if err != nil {
			log.Fatal(err)
		}

		instanceManager := &ActionPadInstanceManager{}

		fmt.Println("* Spawning Server Engine")
		fmt.Println(execPath)
		cmd := exec.Command(execPath, "-engine")
		cmd.Start()
		instanceManager.engine = cmd.Process
		fmt.Print("Engine running on PID ")
		fmt.Println(instanceManager.engine.Pid)

		fmt.Println("* Spawning server UI")
		instanceManager.runInterface()
	}
}
