package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Context struct {
	Context string `json:"context"`
	Width   uint   `json:"width"`
	Height  uint   `json:"height"`
}

type Config struct {
	Server struct {
		Port     string `yaml:"port"`
		Uploader struct {
			MaxSize  string    `yaml:"maxSize"`
			Allow    []string  `yaml:"allow"`
			Contexts []Context `json:"contexts"`
		} `yaml:"uploader"`
	} `yaml:"server"`
	Exporter struct {
		Port    string `yaml:"port"`
		Enabled bool   `yaml:"enabled"`
	} `yaml:"exporter"`
	Storage struct {
		Path string `yaml:"path"`
	} `yaml:"storage"`
}

func NewConfig(configPath string) *Config {
	if err := ValidateConfigPath(configPath); err != nil {
		panic(err)
	}

	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		panic(err)
	}

	file.Close()

	return config
}

func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
