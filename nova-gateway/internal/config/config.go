package config

import (
	conf "github.com/miiy/goc/config"
)

type Config struct {
	App      AppConfig                `yaml:"app"`
	Server   ServerConfig             `yaml:"server"`
	Services map[string]ServiceConfig `yaml:"services"`
	TLS      TLSConfig                `yaml:"tls"`
}

type AppConfig struct {
	Debug bool `yaml:"debug"`
}

type ServerConfig struct {
	HTTP HTTPConfig `yaml:"http"`
}

type HTTPConfig struct {
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

// NewConfig reads and parses the config file
func NewConfig(fileName string) (*Config, error) {
	var cfg Config
	if err := conf.Load(fileName, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
