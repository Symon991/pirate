package config

import (
	"encoding/json"
	"io/fs"
	"os"
)

type Config struct {
	Remotes     []Remote
	SubtitleDir string
}

type Remote struct {
	Name string
	Url  string
}

func ReadConfig() Config {

	var config Config
	configString, _ := os.ReadFile("config.json")
	json.Unmarshal(configString, &config)
	return config
}

func WriteConfig(config Config) {

	configString, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile("config.json", configString, fs.ModePerm)
}

func GetRemote(name string) Remote {
	config := ReadConfig()
	var remoteConfig Remote
	for i := range config.Remotes {
		if config.Remotes[i].Name == name {
			remoteConfig = config.Remotes[i]
			return remoteConfig
		}
	}
	return remoteConfig
}
