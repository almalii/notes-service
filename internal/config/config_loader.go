package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Loader() error {
	viper.AddConfigPath("../internal/config/")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("error reading config file")
		return nil
	}

	// Нужно ли логировать успешную установку конфигов или соедения с бд?
	logrus.Println("configs installed successfully")

	return nil
}
