package memcache

import "github.com/mrumyantsev/currency-converter-app/internal/pkg/models"

type MemCache struct {
	currencies           *models.Currencies
	updateDatetime       *models.UpdateDatetime
	calculatedCurrencies []models.CalculatedCurrency
}

func New() *MemCache {
	return new(MemCache)
}

func (m *MemCache) Currencies() *models.Currencies {
	return m.currencies
}

func (m *MemCache) SetCurrencies(currencies *models.Currencies) {
	m.currencies = currencies
}

func (m *MemCache) UpdateDatetime() *models.UpdateDatetime {
	return m.updateDatetime
}

func (m *MemCache) SetUpdateDatetime(updateDatetime *models.UpdateDatetime) {
	m.updateDatetime = updateDatetime
}

func (m *MemCache) CalculatedCurrencies() []models.CalculatedCurrency {
	return m.calculatedCurrencies
}

func (m *MemCache) SetCalculatedCurrencies(calculatedCurrencies []models.CalculatedCurrency) {
	m.calculatedCurrencies = calculatedCurrencies
}
