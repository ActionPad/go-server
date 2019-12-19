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
	filename := "ActionPadConfig.json"
	createFileIfNotExists(path + "/" + filename)

	viper.AddConfigPath(path)
	viper.SetConfigName(filename)

	err = viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal("Fatal error config file: %s \n", err)
	}

	viper.SetDefault("port", 2960)
	viper.Set("activePort", nil)
	viper.Set("activeHost", nil)
	configSave()
}

func setActiveServer(host string, port int) {
	viper.Set("activePort", port)
	viper.Set("activeHost", host)
	configSave()
}

func configSave() {
	viper.WriteConfig()
}

func watchConfig(run func(e fsnotify.Event)) {
	viper.WatchConfig()
	viper.OnConfigChange(run)
}
