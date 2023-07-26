package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"sync"
)

type Config struct {
	DB            DB         `yaml:"data_base"`
	HTTPServer    HTTPServer `yaml:"http_server"`
	GRPCServer    GRPCServer `yaml:"grpc_server"`
	MigrationsDir string     `yaml:"migrations_dir" env:"MIGRATIONS_DIR"`
	Session       string     `yaml:"sessions" env:"SESSION"`
	Salt          string     `yaml:"salt" env:"SALT"`
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
	Address string `yaml:"address" env:"HTTP_SERVER_ADDRESS"`
}

type GRPCServer struct {
	Address        string `yaml:"address" env:"GRPC_SERVER_ADDRESS"`
	GateWayAddress string `yaml:"gateway_address" env:"GRPC_GATEWAY_ADDRESS"`
}

const (
	FlagConfigPathName = "config"
	EnvConfigPathName  = "CONFIG_PATH"
	FlagEnvPathName    = "env"
	EnvEnvPathName     = "ENV_PATH"
)

var (
	configPath string
	envPath    string
	instance   Config
	once       sync.Once
	onceFlag   sync.Once
)

func InitConfig() Config {
	once.Do(func() {
		// этот кодв выполнится только 1 раз при первом вызове этого метода

		// 1. parse flag
		onceFlag.Do(func() {
			flag.StringVar(
				&configPath,
				FlagConfigPathName,
				"config/config.yml",
				"path to config file",
			)
			flag.StringVar(
				&envPath,
				FlagEnvPathName,
				"config/.env",
				"path to .env file",
			)

			flag.Parse()
		})

		// 2. read env
		if p, ok := os.LookupEnv(EnvConfigPathName); ok {
			configPath = p
		}

		if p, ok := os.LookupEnv(EnvEnvPathName); ok {
			envPath = p
		}

		if err := cleanenv.ReadConfig(configPath, &instance); err != nil {
			help, helpErr := cleanenv.GetDescription(&instance, nil)
			if helpErr != nil {
				logrus.Printf("error get config description due error: %v\n", helpErr)
			} else {
				logrus.Println(help)
			}

			logrus.Fatal(helpErr)
		}

		log.Println("configuration loaded")
	})

	return instance
}

func ResetOnce() {
	once = sync.Once{}
}
