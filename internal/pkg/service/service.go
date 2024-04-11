package service

import (
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/repository"
)

type UpdateDatetime interface {
	Create(datetime string) (models.UpdateDatetime, error)
	GetLatest() (models.UpdateDatetime, error)
}

type Currencies interface {
	Create(currencies models.Currencies, updateDatetimeId int) error
	GetLatest(updateDatetimeId int) (models.Currencies, error)
}

type Service struct {
	UpdateDatetime UpdateDatetime
	Currencies     Currencies
}

func New(cfg *config.Config, repo *repository.Repository) *Service {
	return &Service{
		UpdateDatetime: NewUpdateDatetimeService(cfg, repo.UpdateDatetime),
		Currencies:     NewCurrenciesService(cfg, repo.Currencies),
	}
}
