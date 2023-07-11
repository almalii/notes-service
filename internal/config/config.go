package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

type Config struct {
	DB            `yaml:"data_base"`
	HTTPServer    `yaml:"http_server"`
	MigrationsDir string `yaml:"migrations_dir" env:"MIGRATIONS_DIR"`
	Session       string `yaml:"session" env:"SESSION"`
	Salt          string `yaml:"salt" env:"SALT"`
}

type DB struct {
	Host     string `yaml:"host" env:"DB_HOST"`
	Port     string `yaml:"port" env:"DB_PORT"`
	DBName   string `yaml:"dbname" env:"DB_NAME"`
	SSLMode  string `yaml:"sslmode" env:"DB_SSLMODE"`
	UserName string `yaml:"username" env:"DB_USERNAME"`
	Password string `yaml:"password" env:"DB_PASSWORD"`
	Driver   string `yaml:"driver" env:"DB_DRIVER"`
}

type HTTPServer struct {
	Host string `yaml:"host" env:"HTTP_HOST"`
	Port string `yaml:"port" env:"HTTP_PORT"`
}

func InitConfig() *Config {
	var cfg Config

	if err := gotenv.Load("../config/.env"); err != nil {
		logrus.Fatalf("failed to load env: %v", err)
	}

	if err := cleanenv.ReadConfig("../config/config.yml", &cfg); err != nil {
		logrus.Fatalf("failed to read config file: %v", err)
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		logrus.Fatalf("failed to read env: %v", err)
	}

	return &cfg
}
