package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewDbConfig() *DbConfig {
	return &DbConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}
}

type MigrationsConfig struct {
	Driver        string
	ConnString    string
	MigrationsDir string
}

func NewMigrationsConfig() *MigrationsConfig {
	return &MigrationsConfig{
		Driver:        viper.GetString("goose.driver"),
		ConnString:    viper.GetString("goose.dbstring"),
		MigrationsDir: viper.GetString("goose.dir"),
	}
}

func ServerPort() string {
	port := viper.GetString("server.port")
	if port == "" {
		logrus.Error("server.port is empty")
		return ""
	}
	return port
}

func SessionKey() string {
	session := viper.GetString("session")
	//if session == "" {
	//	logrus.Error("session key is empty")
	//	return ""
	//}
	return session
}

func SaltKey() string {
	salt := viper.GetString("salt")
	if salt == "" {
		logrus.Error("salt is empty")
		return ""
	}

	return salt
}
