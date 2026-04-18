package config

import (
	"log/slog"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Host string `yaml:"host" env-default:"0.0.0.0" env-required:"true"`
	Port string `yaml:"port" env-default:"8000" env-required:"true"`
}

type Db struct {
	Host     string `yaml:"host" env-default:"example"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"example"`
	Password string `yaml:"password" env-default:"example"`
	Name     string `yaml:"db_name" env-default:"example"`
	SslMode  string `yaml:"sslmode"`
}

type Auth struct {
	JwtSecret     string `yaml:"jwt_secret" env-default:"example"`
	JwtTTLMinutes int    `yaml:"jwt_ttl_minutes" env-default:"60"`
}

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	HttpServer HttpServer `yaml:"http_server"`
	Db         Db         `yaml:"db"`
	Auth       Auth       `yaml:"auth"`
}

func Make() *Config {

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist " + path)
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		slog.Error("cannot read config file", err)
	}

	return &config
}

func (s *Config) HttpHost() string {
	return s.HttpServer.Host
}

func (s *Config) HttpPort() string {
	return s.HttpServer.Port
}

func (s *Config) PostgresHost() string {
	return s.Db.Host
}

func (s *Config) PostgresPort() string {
	return s.Db.Port
}

func (s *Config) PostgresUser() string {
	return s.Db.User
}

func (s *Config) PostgresPassword() string {
	return s.Db.Password
}

func (s *Config) PostgresName() string {
	return s.Db.Name
}

func (s *Config) EnvLevel() string {
	return s.Env
}

func (s *Config) SslMode() string {
	return s.Db.SslMode
}

func (s *Config) JwtSecret() string {
	return s.Auth.JwtSecret
}

func (s *Config) JwtTTL() time.Duration {
	if s.Auth.JwtTTLMinutes <= 0 {
		return time.Hour
	}
	return time.Duration(s.Auth.JwtTTLMinutes) * time.Minute
}
