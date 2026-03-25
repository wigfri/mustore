package postgres_driver

import (
	"fmt"

	"github.com/wigfri/mustore/domain"
	"github.com/wigfri/mustore/domain/repositories"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type connection struct {
	db *gorm.DB

	exampleRepository repositories.Example
}

func makeConnection(db *gorm.DB) *connection {
	return &connection{
		db:                db,
		exampleRepository: &exampleRepository{db},
	}
}

func Make(user, password, host, port, database, sslmode string) (domain.Connection, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		database,
		sslmode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to open database due [%s]", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("unable to get DB object due [%s]", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("unable to ping DB due [%s]", err)
	}

	return makeConnection(db), nil
}

func (c connection) Example() repositories.Example {
	return c.exampleRepository
}
