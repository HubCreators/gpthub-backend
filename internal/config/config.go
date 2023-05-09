package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	HTTP     HTTPConfig
	Auth     AuthConfig
	Postgres PostgresConfig
	GIN      GINConfig
}

type HTTPConfig struct {
	Host               string        `yaml:"host"`
	Port               string        `yaml:"port"`
	ReadTimeout        time.Duration `yaml:"readTimeout"`
	WriteTimeout       time.Duration `yaml:"writeTimeout"`
	MaxHeaderMegabytes int           `yaml:"maxHeaderBytes"`
}

type PostgresConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DBName   string `yaml:"dbname"`
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration `yaml:"accessTokenTTL"`
	RefreshTokenTTL time.Duration `yaml:"refreshTokenTTL"`
	SigningKey      string
	Salt            string
}

type GINConfig struct {
	Mode string
}

func Init(configDir string) (*Config, error) {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("main")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("gpthub-backend", &cfg.Auth); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("postgres", &cfg.Postgres); err != nil {
		return nil, err
	}

	setEnvVariables(&cfg)

	logrus.Info(cfg)
	return &cfg, nil
}

func setEnvVariables(cfg *Config) {
	cfg.Auth.Salt = os.Getenv("PASSWORD_SALT")
	cfg.Auth.SigningKey = os.Getenv("JWT_SIGNING_KEY")
	cfg.Postgres.Password = os.Getenv("DB_PASSWORD")
	cfg.GIN.Mode = os.Getenv("GIN_MODE")
}
