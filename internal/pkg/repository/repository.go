package repository

import (
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/database"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/repository/postgres"
)

type UpdateDatetime interface {
	Create(datetime string) (models.UpdateDatetime, error)
	GetLatest() (models.UpdateDatetime, error)
}

type Currencies interface {
	Create(currencies models.Currencies, updateDatetimeId int) error
	GetLatest(updateDatetimeId int) (models.Currencies, error)
}

type Repository struct {
	UpdateDatetime UpdateDatetime
	Currencies     Currencies
}

func New(cfg *config.Config, db *database.Database) *Repository {
	return &Repository{
		UpdateDatetime: postgres.NewUpdateDatetimeRepository(cfg, db),
		Currencies:     postgres.NewCurrenciesRepository(cfg, db),
	}
}
