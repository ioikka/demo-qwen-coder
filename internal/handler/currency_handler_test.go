package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"currency-api/internal/model"
)

// MockCurrencyService is a test double for service.CurrencyService
type MockCurrencyService struct {
	currencies map[string]model.Currency
}

func NewMockCurrencyService() *MockCurrencyService {
	return &MockCurrencyService{
		currencies: make(map[string]model.Currency),
	}
}

func (m *MockCurrencyService) GetCurrency(code string) (model.Currency, bool) {
	currency, exists := m.currencies[code]
	return currency, exists
}

func (m *MockCurrencyService) GetAllCurrencies() []model.Currency {
	currencies := make([]model.Currency, 0, len(m.currencies))
	for _, currency := range m.currencies {
		currencies = append(currencies, currency)
	}
	return currencies
}

func (m *MockCurrencyService) GetCurrencyCount() int {
	return len(m.currencies)
}

func (m *MockCurrencyService) RefreshData() error {
	return nil
}

func TestNewCurrencyHandler(t *testing.T) {
	mockService := NewMockCurrencyService()
	
	handler := NewCurrencyHandler(mockService)
	
	if handler == nil {
		t.Fatal("Expected handler to be created")
	}
}

func TestHealth(t *testing.T) {
	mockService := NewMockCurrencyService()
	mockService.currencies["USD"] = model.Currency{Code: "USD"}
	mockService.currencies["EUR"] = model.Currency{Code: "EUR"}
	
	handler := NewCurrencyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	
	handler.Health(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}
	
	var status map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if status["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", status["status"])
	}
	
	if int(status["currency_count"].(float64)) != 2 {
		t.Errorf("Expected currency_count 2, got %v", status["currency_count"])
	}
}

func TestHealth_MethodNotAllowed(t *testing.T) {
	mockService := NewMockCurrencyService()
	handler := NewCurrencyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodPost, "/health", nil)
	w := httptest.NewRecorder()
	
	handler.Health(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed, got %d", resp.StatusCode)
	}
}

func TestGetAllCurrencies(t *testing.T) {
	mockService := NewMockCurrencyService()
	mockService.currencies["USD"] = model.Currency{
		Code:        "USD",
		NumericCode: "840",
		Name:        "US Dollar",
		Countries:   []string{"UNITED STATES"},
		MinorUnits:  2,
	}
	mockService.currencies["EUR"] = model.Currency{
		Code:        "EUR",
		NumericCode: "978",
		Name:        "Euro",
		Countries:   []string{"GERMANY"},
		MinorUnits:  2,
	}
	
	handler := NewCurrencyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/currencies", nil)
	w := httptest.NewRecorder()
	
	handler.GetAllCurrencies(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %d", resp.StatusCode)
	}
	
	var currencies []model.Currency
	if err := json.NewDecoder(resp.Body).Decode(&currencies); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if len(currencies) != 2 {
		t.Errorf("Expected 2 currencies, got %d", len(currencies))
	}
	
	// Verify currencies are sorted by code
	if currencies[0].Code > currencies[1].Code {
		t.Error("Expected currencies to be sorted by code")
	}
}

func TestGetAllCurrencies_MethodNotAllowed(t *testing.T) {
	mockService := NewMockCurrencyService()
	handler := NewCurrencyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodPost, "/currencies", nil)
	w := httptest.NewRecorder()
	
	handler.GetAllCurrencies(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed, got %d", resp.StatusCode)
	}
}

func TestGetCurrency(t *testing.T) {
	mockService := NewMockCurrencyService()
	mockService.currencies["USD"] = model.Currency{
		Code:        "USD",
		NumericCode: "840",
		Name:        "US Dollar",
		Countries:   []string{"UNITED STATES"},
		MinorUnits:  2,
	}
	
	handler := NewCurrencyHandler(mockService)
	
	tests := []struct {
		name           string
		url            string
		expectedCode   string
		expectedStatus int
	}{
		{
			name:           "query parameter",
			url:            "/currencies?code=USD",
			expectedCode:   "USD",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "path parameter",
			url:            "/currencies/USD",
			expectedCode:   "USD",
			expectedStatus: http.StatusOK,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			
			handler.GetCurrency(w, req)
			
			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
			
			if resp.StatusCode == http.StatusOK {
				var currency model.Currency
				if err := json.NewDecoder(resp.Body).Decode(&currency); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				
				if currency.Code != tt.expectedCode {
					t.Errorf("Expected code %s, got %s", tt.expectedCode, currency.Code)
				}
			}
		})
	}
}

func TestGetCurrency_NotFound(t *testing.T) {
	mockService := NewMockCurrencyService()
	handler := NewCurrencyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/currencies?code=XXX", nil)
	w := httptest.NewRecorder()
	
	handler.GetCurrency(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status Not Found, got %d", resp.StatusCode)
	}
	
	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	if response.Success != false {
		t.Error("Expected success to be false")
	}
	if response.Error != "Currency not found" {
		t.Errorf("Expected error 'Currency not found', got %s", response.Error)
	}
}

func TestGetCurrency_MissingCode(t *testing.T) {
	mockService := NewMockCurrencyService()
	handler := NewCurrencyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/currencies", nil)
	w := httptest.NewRecorder()
	
	handler.GetCurrency(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %d", resp.StatusCode)
	}
}

func TestGetCurrency_MethodNotAllowed(t *testing.T) {
	mockService := NewMockCurrencyService()
	handler := NewCurrencyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodPost, "/currencies/USD", nil)
	w := httptest.NewRecorder()
	
	handler.GetCurrency(w, req)
	
	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status Method Not Allowed, got %d", resp.StatusCode)
	}
}

func TestResponseHeaders(t *testing.T) {
	mockService := NewMockCurrencyService()
	mockService.currencies["USD"] = model.Currency{Code: "USD"}
	handler := NewCurrencyHandler(mockService)
	
	req := httptest.NewRequest(http.MethodGet, "/currencies", nil)
	w := httptest.NewRecorder()
	
	handler.GetAllCurrencies(w, req)
	
	resp := w.Result()
	
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
	}
	
	cacheControl := resp.Header.Get("Cache-Control")
	if cacheControl != "public, max-age=3600" {
		t.Errorf("Expected Cache-Control 'public, max-age=3600', got %s", cacheControl)
	}
}
