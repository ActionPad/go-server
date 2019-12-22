package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

func configLoad() {
	executablePath, err := getExecPath()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Dir(executablePath)
	filename := "ActionPadConfig.yml"

	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	err = viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		log.Fatal("Fatal error config file: %s \n", err)
	}
}

func configGenerateServerSecret() {
	viper.Set("serverSecret", generateRandomStr(8))
}

func configInitialize() {
	executablePath, err := getExecPath()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Dir(executablePath)
	filename := "ActionPadConfig.yml"
	createFileIfNotExists(path + "/" + filename)

	configLoad()

	if !viper.IsSet("serverSecret") {
		configGenerateServerSecret()
	}

	viper.SetDefault("devices", map[string]string{})
	viper.SetDefault("port", 2960)
	viper.SetDefault("saveDevices", true)

	viper.Set("ip", "")
	viper.Set("runningPort", nil)
	viper.Set("runningHost", nil)
	configSave()
}

func clearActiveServer() {
	viper.Set("runningPort", nil)
	viper.Set("runningHost", nil)
	configSave()
}

func setActiveServer(host string, port int) {
	viper.Set("runningPort", port)
	viper.Set("runningHost", host)
	configSave()
}

func setDesiredServer(host string, port int) {
	viper.Set("port", port)
	viper.Set("ip", host)
	configSave()
}

func clearDesiredServer() {
	viper.Set("port", 2960)
	viper.Set("ip", nil)
	configSave()
}

func configSave() {
	viper.WriteConfig()
}

func watchConfig(run func(e fsnotify.Event)) {
	viper.WatchConfig()
	viper.OnConfigChange(run)
}

func configCheckDevice(deviceID string) bool {
	configLoad()
	fmt.Println("Config file used", viper.ConfigFileUsed())
	if viper.GetBool("saveDevices") {
		devices := viper.GetStringMap("devices")
		fmt.Println("Check device:", strings.ToLower(deviceID), devices[deviceID], devices)
		if devices[strings.ToLower(deviceID)] != nil {
			return true
		}
	}
	return false
}

func configSaveDevice(deviceName string, deviceID string) {
	configLoad()
	if viper.GetBool("saveDevices") {
		devices := viper.GetStringMap("devices")
		devices[deviceID] = deviceName
		viper.Set("devices", devices)
		configSave()
	}
}

func configUnsaveDevice(deviceID string) {
	configLoad()
	devices := viper.GetStringMap("devices")
	delete(devices, strings.ToLower(deviceID))
	viper.Set("devices", devices)
	configSave()
}
