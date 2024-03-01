package config

import (
	"github.com/spf13/viper"
	"os"
)

// ReadEnv Чтение переменок окружения
func ReadEnv() error {
	root, err := os.Getwd()

	if err != nil {
		return err
	}

	viper.AddConfigPath(root)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	if err != nil {
		return err
	}

	return nil
}

// InitConfig Инициализирует конфиг из переменок окружения
func InitConfig[E any](config *E) error {
	err := viper.Unmarshal(config)

	if err != nil {
		return err
	}

	return nil
}
