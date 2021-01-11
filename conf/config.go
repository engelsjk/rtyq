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

	fmt.Println(configFilename)
	
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
		fmt.Println(err.Error())
	}
}
