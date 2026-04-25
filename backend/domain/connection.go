package domain

import "github.com/wigfri/mustore/domain/repositories"

type Connection interface {
	User() repositories.User
	Instrument() repositories.Instrument
}
