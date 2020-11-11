package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	// ErrConfigOpenFile error is returned when the config file cannot be opened
	ErrConfigOpenFile error = fmt.Errorf("unable to open config file")
	// ErrConfigReadFile error is reutnred when the config file cannot be read
	ErrConfigReadFile error = fmt.Errorf("unable to read config file")
	// ErrConfigInvalidStructure error is returned when the config file does not have the required data structure
	ErrConfigInvalidStructure error = fmt.Errorf("config file structure is invalid")
)

// Set is the combined set config (data/database/service) for one data type
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

// Config defines the array of data sets to create and run,
// along with the single port to run the service on
type Config struct {
	Port int   `json:"port"`
	Sets []Set `json:"sets"`
}

// New creates a new config from a single data set
func New(set Set) *Config {
	return &Config{
		Sets: []Set{set},
	}
}

// Load creates a new config from a JSON file
func Load(path string) (*Config, error) {

	if path == "" {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, ErrConfigOpenFile
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, ErrConfigReadFile
	}

	var config *Config

	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, ErrConfigInvalidStructure
	}

	err = Validate(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks a config to ensure that it is properly instantiated
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
