package memstorage

import "github.com/mrumyantsev/currency-converter/internal/pkg/models"

type MemStorage struct {
	currencyStorage *models.CurrencyStorage
	updateDatetime  *models.UpdateDatetime
}

func New() *MemStorage {
	return &MemStorage{}
}

func (s *MemStorage) GetCurrencyStorage() *models.CurrencyStorage {
	return s.currencyStorage
}

func (s *MemStorage) SetCurrencyStorage(currencyStorage *models.CurrencyStorage) {
	s.currencyStorage = currencyStorage
}

func (s *MemStorage) GetUpdateDatetime() *models.UpdateDatetime {
	return s.updateDatetime
}

func (s *MemStorage) SetUpdateDatetime(updateDatetime *models.UpdateDatetime) {
	s.updateDatetime = updateDatetime
}
