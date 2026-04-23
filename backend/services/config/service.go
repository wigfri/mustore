package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HttpServer struct {
	Host               string `yaml:"host" env-default:"0.0.0.0" env-required:"true"`
	Port               string `yaml:"port" env-default:"8000" env-required:"true"`
	CorsAllowedOrigins string `yaml:"cors_allowed_origins"`
}

type Db struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"db_name"`
	SslMode  string `yaml:"sslmode"`
}

type Auth struct {
	JwtSecret     string `yaml:"jwt_secret"`
	JwtTTLMinutes int    `yaml:"jwt_ttl_minutes" env-default:"60"`
}

type RabbitMQ struct {
	Enabled bool   `yaml:"enabled"`
	URL     string `yaml:"url"`
}

type Mail struct {
	SkipSend     bool   `yaml:"skip_send"`
	SMTPHost     string `yaml:"smtp_host"`
	SMTPPort     int    `yaml:"smtp_port"`
	SMTPUser     string `yaml:"smtp_user"`
	SMTPPassword string `yaml:"smtp_password"`
	FromAddress  string `yaml:"from_address"`
}

type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	HttpServer HttpServer `yaml:"http_server"`
	Db         Db         `yaml:"db"`
	Auth       Auth       `yaml:"auth"`
	RabbitMQ   RabbitMQ   `yaml:"rabbitmq"`
	Mail       Mail       `yaml:"mail"`
	Redis      Redis      `yaml:"redis"`
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

	fmt.Println(config)

	return &config
}

func (s *Config) HttpHost() string {
	return s.HttpServer.Host
}

func (s *Config) HttpPort() string {
	return s.HttpServer.Port
}

// CorsAllowedOrigins returns comma-separated origins for CORS (required when the browser sends credentials).
func (s *Config) CorsAllowedOrigins() string {
	o := strings.TrimSpace(s.HttpServer.CorsAllowedOrigins)
	if o == "" {
		return "http://localhost:5173"
	}
	return o
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

func (s *Config) RabbitMQEnabled() bool {
	return s.RabbitMQ.Enabled
}

func (s *Config) RabbitMQURL() string {
	return s.RabbitMQ.URL
}

func (s *Config) MailSkipSend() bool {
	return s.Mail.SkipSend
}

func (s *Config) SMTPHost() string {
	return s.Mail.SMTPHost
}

func (s *Config) SMTPPort() int {
	if s.Mail.SMTPPort <= 0 {
		return 587
	}
	return s.Mail.SMTPPort
}

func (s *Config) SMTPUser() string {
	return s.Mail.SMTPUser
}

func (s *Config) SMTPPassword() string {
	return s.Mail.SMTPPassword
}

func (s *Config) MailFrom() string {
	return s.Mail.FromAddress
}

func (s *Config) RedisAddr() string {
	return s.Redis.Addr
}

func (s *Config) RedisPassword() string {
	return s.Redis.Password
}

func (s *Config) RedisDB() int {
	return s.Redis.DB
}
