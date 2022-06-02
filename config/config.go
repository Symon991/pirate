package config

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
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
	basepath, _ := os.Executable()
	configString, _ := os.ReadFile(filepath.Join(filepath.Dir(basepath), "config.json"))
	json.Unmarshal(configString, &config)
	return config
}

func WriteConfig(config Config) {

	configString, _ := json.MarshalIndent(config, "", "  ")
	basepath, _ := os.Executable()
	os.WriteFile(filepath.Join(filepath.Dir(basepath), "config.json"), configString, fs.ModePerm)
}

func GetRemote(name string) Remote {

	userConfig := ReadConfig()
	var remoteConfig Remote
	for i := range userConfig.Remotes {
		if userConfig.Remotes[i].Name == name {
			remoteConfig = userConfig.Remotes[i]
			return remoteConfig
		}
	}
	return remoteConfig
}

func GetSubtitleDir() string {

	userConfig := ReadConfig()

	if !filepath.IsAbs(userConfig.SubtitleDir) {
		basepath, _ := os.Executable()
		return filepath.Join(basepath, userConfig.SubtitleDir)
	}

	return userConfig.SubtitleDir
}
