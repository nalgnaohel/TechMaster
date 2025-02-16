package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Postgres Postgres
}

// DatabaseConfig - database configuration struct
type Postgres struct {
	Host            string
	Port            string
	User            string
	Password        string
	DbName          string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	ConnMaxIdleTime int
	ReadTimeout     int
	WriteTimeout    int
	PgDriver        string
}

// ServerConfig - server configuration struct
type ServerConfig struct {
	Port string
}

// Load config
func LoadConfig(configFileName string) (*viper.Viper, error) {
	vip := viper.New()
	fmt.Println("Loading config file: ", configFileName)
	vip.SetConfigName(configFileName)
	vip.AddConfigPath(".")
	vip.AutomaticEnv()

	err := vip.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %v !", err)
		return nil, errors.New("error reading config file")
	}

	return vip, nil
}

func ParseConfig(vip *viper.Viper) (*Config, error) {
	var config Config

	err := vip.Unmarshal(&config)
	if err != nil {
		return nil, errors.New("error parsing config file")
	}

	return &config, nil
}
