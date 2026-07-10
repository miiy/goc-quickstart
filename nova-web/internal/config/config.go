package config

import (
	"fmt"

	timeutil "github.com/miiy/goc/utils/time"
	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig                `yaml:"app"`
	Redis    RedisConfig              `yaml:"redis"`
	Session  SessionConfig            `yaml:"session"`
	Server   ServerConfig             `yaml:"server"`
	Static   StaticConfig             `yaml:"static"`
	Gateway  GatewayConfig            `yaml:"gateway"`
	Services map[string]ServiceConfig `yaml:"services"`
	TLS      TLSConfig                `yaml:"tls"`
	Storage  StorageConfig            `yaml:"storage"`
}

type AppConfig struct {
	Name            string `yaml:"name"`
	Description     string `yaml:"description"`
	Env             string `yaml:"env"`
	Debug           bool   `yaml:"debug"`
	RegisterEnabled bool   `yaml:"registerEnabled"`
	Locale          string `yaml:"locale"`
	Timezone        string `yaml:"timezone"`
	Url             string `yaml:"url"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type SessionConfig struct {
	Name   string `yaml:"name"`
	Secret string `yaml:"secret"`
	Domain string `yaml:"domain"`
	MaxAge int    `yaml:"maxAge"`
	Secure bool   `yaml:"secure"`
}

type ServerConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"`
}

// StaticConfig controls where nova-web reads frontend build artifacts from.
type StaticConfig struct {
	Root string `yaml:"root"`
}

type GatewayConfig struct {
	Addr string `yaml:"addr"`
}

type ServiceConfig struct {
	Endpoint string `yaml:"endpoint"`
}

type TLSConfig struct {
	Enabled    bool   `yaml:"enabled"`
	ServerName string `yaml:"serverName"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	CaFile     string `yaml:"caFile"`
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
	v.SetDefault("app.registerEnabled", true)
	v.SetDefault("app.timezone", "Asia/Shanghai")
	v.SetDefault("static.root", "dist")
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
