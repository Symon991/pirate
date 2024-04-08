package config

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

type ConfigHandler struct {
	Path   string
	Config Config
}

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

var configHandler *ConfigHandler

func LoadConfig() (*ConfigHandler, error) {

	if configHandler != nil {
		return configHandler, nil
	}

	var path string = "config.json"

	basepath, _ := os.Executable()
	configString, err := os.ReadFile(path)
	if err != nil {
		path = filepath.Join(filepath.Dir(basepath), "config.json")
		configString, err = os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("can't load config.json: %w", err)
		}
	}

	var config Config
	err = json.Unmarshal(configString, &config)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal config.json: %w", err)
	}
	return &ConfigHandler{
		Config: config,
		Path:   path,
	}, nil
}

func (c *ConfigHandler) WriteConfig() {

	configString, _ := json.MarshalIndent(c.Config, "", "  ")
	os.WriteFile(c.Path, configString, fs.ModePerm)
}

func (c *ConfigHandler) GetRemote(name string) (*Remote, error) {

	var remoteConfig *Remote
	for i := range c.Config.Remotes {
		if c.Config.Remotes[i].Name == name {
			remoteConfig = &c.Config.Remotes[i]
			return remoteConfig, nil
		}
	}
	return nil, fmt.Errorf("remote %s not found", name)
}

func (c *ConfigHandler) AddRemote(remote Remote) error {

	c.Config.Remotes = append(c.Config.Remotes, remote)
	return nil
}

func (c *ConfigHandler) DeleteRemote(name string) error {

	var index int = -1
	for i := range c.Config.Remotes {
		if c.Config.Remotes[i].Name == name {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("remote %s not found", name)
	}

	c.Config.Remotes = slices.Delete(c.Config.Remotes, index, index+1)
	return nil
}

func (c *ConfigHandler) GetSubtitleDir() string {

	if !filepath.IsAbs(c.Config.SubtitleDir) {
		basepath, _ := os.Executable()
		return filepath.Join(filepath.Dir(basepath), c.Config.SubtitleDir)
	}

	return c.Config.SubtitleDir
}
