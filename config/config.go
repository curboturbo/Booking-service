package config

import "time"

type AppConfig struct {
	Token time.Duration `yaml:"token_ttl"`
	DataBase PostgreConfig `yaml:"db"`
	Server ServerConfig `yaml:"http"`
}

type PostgreConfig struct {
	Port     string `yaml:"port"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"app"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}