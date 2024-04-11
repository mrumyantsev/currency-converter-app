package service

import (
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/repository"
)

type CurrenciesService struct {
	config     *config.Config
	repository repository.Currencies
}

func NewCurrenciesService(cfg *config.Config, repo repository.Currencies) *CurrenciesService {
	return &CurrenciesService{
		config:     cfg,
		repository: repo,
	}
}

func (s *CurrenciesService) Create(currencies models.Currencies, updateDatetimeId int) error {
	return s.repository.Create(currencies, updateDatetimeId)
}

func (s *CurrenciesService) GetLatest(updateDatetimeId int) (models.Currencies, error) {
	return s.repository.GetLatest(updateDatetimeId)
}
