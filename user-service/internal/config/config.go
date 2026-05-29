package config

import (
	conf "github.com/miiy/goc/config"
)

type Config struct {
	App        App        `yaml:"app"`
	Database   Database   `yaml:"database"`
	Redis      Redis      `yaml:"redis"`
	Server     Server     `yaml:"server"`
	GrpcClient GrpcClient `yaml:"grpcClient"`
	Jwt        Jwt        `yaml:"jwt"`
	Snowflake  Snowflake  `yaml:"snowflake"`
}

type App struct {
	Name    string `yaml:"name"`
	Env     string `yaml:"env"`
	Debug   bool   `yaml:"debug"`
	Version string `yaml:"version"`
}

type Database struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Redis struct {
	Addrs    []string `yaml:"addrs"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	DB       int      `yaml:"db"`
}

type Jwt struct {
	Secret    string `yaml:"secret"`
	Issuer    string `yaml:"issuer"`
	ExpiresIn int64  `yaml:"expiresIn"`
}

type Server struct {
	Http ServerHttp `yaml:"http"`
	Grpc ServerGrpc `yaml:"grpc"`
}

type ServerHttp struct {
	Addr string `yaml:"addr"`
	Url  string `yaml:"url"`
}

type ServerGrpc struct {
	Addr string        `yaml:"addr"`
	Tls  ServerGrpcTLS `yaml:"tls"`
}

type ServerGrpcTLS struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
	CaFile   string `yaml:"caFile"`
}

type GrpcClient struct {
	Endpoint string        `yaml:"endpoint"`
	Tls      GrpcClientTLS `yaml:"tls"`
}

type GrpcClientTLS struct {
	ServerName string `yaml:"serverName"`
	CertFile   string `yaml:"certFile"`
	KeyFile    string `yaml:"keyFile"`
	CaFile     string `yaml:"caFile"`
}

type Snowflake struct {
	Node int64 `yaml:"node"`
}

var config *Config

func NewConfig(fileName string) (*Config, error) {
	if err := conf.Load(fileName, &config); err != nil {
		return nil, err
	}
	return config, nil
}
