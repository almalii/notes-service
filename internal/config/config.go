package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type Config struct {
	DB            DB            `yaml:"data_base"`
	HTTPServer    HTTPServer    `yaml:"http_server"`
	GRPCServer    GRPCServer    `yaml:"grpc_server"`
	GatewayServer GatewayServer `yaml:"grpc_gateway"`
	Redis         Redis         `yaml:"redis"`
	MigrationsDir string        `yaml:"migrations_dir" env:"MIGRATIONS_DIR"`
	JwtSigning    string        `yaml:"jwt_signing" env:"JWT_SIGNING"`
	SaltHash      string        `yaml:"salt_hash" env:"SALT_HASH"`
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

type Redis struct {
	Address  string `yaml:"address" env:"REDIS_ADDRESS"`
	Password string `yaml:"password" env:"REDIS_PASSWORD"`
	DB       int    `yaml:"db" env:"REDIS_DB"`
}

type HTTPServer struct {
	Address string `yaml:"address" env:"HTTP_SERVER_ADDRESS"`
	//ReadTimeout    time.Duration `yaml:"read_timeout" env:"HTTP_SERVER_READ_TIME_OUT"`
	//WriteTimeout   time.Duration `yaml:"write_timeout" env:"HTTP_SERVER_WRITE_TIME_OUT"`
	//MaxHeaderBytes int           `yaml:"max_header_bytes" env:"HTTP_SERVER_MAX_HEADER"`
}

type GRPCServer struct {
	Address string `yaml:"address" env:"GRPC_SERVER_ADDRESS"`
}

type GatewayServer struct {
	Address string `yaml:"address" env:"GATEWAY_SERVER_ADDRESS"`
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

		logrus.Println("configuration loaded")
	})

	return instance
}

func ResetOnce() {
	once = sync.Once{}
}
