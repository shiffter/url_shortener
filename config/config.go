package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server   Server
	Postgres Postgres
}

type Server struct {
	Host        string
	Port        string
	StorageMode string
}

type Postgres struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Pass     string `json:"pass"`
	DbName   string `json:"dbName"`
	PgDriver string `json:"pgDriver"`
	SslMode  string `json:"sslMode"`
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()

	v.AddConfigPath("./config")
	v.SetConfigName("config")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatalf("cant decode config into struct, %v", err)
		return nil, err
	}
	return &c, nil
}
