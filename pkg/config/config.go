package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	ErrLoadConfigFile error = fmt.Errorf("unable to load config file")
)

// Service ...
type Set struct {
	Data struct {
		Path      string `json:"path"`
		Extension string `json:"extension"`
		ID        string `json:"id"`
	} `json:"data"`
	Database struct {
		Path  string `json:"path"`
		Index string `json:"index"`
	} `json:"database"`
	Service struct {
		Path      string `json:"path"`
		ZoomLimit int    `json:"zoom_limit"`
	} `json:"service"`
}

// Config ...
type Config struct {
	Port int   `json:"port"`
	Sets []Set `json:"sets"`
}

// New ...
func New(set Set) *Config {
	return &Config{
		Sets: []Set{set},
	}
}

// Load ...
func Load(path string) (*Config, error) {

	if path == "" {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, ErrLoadConfigFile
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, ErrLoadConfigFile
	}

	var config *Config

	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, ErrLoadConfigFile
	}

	err = Validate(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Validate ...
func Validate(cfg *Config) error {

	if cfg == nil {
		return fmt.Errorf("config not initialized")
	}

	if cfg.Sets == nil {
		return fmt.Errorf("no sets provided in config ()")
	}

	// todo: add string cleaning/checking for each item below?

	for _, set := range cfg.Sets {
		if set.Data.Path == "" {
			return fmt.Errorf("no data path in set")
		}
		if set.Data.Extension == "" {
			return fmt.Errorf("no data extension in set")
		}
		if set.Data.ID == "" {
			return fmt.Errorf("no data id in set")
		}
		if set.Database.Path == "" {
			return fmt.Errorf("no database path in set")
		}
		if set.Database.Index == "" {
			return fmt.Errorf("no database index in set")
		}
		if set.Service.Path == "" {
			return fmt.Errorf("no service path set")
		}
		if set.Service.ZoomLimit == 0 {
			fmt.Printf("warning: no zoom limit provided (%s), default z=0\n", set.Service.Path)
			return nil
		}
	}

	return nil
}
