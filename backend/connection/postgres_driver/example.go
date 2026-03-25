package postgres_driver

import (
	"github.com/wigfri/mustore/domain/models"
	"gorm.io/gorm"
)

type exampleRepository struct {
	db *gorm.DB
}

type example struct {
	Id        string
	FirstRow  int
	SecondRow int
}

func (e example) model() *models.Example {
	return &models.Example{
		Id:        e.Id,
		FirstRow:  e.FirstRow,
		SecondRow: e.SecondRow,
	}
}

func makeExample(e *models.Example) example {
	return example{
		Id:        e.Id,
		FirstRow:  e.FirstRow,
		SecondRow: e.SecondRow,
	}
}

func (e *exampleRepository) Insert(example *models.Example) (string, error) {
	exampleModel := makeExample(example)

	if err := e.db.Create(exampleModel).Error; err != nil {
		return "", err
	}

	return exampleModel.Id, nil
}

func (e *exampleRepository) GetExample(id string) (*models.Example, error) {
	var result example

	if err := e.db.Where(models.Example{Id: id}).First(&result).Error; err != nil {
		return nil, err
	}

	return result.model(), nil
}

func (e *exampleRepository) All() ([]*models.Example, error) {
	var result []example

	if err := e.db.Find(&result).Error; err != nil {
		return nil, err
	}

	out := make([]*models.Example, len(result))

	for i, em := range result {
		out[i] = em.model()
	}

	return out, nil
}
