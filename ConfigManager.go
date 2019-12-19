package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func createFileIfNotExists(path string) error {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

func configInitialize() {
	executablePath, err := getExecPath()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Dir(executablePath)
	createFileIfNotExists(path + "/ActionPadConfig.json")
}

func configGetValue(key string) {

}

func configPutValue(key string) {
	viper.SafeWriteConfig()
}

func watchConfig(run func(e fsnotify.Event)) {
	viper.WatchConfig()
	viper.OnConfigChange(run)
}
