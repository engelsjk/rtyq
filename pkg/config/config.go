package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Service ...
type Service struct {
	Data struct {
		Path      string `json:"path"`
		Extension string `json:"extension"`
	} `json:"data"`
	Database struct {
		Path      string `json:"path"`
		Extension string `json:"index"`
	} `json:"database"`
}

// Config ...
type Config struct {
	Port     int       `json:"port"`
	Services []Service `json:"services"`
}

// Load ...
func Load(path string) (*Config, error) {
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

	json.Unmarshal(b, &config)

	fmt.Printf("%v\n", config)

	return config, nil
}
