package config

import (
	"fmt"

	timeutil "github.com/miiy/goc/utils/time"
	"github.com/spf13/viper"
)

type Config struct {
	App     AppConfig     `yaml:"app"`
	Redis   RedisConfig   `yaml:"redis"`
	Session SessionConfig `yaml:"session"`
	Server  ServerConfig  `yaml:"server"`
	Gateway GatewayConfig `yaml:"gateway"`
	Storage StorageConfig `yaml:"storage"`
}

type AppConfig struct {
	Name            string `yaml:"name"`
	Env             string `yaml:"env"`
	Debug           bool   `yaml:"debug"`
	Locale          string `yaml:"locale"`
	Timezone        string `yaml:"timezone"`
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
	MaxAge int    `yaml:"maxAge"`
	Secure bool   `yaml:"secure"`
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

// StorageConfig is the on-disk root nova-web serves uploaded files from
// (the files themselves are owned/written by nova-file). The path is relative
// to nova-web's working directory.
type StorageConfig struct {
	Root string `yaml:"root"`
}

var config Config

var v *viper.Viper

func NewConfig(fileName string) (*Config, error) {
	var err error
	v = viper.New()
	v.SetConfigFile(fileName)
	v.SetDefault("app.timezone", "Asia/Shanghai")
	if err = v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err = v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v \n", err)
	}
	if _, err = timeutil.LoadLocation(config.App.Timezone); err != nil {
		return nil, fmt.Errorf("invalid app timezone %q: %w", config.App.Timezone, err)
	}

	return &config, nil
}

func GetConfig(key string) any {
	return v.Get(key)
}
