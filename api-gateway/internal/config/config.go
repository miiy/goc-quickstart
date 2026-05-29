package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig             `yaml:"server"`
	Services map[string]ServiceConfig `yaml:"services"`
	TLS      TLSConfig                `yaml:"tls"`
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

// Load reads and parses the config file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
