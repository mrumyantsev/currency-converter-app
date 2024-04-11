package service

import (
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/repository"
)

type UpdateDatetimeService struct {
	config     *config.Config
	repository repository.UpdateDatetime
}

func NewUpdateDatetimeService(cfg *config.Config, repo repository.UpdateDatetime) *UpdateDatetimeService {
	return &UpdateDatetimeService{
		config:     cfg,
		repository: repo,
	}
}

func (s *UpdateDatetimeService) Create(datetime string) (models.UpdateDatetime, error) {
	return s.repository.Create(datetime)
}

func (s *UpdateDatetimeService) GetLatest() (models.UpdateDatetime, error) {
	return s.repository.GetLatest()
}
