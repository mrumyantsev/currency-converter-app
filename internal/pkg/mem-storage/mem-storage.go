package memstorage

import "github.com/mrumyantsev/currency-converter-app/internal/pkg/models"

type MemStorage struct {
	currencies           *models.Currencies
	updateDatetime       *models.UpdateDatetime
	calculatedCurrencies []models.CalculatedCurrency
}

func New() *MemStorage {
	return &MemStorage{}
}

func (s *MemStorage) Currencies() *models.Currencies {
	return s.currencies
}

func (s *MemStorage) SetCurrencies(currencies *models.Currencies) {
	s.currencies = currencies
}

func (s *MemStorage) UpdateDatetime() *models.UpdateDatetime {
	return s.updateDatetime
}

func (s *MemStorage) SetUpdateDatetime(updateDatetime *models.UpdateDatetime) {
	s.updateDatetime = updateDatetime
}

func (s *MemStorage) CalculatedCurrencies() []models.CalculatedCurrency {
	return s.calculatedCurrencies
}

func (s *MemStorage) SetCalculatedCurrencies(calculatedCurrencies []models.CalculatedCurrency) {
	s.calculatedCurrencies = calculatedCurrencies
}
