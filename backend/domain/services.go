package domain

import "github.com/wigfri/mustore/domain/services"

type Services interface {
	Config() services.Config
	Logger() services.Logger
}
