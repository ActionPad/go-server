package main

import (
	"flag"
	"fmt"
	"os"
)

type ActionPadInstanceManager struct {
	configurator  *os.Process
	engine        *os.Process
	statusMessage string
}

var CURRENT_VERSION = "2.0.1"

func main() {
	engine := flag.Bool("engine", false, "Start core server engine")
	configurator := flag.Bool("configurator", false, "Start configurator UI")

	flag.Parse()

	if *engine {
		launchEngine()
	} else if *configurator {
		fmt.Println("* Spawning configurator")
		launchConfigurator()
	} else {
		fmt.Println("====== ActionPad Server ======")

		instanceManager := &ActionPadInstanceManager{}

		instanceManager.spawnEngine()

		fmt.Println("* Spawning server UI")

		instanceManager.runInterface()
	}
}
