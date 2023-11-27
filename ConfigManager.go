package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	yaml "gopkg.in/yaml.v3"
)

var (
	config *ConfigManager
	once   sync.Once
)

type ConfigManager struct {
	data           map[string]interface{}
	mu             sync.RWMutex
	configFilePath string
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{
		data: make(map[string]interface{}),
	}
}

func (c *ConfigManager) load(filePath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(fileData, &c.data); err != nil {
		return err
	}

	c.configFilePath = filePath

	return nil
}

func GetConfigManager() *ConfigManager {
	once.Do(func() {
		config = &ConfigManager{
			data: make(map[string]interface{}),
		}
	})
	return config
}

func Load(filename string) error {
	return GetConfigManager().load(filename)
}

func ConfigFileUsed() string {
	return GetConfigManager().configFilePath
}

func GetString(key string) string {
	return GetConfigManager().GetString(key)
}

func GetBool(key string) bool {
	return GetConfigManager().GetBool(key)
}

func GetInt(key string) int {
	return GetConfigManager().GetInt(key)
}

func GetStringMap(key string) map[string]string {
	return GetConfigManager().GetStringMap(key)
}

func SetString(key string, value string) error {
	return GetConfigManager().SetString(key, value)
}

func SetBool(key string, value bool) error {
	return GetConfigManager().SetBool(key, value)
}

func SetInt(key string, value int) error {
	return GetConfigManager().SetInt(key, value)
}

func SetStringMap(key string, value map[string]string) error {
	return GetConfigManager().SetStringMap(key, value)
}

func (c *ConfigManager) GetString(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.data[strings.ToLower(key)]; ok {
		return val.(string)
	}
	return ""
}

func (c *ConfigManager) GetBool(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.data[strings.ToLower(key)]; ok {
		return val.(bool)
	}
	return false
}

func (c *ConfigManager) GetInt(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if val, ok := c.data[strings.ToLower(key)]; ok {
		return val.(int)
	}
	return 0
}

func (c *ConfigManager) GetStringMap(key string) map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.data[strings.ToLower(key)]; ok {
		resultMap := make(map[string]string)
		switch v := val.(type) {
		case map[string]interface{}:
			for k, innerVal := range v {
				resultMap[k] = fmt.Sprintf("%v", innerVal)
			}
		case map[string]string:
			return val.(map[string]string)
		default:
			// Handle other types or error out
		}
		return resultMap
	}

	return nil
}

func (c *ConfigManager) SetString(key string, value string) error {
	return c.set(strings.ToLower(key), value)
}

func (c *ConfigManager) SetBool(key string, value bool) error {
	return c.set(strings.ToLower(key), value)
}

func (c *ConfigManager) SetInt(key string, value int) error {
	return c.set(strings.ToLower(key), value)
}

func (c *ConfigManager) SetStringMap(key string, value map[string]string) error {
	return c.set(strings.ToLower(key), value)
}

func (c *ConfigManager) set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value

	fileData, err := yaml.Marshal(&c.data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.configFilePath, fileData, 0644)
}

func createFileIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func configLoad() {
	path := configGetPath()
	filename := path + "/ActionPadConfig.yml"

	err := Load(filename)
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
}

func configGenerateServerSecret() {
	SetString("serverSecret", generateRandomStr(8))
}

func configGetPath() string {
	path := os.Getenv("APPDATA")

	if runtime.GOOS == "darwin" {
		path = os.Getenv("HOME") + "/Library/Application Support"
	}

	path += "/ActionPad"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0744)
	}
	return path
}

func configInitialize() {
	path := configGetPath()
	filePath := path + "/ActionPadConfig.yml"
	createFileIfNotExists(filePath)

	Load(filePath)

	if GetString("serverSecret") == "" {
		configGenerateServerSecret()
	}

	if GetStringMap("devices") == nil {
		SetStringMap("devices", map[string]string{})
	}

	if GetInt("port") == 0 {
		SetInt("port", 2960)
	}

	if len(GetStringMap("devices")) == 0 {
		SetBool("pairingEnabled", true)
	} else {
		SetBool("pairingEnabled", false)
	}

	if GetInt("keyDelay") == 0 {
		SetInt("keyDelay", 100)
	}
	if GetInt("mouseDelay") == 0 {
		SetInt("mouseDelay", 25)
	}
	if GetString("ipOverride") == "" {
		log.Println("Resetting IP")
		SetString("ipOverride", "")
	}
}

func setPairingEnabled(enabled bool) {
	configLoad()
	SetBool("pairingEnabled", enabled)
}

func configCheckDevice(deviceID string) bool {
	configLoad()
	devices := GetStringMap("devices")
	if devices == nil {
		return false
	}
	if _, exists := devices[strings.ToLower(deviceID)]; exists {
		return true
	}
	return false
}

func configSaveDevice(deviceName string, deviceID string) {
	configLoad()
	devices := GetStringMap("devices")
	if devices == nil {
		devices = make(map[string]string)
	}
	devices[strings.ToLower(deviceID)] = deviceName
	SetStringMap("devices", devices)
}

func configUnsaveDevice(deviceID string) {
	configLoad()
	devices := GetStringMap("devices")
	if devices != nil {
		delete(devices, strings.ToLower(deviceID))
		SetStringMap("devices", devices)
	}
}
