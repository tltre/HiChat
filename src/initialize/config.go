package initialize

import (
	"HiChat/src/global"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig(configType string) {
	// Get Config File Path By Type
	var configPath string
	if configType == "debug" {
		configPath = "./src/config-debug.yaml"
	} else if configType == "deploy" {
		configPath = "./src/config-deploy.yaml"
	}

	// read Config File
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(&global.ServiceConfig); err != nil {
		panic(err)
	}

	zap.S().Info("Config Information: ", global.ServiceConfig)
}
