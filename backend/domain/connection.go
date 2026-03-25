package domain

import "github.com/wigfri/mustore/domain/repositories"

type Connection interface {
	Example() repositories.Example
}
