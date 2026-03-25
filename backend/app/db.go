package app

import (
	"github.com/wigfri/mustore/connection/postgres_driver"
	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/services/config"
)

func InitDB(cfg *config.Config) (domain.Connection, error) {
	return postgres_driver.Make(cfg.PostgresUser(), cfg.PostgresPassword(), cfg.PostgresHost(), cfg.PostgresPort(), cfg.PostgresName(), cfg.SslMode())
}
