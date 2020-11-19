package rtyq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	ErrNoConfigProvided       error  = fmt.Errorf("no config provided")
	ErrConfigOpenFile         error  = fmt.Errorf("unable to open config file")
	ErrConfigReadFile         error  = fmt.Errorf("unable to read config file")
	ErrConfigInvalidStructure error  = fmt.Errorf("config file structure is invalid")
	ErrConfigNotInitialized   error  = fmt.Errorf("config validation: config not initialized")
	ErrConfigNoLayersProvided error  = fmt.Errorf("config validation: no layers provided in config")
	ErrLayerNoName            error  = fmt.Errorf("config validation: no layer name provided")
	ErrLayerNoDataPath        error  = fmt.Errorf("config validation: no data path in layer")
	ErrLayerNoDataExtension   error  = fmt.Errorf("config validation: no data extension in layer")
	ErrLayerNoDataID          error  = fmt.Errorf("config validation: no data id in layer")
	ErrLayerNoDatabasePath    error  = fmt.Errorf("config validation: no database path in layer")
	ErrLayerNoDatabaseIndex   error  = fmt.Errorf("config validation: no database index in layer")
	ErrLayerNoServiceEndpoint error  = fmt.Errorf("config validation: no service endpoint in layer")
	ErrServiceNoPort          error  = fmt.Errorf("config validation: no service port provided")
	ErrServiceNoThrottleLimit error  = fmt.Errorf("config validation: no throttle limit provided")
	WarningLayerNoZoomLimit   string = fmt.Sprintf("warning: config validation: no zoom limit provided (default z=0) in layer")
)

// ConfigLayer is the combined layer config (data/database/service) for one data type
type ConfigLayer struct {
	Name string
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
		Endpoint  string `json:"endpoint"`
		ZoomLimit int    `json:"zoom_limit"`
	} `json:"service"`
}

// Config defines the array of data layers to create and run,
// along with the single port to run the service on
type Config struct {
	Port          int           `json:"port"`
	EnableDebug   bool          `json:"enable_debug"`
	ThrottleLimit int           `json:"throttle_limit"`
	Layers        []ConfigLayer `json:"layers"`
}

// NewConfig creates a new config from a single data layer
func NewConfig(layer ConfigLayer) Config {
	config := Config{}
	config.Layers = append(config.Layers, layer)
	return config
}

// LoadConfig creates a new config from a JSON file
func LoadConfig(path string) (Config, error) {

	if path == "" {
		return Config{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer func(f io.Closer) {
		if err := f.Close(); err != nil {
			log.Printf("%s", ErrUnableToCloseDataFile.Error())
		}
	}(file)

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, file)
	if err != nil {
		return Config{}, err
	}

	b := buf.Bytes()

	config := Config{}

	err = json.Unmarshal(b, &config)
	if err != nil {
		return Config{}, ErrConfigInvalidStructure
	}

	return config, nil
}

// ValidateConfigData checks a config to ensure
// that it is properly instantiated for data
func ValidateConfigData(cfg Config) error {

	if cfg.Layers == nil {
		return ErrConfigNoLayersProvided
	}

	// todo: add string cleaning/checking for each item below?

	for _, layer := range cfg.Layers {
		err := ValidateConfigLayerData(layer)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), layer.Name)
		}
	}

	return nil
}

// ValidateConfigLayerData checks a single config layer to ensure
// that it is properly instantiated for data
func ValidateConfigLayerData(layer ConfigLayer) error {

	if layer.Name == "" {
		return ErrLayerNoName
	}
	if layer.Data.Path == "" {
		return ErrLayerNoDataPath
	}
	if layer.Data.Extension == "" {
		return ErrLayerNoDataExtension
	}
	if layer.Data.ID == "" {
		return ErrLayerNoDataID
	}
	return nil
}

// ValidateConfigDatabase checks a config to ensure
// that it is properly instantiated for database
func ValidateConfigDatabase(cfg Config) error {

	if cfg.Layers == nil {
		return ErrConfigNoLayersProvided
	}

	// todo: add string cleaning/checking for each item below?

	for _, layer := range cfg.Layers {
		err := ValidateConfigLayerDatabase(layer)
		if err != nil {
			return fmt.Errorf("%s (%s)", err.Error(), layer.Name)
		}
	}

	return nil
}

// ValidateConfigLayerDatabase checks a single config layer to ensure
// that it is properly instantiated for database
func ValidateConfigLayerDatabase(layer ConfigLayer) error {

	//todo: if only 1 out of 3 provided (name/index/endpoint), fill in others w/ warning

	if layer.Name == "" {
		return ErrLayerNoName
	}
	if layer.Database.Path == "" {
		return ErrLayerNoDatabasePath
	}
	if layer.Database.Index == "" {
		return ErrLayerNoDatabaseIndex
	}
	return nil
}

// ValidateConfigServiceOnly checks a single config layer to ensure
// that it is properly instantiated for api service
func ValidateConfigServiceOnly(cfg Config) error {

	if cfg.Layers == nil {
		return ErrConfigNoLayersProvided
	}

	if cfg.Port == 0 {
		return ErrServiceNoPort
	}

	if cfg.ThrottleLimit == 0 {
		return ErrServiceNoThrottleLimit
	}

	//todo: if only 1 out of 3 provided (name/index/endpoint), fill in others w/ warning

	for _, layer := range cfg.Layers {
		if layer.Service.Endpoint == "" {
			return ErrLayerNoServiceEndpoint
		}
		if layer.Service.ZoomLimit == 0 {
			fmt.Println(fmt.Sprintf("%s (%s)", WarningLayerNoZoomLimit, layer.Name))
		}
	}

	return nil
}
