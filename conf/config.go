package conf

import (
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// var (
// 	ErrNoConfigProvided       error  = fmt.Errorf("no config provided")
// 	ErrConfigOpenFile         error  = fmt.Errorf("unable to open config file")
// 	ErrConfigReadFile         error  = fmt.Errorf("unable to read config file")
// 	ErrConfigInvalidStructure error  = fmt.Errorf("config file structure is invalid")
// 	ErrConfigNotInitialized   error  = fmt.Errorf("config validation: config not initialized")
// 	ErrConfigNoLayersProvided error  = fmt.Errorf("config validation: no layers provided in config")
// 	ErrLayerNoName            error  = fmt.Errorf("config validation: no layer name provided")
// 	ErrLayerNoDataPath        error  = fmt.Errorf("config validation: no data path in layer")
// 	ErrLayerNoDataExtension   error  = fmt.Errorf("config validation: no data extension in layer")
// 	ErrLayerNoDataID          error  = fmt.Errorf("config validation: no data id in layer")
// 	ErrLayerNoDatabasePath    error  = fmt.Errorf("config validation: no database path in layer")
// 	ErrLayerNoDatabaseIndex   error  = fmt.Errorf("config validation: no database index in layer")
// 	ErrLayerNoServiceEndpoint error  = fmt.Errorf("config validation: no service endpoint in layer")
// 	ErrServiceNoPort          error  = fmt.Errorf("config validation: no service port provided")
// 	ErrServiceNoThrottleLimit error  = fmt.Errorf("config validation: no throttle limit provided")
// 	WarningLayerNoZoomLimit   string = fmt.Sprintf("warning: config validation: no zoom limit provided (default z=0) in layer")
// )

var Configuration Config

func setDefaultConfig() {
	viper.SetDefault("Server.Host", "0.0.0.0")
	viper.SetDefault("Server.Port", 5500)
	viper.SetDefault("Server.CORSOrigins", "*")
	viper.SetDefault("Server.ReadTimeoutSec", 5)
	viper.SetDefault("Server.WriteTimeoutSec", 30)
	viper.SetDefault("Server.ThrottleLimit", 1000)
	viper.SetDefault("Server.Debug", false)
	viper.SetDefault("Server.Logs", false)
}

type Config struct {
	Server Server
	Layers []Layer
}

type Server struct {
	Host            string
	Port            int
	CORSOrigin      string
	ReadTimeoutSec  int
	WriteTimeoutSec int
	ThrottleLimit   int
	Debug           bool
	Logs            bool
}

type Layer struct {
	Name      string
	Data      LayerData
	Database  LayerDatabase
	Endpoint  string
	ZoomLimit int
}

type LayerData struct {
	Dir string
	Ext string
	ID  string
}

type LayerDatabase struct {
	Filepath string
	Index    string
}

func InitConfig(configFilename string) {

	setDefaultConfig()

	isExplicitConfigFile := configFilename != ""
	confFile := AppConfig.Name
	if configFilename != "" {
		viper.SetConfigFile(configFilename)
		confFile = configFilename
	} else {
		viper.SetConfigName(confFile)
		viper.AddConfigPath("./config")
		viper.AddConfigPath("/config")
		viper.AddConfigPath("/etc")
	}
	err := viper.ReadInConfig()
	if err != nil {
		_, isConfigFileNotFound := err.(viper.ConfigFileNotFoundError)
		errConfRead := fmt.Errorf("Fatal error reading config file: %s", err)
		isUseDefaultConfig := isConfigFileNotFound && !isExplicitConfigFile
		if isUseDefaultConfig {
			confFile = "DEFAULT"
			zap.L().Error(errConfRead.Error())
		} else {
			zap.L().Fatal(errConfRead.Error())
		}
	}

	zap.L().Info("using config file", zap.String("filename", confFile))
	err = viper.Unmarshal(&Configuration)
	if err != nil {
		// log error
	}
}

// // ValidateConfigData checks a config to ensure
// // that it is properly instantiated for data
// func ValidateConfigData(cfg Config) error {

// 	if cfg.Layers == nil {
// 		return ErrConfigNoLayersProvided
// 	}

// 	// todo: add string cleaning/checking for each item below?

// 	for _, layer := range cfg.Layers {
// 		err := ValidateConfigLayerData(layer)
// 		if err != nil {
// 			return fmt.Errorf("%s (%s)", err.Error(), layer.Name)
// 		}
// 	}

// 	return nil
// }

// // ValidateConfigLayerData checks a single config layer to ensure
// // that it is properly instantiated for data
// func ValidateConfigLayerData(layer ConfigLayer) error {

// 	if layer.Name == "" {
// 		return ErrLayerNoName
// 	}
// 	if layer.Data.Path == "" {
// 		return ErrLayerNoDataPath
// 	}
// 	if layer.Data.Extension == "" {
// 		return ErrLayerNoDataExtension
// 	}
// 	if layer.Data.ID == "" {
// 		return ErrLayerNoDataID
// 	}
// 	return nil
// }

// // ValidateConfigDatabase checks a config to ensure
// // that it is properly instantiated for database
// func ValidateConfigDatabase(cfg Config) error {

// 	if cfg.Layers == nil {
// 		return ErrConfigNoLayersProvided
// 	}

// 	// todo: add string cleaning/checking for each item below?

// 	for _, layer := range cfg.Layers {
// 		err := ValidateConfigLayerDatabase(layer)
// 		if err != nil {
// 			return fmt.Errorf("%s (%s)", err.Error(), layer.Name)
// 		}
// 	}

// 	return nil
// }

// // ValidateConfigLayerDatabase checks a single config layer to ensure
// // that it is properly instantiated for database
// func ValidateConfigLayerDatabase(layer ConfigLayer) error {

// 	//todo: if only 1 out of 3 provided (name/index/endpoint), fill in others w/ warning

// 	if layer.Name == "" {
// 		return ErrLayerNoName
// 	}
// 	if layer.Database.Path == "" {
// 		return ErrLayerNoDatabasePath
// 	}
// 	if layer.Database.Index == "" {
// 		return ErrLayerNoDatabaseIndex
// 	}
// 	return nil
// }

// // ValidateConfigServiceOnly checks a single config layer to ensure
// // that it is properly instantiated for api service
// func ValidateConfigServiceOnly(cfg Config) error {

// 	if cfg.Layers == nil {
// 		return ErrConfigNoLayersProvided
// 	}

// 	if cfg.Port == 0 {
// 		return ErrServiceNoPort
// 	}

// 	if cfg.ThrottleLimit == 0 {
// 		return ErrServiceNoThrottleLimit
// 	}

// 	//todo: if only 1 out of 3 provided (name/index/endpoint), fill in others w/ warning

// 	for _, layer := range cfg.Layers {
// 		if layer.Service.Endpoint == "" {
// 			return ErrLayerNoServiceEndpoint
// 		}
// 		if layer.Service.ZoomLimit == 0 {
// 			fmt.Println(fmt.Sprintf("%s (%s)", WarningLayerNoZoomLimit, layer.Name))
// 		}
// 	}

// 	return nil
// }
