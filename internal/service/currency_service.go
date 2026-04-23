package service

import (
	"currency-api/internal/model"
)

// CurrencyRepository defines the interface for currency data access
type CurrencyRepository interface {
	GetCurrency(code string) (model.Currency, bool)
	GetAllCurrencies() []model.Currency
	GetCurrencyCount() int
	RefreshData() error
}

// CurrencyService provides currency data operations with in-memory caching
type CurrencyService struct {
	repo CurrencyRepository
}

// NewCurrencyService creates a new currency service with data from repository
func NewCurrencyService(repo CurrencyRepository) *CurrencyService {
	return &CurrencyService{
		repo: repo,
	}
}

// GetCurrency returns a currency by its ISO code
func (s *CurrencyService) GetCurrency(code string) (model.Currency, bool) {
	return s.repo.GetCurrency(code)
}

// GetAllCurrencies returns all available currencies
func (s *CurrencyService) GetAllCurrencies() []model.Currency {
	return s.repo.GetAllCurrencies()
}

// GetCurrencyCount returns the total number of currencies loaded
func (s *CurrencyService) GetCurrencyCount() int {
	return s.repo.GetCurrencyCount()
}

// RefreshData forces a refresh of data from the remote source
func (s *CurrencyService) RefreshData() error {
	return s.repo.RefreshData()
}
