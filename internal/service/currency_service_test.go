package service

import (
	"testing"

	"currency-api/internal/model"
)

// MockCurrencyRepository is a test double for CurrencyRepository
type MockCurrencyRepository struct {
	currencies map[string]model.Currency
}

func NewMockCurrencyRepository() *MockCurrencyRepository {
	return &MockCurrencyRepository{
		currencies: make(map[string]model.Currency),
	}
}

func (m *MockCurrencyRepository) GetCurrency(code string) (model.Currency, bool) {
	currency, exists := m.currencies[code]
	return currency, exists
}

func (m *MockCurrencyRepository) GetAllCurrencies() []model.Currency {
	currencies := make([]model.Currency, 0, len(m.currencies))
	for _, currency := range m.currencies {
		currencies = append(currencies, currency)
	}
	return currencies
}

func (m *MockCurrencyRepository) GetCurrencyCount() int {
	return len(m.currencies)
}

func (m *MockCurrencyRepository) RefreshData() error {
	return nil
}

func TestNewCurrencyService(t *testing.T) {
	mockRepo := NewMockCurrencyRepository()
	
	service := NewCurrencyService(mockRepo)
	
	if service == nil {
		t.Fatal("Expected service to be created")
	}
}

func TestGetCurrency(t *testing.T) {
	mockRepo := NewMockCurrencyRepository()
	mockRepo.currencies["USD"] = model.Currency{
		Code:        "USD",
		NumericCode: "840",
		Name:        "US Dollar",
		Countries:   []string{"UNITED STATES"},
		MinorUnits:  2,
	}
	
	service := NewCurrencyService(mockRepo)
	
	currency, exists := service.GetCurrency("USD")
	if !exists {
		t.Error("Expected USD to exist")
	}
	if currency.Code != "USD" {
		t.Errorf("Expected code USD, got %s", currency.Code)
	}
	if currency.Name != "US Dollar" {
		t.Errorf("Expected name 'US Dollar', got %s", currency.Name)
	}
	
	// Test non-existent currency
	_, exists = service.GetCurrency("XXX")
	if exists {
		t.Error("Expected XXX to not exist")
	}
}

func TestGetAllCurrencies(t *testing.T) {
	mockRepo := NewMockCurrencyRepository()
	mockRepo.currencies["USD"] = model.Currency{Code: "USD", Name: "US Dollar"}
	mockRepo.currencies["EUR"] = model.Currency{Code: "EUR", Name: "Euro"}
	mockRepo.currencies["GBP"] = model.Currency{Code: "GBP", Name: "British Pound"}
	
	service := NewCurrencyService(mockRepo)
	
	currencies := service.GetAllCurrencies()
	
	if len(currencies) != 3 {
		t.Errorf("Expected 3 currencies, got %d", len(currencies))
	}
	
	// Verify all currencies are present
	codes := make(map[string]bool)
	for _, c := range currencies {
		codes[c.Code] = true
	}
	
	expectedCodes := []string{"USD", "EUR", "GBP"}
	for _, code := range expectedCodes {
		if !codes[code] {
			t.Errorf("Expected currency %s to be present", code)
		}
	}
}

func TestGetCurrencyCount(t *testing.T) {
	mockRepo := NewMockCurrencyRepository()
	service := NewCurrencyService(mockRepo)
	
	if service.GetCurrencyCount() != 0 {
		t.Errorf("Expected 0 currencies initially, got %d", service.GetCurrencyCount())
	}
	
	mockRepo.currencies["USD"] = model.Currency{Code: "USD"}
	mockRepo.currencies["EUR"] = model.Currency{Code: "EUR"}
	mockRepo.currencies["JPY"] = model.Currency{Code: "JPY"}
	
	if service.GetCurrencyCount() != 3 {
		t.Errorf("Expected 3 currencies, got %d", service.GetCurrencyCount())
	}
}

func TestRefreshData(t *testing.T) {
	mockRepo := NewMockCurrencyRepository()
	service := NewCurrencyService(mockRepo)
	
	err := service.RefreshData()
	if err != nil {
		t.Errorf("Expected no error from RefreshData, got %v", err)
	}
}
