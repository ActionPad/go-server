package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

type ActionPadInstanceManager struct {
	engine        *os.Process
	statusMessage string
}

var CURRENT_VERSION = "3.0.0"

func logInit() {
	configFolderPath := configGetPath()
	logPath := configFolderPath + "/ActionPadServerLog.txt"
	createFileIfNotExists(logPath)

	f, ferr := os.OpenFile(logPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if ferr != nil {
		log.Fatal(ferr)
	}
	log.SetOutput(f)
}

func main() {
	engine := flag.Bool("engine", false, "Start core server engine")

	logInit()

	flag.Parse()

	if *engine {
		launchEngine()
	} else {
		log.Println("====== ActionPad Server ======")
		configInitialize()

		instanceManager := &ActionPadInstanceManager{}

		instanceManager.spawnEngine()

		log.Println("* Spawning server UI")

		instanceManager.runInterface()
	}
}
