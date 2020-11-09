package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Service ...
type Service struct {
	Data struct {
		Path      string `json:"path"`
		Extension string `json:"extension"`
		ID        string `json:"id"`
	} `json:"data"`
	Database struct {
		Path  string `json:"path"`
		Index string `json:"index"`
	} `json:"database"`
}

// Config ...
type Config struct {
	Port     int       `json:"port"`
	Services []Service `json:"services"`
}

// Create ...
func New(svc Service) *Config {
	return &Config{
		Services: []Service{svc},
	}
}

// Load ...
func Load(path string) (*Config, error) {

	if path == "" {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config *Config

	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Validate(cfg *Config) error {
	// todo: make sure all required config values are not nil
	return nil
}
