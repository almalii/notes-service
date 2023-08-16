package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

const configPath = "../config/config.yml"

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

func InitConfig() Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		logrus.Error("error read config")
		panic(err)
	}

	return cfg
}
