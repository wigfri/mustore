package app

import (
	"context"
	"log"

	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/services"
	"github.com/wigfri/mustore/services/config"
	"github.com/wigfri/mustore/services/logger"
	"github.com/wigfri/mustore/services/mailqueue"
	"github.com/wigfri/mustore/services/redisotp"
)

type ctx struct {
	services   domain.Services
	connection domain.Connection
}

type svs struct {
	config   services.Config
	logger   services.Logger
	mail     services.MailQueue
	otp      services.OTPStore
}

func (s *svs) Logger() services.Logger {
	return s.logger
}

func (s *svs) Config() services.Config {
	return s.config
}

func (s *svs) MailQueue() services.MailQueue {
	return s.mail
}

func (s *svs) OTPStore() services.OTPStore {
	return s.otp
}

func (c *ctx) Services() domain.Services {
	return c.services
}

func (c *ctx) Connection() domain.Connection {
	return c.connection
}

func (c *ctx) Make() domain.Context {
	return &ctx{
		services:   c.services,
		connection: c.connection,
	}
}

func InitCtx(cfg *config.Config) *ctx {
	connection, err := InitDB(cfg)
	if err != nil {
		log.Fatalf("cant initialize connection context due [%s]", err)
	}

	logSvc := logger.Init(cfg.EnvLevel())
	mailDisp := mailqueue.NewDispatcher(cfg, logSvc)
	mailDisp.StartConsumer()

	otpStore := redisotp.New(cfg.RedisAddr(), cfg.RedisPassword(), cfg.RedisDB())
	if err := otpStore.Ping(context.Background()); err != nil {
		log.Fatalf("redis otp store: %v", err)
	}

	return &ctx{
		services: &svs{
			config: cfg,
			logger: logSvc,
			mail:   mailDisp,
			otp:    otpStore,
		},
		connection: connection,
	}
}
