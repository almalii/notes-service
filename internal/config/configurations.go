package config

import (
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
		return "server.port is empty"
	}
	return port
}

func SessionKey() string {
	session := viper.GetString("SESSION_KEY")
	if session == "" {
		return "SESSION_KEY is empty"
	}
	return session
}
