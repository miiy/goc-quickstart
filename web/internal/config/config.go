package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig     `yaml:"app"`
	Redis   RedisConfig   `yaml:"redis"`
	Session SessionConfig `yaml:"session"`
	Server  ServerConfig  `yaml:"server"`
	Gateway GatewayConfig `yaml:"gateway"`
}

type AppConfig struct {
	Name            string `yaml:"name"`
	Env             string `yaml:"env"`
	Debug           bool   `yaml:"debug"`
	Locale          string `yaml:"locale"`
	Url             string `yaml:"url"`
	FooterCopyright string `yaml:"footerCopyright"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type SessionConfig struct {
	Name   string `yaml:"name"`
	Secret string `yaml:"secret"`
}

type ServerConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

type GatewayConfig struct {
	Addr string `yaml:"addr"`
}

var config Config

var v *viper.Viper

func NewConfig(fileName string) (*Config, error) {
	var err error
	v = viper.New()
	v.SetConfigFile(fileName)
	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v \n", err)
	}

	return &config, nil
}

func GetConfig(key string) any {
	return v.Get(key)
}
