package ip

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var ip string

type Config struct {
	BindingHost string `json:"binding_host"`
}

func GetIp() string {
	if ip != "" {
		return ip
	}

	configPath := filepath.Join("config", "config.json")
	
	// Check if the file exists in the current working directory
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If not found, try to find it relative to the executable
		exePath, err := os.Executable()
		if err != nil {
			log.Printf("Error getting executable path: %v", err)
			return "0.0.0.0"
		}
		exeDir := filepath.Dir(exePath)
		configPath = filepath.Join(exeDir, configPath)
	}

	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("Error reading config file: %v", err)
		return "0.0.0.0"
	}

	var config Config
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Printf("Error unmarshaling config file: %v", err)
		return "0.0.0.0"
	}

	ip = config.BindingHost
	return ip
}
