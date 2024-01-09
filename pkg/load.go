package pkg

import (
	"fmt"

	"github.com/spf13/viper"
)

type config struct {
	Host      string `yaml:"Host" mapstructure:"Host" json:"Host"`
	Port      string `yaml:"Port" mapstructure:"Port" json:"Port"`
	SecretKey string `yaml:"SecretKey" mapstructure:"SecretKey" json:"SecretKey"`
	ResultCLi string `yaml:"ResultCLi" mapstructure:"ResultCLi" json:"ResultCLi"`
}

var Config config

func Load(configPath string) error {
	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		err = fmt.Errorf("error reading config file, %s", err)
		return err
	}

	Config.Host = viper.GetString("Host")
	Config.Port = viper.GetString("Port")
	Config.SecretKey = viper.GetString("SecretKey")
	Config.ResultCLi = viper.GetString("ResultCLi")

	return nil
}

func GetConfig() config {
	return Config
}

func IsValidConfig() bool {
	return Config.SecretKey != ""
}
