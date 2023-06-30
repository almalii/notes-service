package config

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func Loader() error {

	//go run file_name.go -configDir=/path/to/config/dir
	configDir := flag.String("configDir", "../internal/config/", "Directory path for configuration files")
	flag.Parse()

	if *configDir == "" {
		logrus.Fatal("config directory is not specified")
		return nil
	}

	_, err := os.Stat(*configDir)
	if os.IsNotExist(err) {
		logrus.Fatalf("config directory does not exist: %s", *configDir)
		return err
	}

	viper.AddConfigPath("../internal/config/")
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("error reading config file")
		return nil
	}

	logrus.Println("configs installed successfully")

	return nil
}
