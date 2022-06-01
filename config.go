package main

import (
	"encoding/json"
	"io/fs"
	"os"
)

type Config struct {
	Remotes []Remote
}

type Remote struct {
	Name string
	Url  string
}

func readConfig() Config {

	var config Config
	configString, _ := os.ReadFile("config.json")
	json.Unmarshal(configString, &config)
	return config
}

func writeConfig(config Config) {

	configString, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile("config.json", configString, fs.ModePerm)
}

func getRemote(name string) Remote {
	config := readConfig()
	var remoteConfig Remote
	for i := range config.Remotes {
		if config.Remotes[i].Name == name {
			remoteConfig = config.Remotes[i]
		}
	}
	return remoteConfig
}
