package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Config struct {
	Remotes     []Remote
	SubtitleDir string
	Sites       Sites
}

type Remote struct {
	Name     string
	Url      string
	UserName string
	Password string
}

type Sites struct {
	NyaaUrlTemplate          string
	OpensubtitlesUrlTemplate string
	PirateBayUrlTemplate     string
	LeetxUrlTemplate         string
}

var config *Config

func LoadConfig() (*Config, error) {

	if config != nil {
		return config, nil
	}

	basepath, _ := os.Executable()
	configString, _ := os.ReadFile(filepath.Join(filepath.Dir(basepath), "config.json"))
	json.Unmarshal(configString, &config)
	return config, nil
}

func GetConfig() *Config {

	return config
}

func WriteConfig() {

	configString, _ := json.MarshalIndent(config, "", "  ")
	basepath, _ := os.Executable()
	os.WriteFile(filepath.Join(filepath.Dir(basepath), "config.json"), configString, fs.ModePerm)
}

func GetRemote(name string) (*Remote, error) {

	var remoteConfig *Remote
	for i := range config.Remotes {
		if config.Remotes[i].Name == name {
			remoteConfig = &config.Remotes[i]
			return remoteConfig, nil
		}
	}
	return nil, fmt.Errorf("remote %s not found", name)
}

func GetSubtitleDir() string {

	if !filepath.IsAbs(config.SubtitleDir) {
		basepath, _ := os.Executable()
		return filepath.Join(filepath.Dir(basepath), config.SubtitleDir)
	}

	return config.SubtitleDir
}
