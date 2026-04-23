package service

import (
	"sync"

	"currency-api/internal/model"
)

// CurrencyService provides currency data operations with in-memory caching
type CurrencyService struct {
	currencies map[string]model.Currency
	mu         sync.RWMutex
}

// NewCurrencyService creates a new currency service with preloaded ISO 4217 data
func NewCurrencyService() *CurrencyService {
	s := &CurrencyService{
		currencies: make(map[string]model.Currency),
	}
	s.loadISOData()
	return s
}

// loadISOData loads the official ISO 4217 currency codes
func (s *CurrencyService) loadISOData() {
	// Major currencies
	s.addCurrency(model.Currency{
		Code:        "USD",
		NumericCode: "840",
		Name:        "US Dollar",
		Symbol:      "$",
		Countries:   []string{"US", "EC", "SV", "ZW", "TL"},
	})
	s.addCurrency(model.Currency{
		Code:        "EUR",
		NumericCode: "978",
		Name:        "Euro",
		Symbol:      "€",
		Countries:   []string{"AD", "AT", "BE", "CY", "EE", "FI", "FR", "DE", "GR", "IE", "IT", "LV", "LT", "LU", "MT", "MC", "ME", "NL", "PT", "SM", "SK", "SI", "ES", "VA"},
	})
	s.addCurrency(model.Currency{
		Code:        "GBP",
		NumericCode: "826",
		Name:        "British Pound Sterling",
		Symbol:      "£",
		Countries:   []string{"GB", "IM", "GS", "IO"},
	})
	s.addCurrency(model.Currency{
		Code:        "JPY",
		NumericCode: "392",
		Name:        "Japanese Yen",
		Symbol:      "¥",
		Countries:   []string{"JP"},
	})
	s.addCurrency(model.Currency{
		Code:        "CHF",
		NumericCode: "756",
		Name:        "Swiss Franc",
		Symbol:      "Fr",
		Countries:   []string{"CH", "LI"},
	})
	s.addCurrency(model.Currency{
		Code:        "CAD",
		NumericCode: "124",
		Name:        "Canadian Dollar",
		Symbol:      "$",
		Countries:   []string{"CA"},
	})
	s.addCurrency(model.Currency{
		Code:        "AUD",
		NumericCode: "036",
		Name:        "Australian Dollar",
		Symbol:      "$",
		Countries:   []string{"AU", "CX", "CC", "HM", "KI", "NR", "NF", "TV"},
	})
	s.addCurrency(model.Currency{
		Code:        "CNY",
		NumericCode: "156",
		Name:        "Chinese Yuan",
		Symbol:      "¥",
		Countries:   []string{"CN"},
	})
	s.addCurrency(model.Currency{
		Code:        "INR",
		NumericCode: "356",
		Name:        "Indian Rupee",
		Symbol:      "₹",
		Countries:   []string{"IN", "BT", "NP"},
	})
	s.addCurrency(model.Currency{
		Code:        "BRL",
		NumericCode: "986",
		Name:        "Brazilian Real",
		Symbol:      "R$",
		Countries:   []string{"BR"},
	})
	s.addCurrency(model.Currency{
		Code:        "ZAR",
		NumericCode: "710",
		Name:        "South African Rand",
		Symbol:      "R",
		Countries:   []string{"ZA", "LS", "NA", "SZ"},
	})
	s.addCurrency(model.Currency{
		Code:        "RUB",
		NumericCode: "643",
		Name:        "Russian Ruble",
		Symbol:      "₽",
		Countries:   []string{"RU"},
	})
	s.addCurrency(model.Currency{
		Code:        "KRW",
		NumericCode: "410",
		Name:        "South Korean Won",
		Symbol:      "₩",
		Countries:   []string{"KR"},
	})
	s.addCurrency(model.Currency{
		Code:        "MXN",
		NumericCode: "484",
		Name:        "Mexican Peso",
		Symbol:      "$",
		Countries:   []string{"MX"},
	})
	s.addCurrency(model.Currency{
		Code:        "SGD",
		NumericCode: "702",
		Name:        "Singapore Dollar",
		Symbol:      "$",
		Countries:   []string{"SG"},
	})
	s.addCurrency(model.Currency{
		Code:        "HKD",
		NumericCode: "344",
		Name:        "Hong Kong Dollar",
		Symbol:      "$",
		Countries:   []string{"HK"},
	})
	s.addCurrency(model.Currency{
		Code:        "NOK",
		NumericCode: "578",
		Name:        "Norwegian Krone",
		Symbol:      "kr",
		Countries:   []string{"NO", "SJ", "BV"},
	})
	s.addCurrency(model.Currency{
		Code:        "SEK",
		NumericCode: "752",
		Name:        "Swedish Krona",
		Symbol:      "kr",
		Countries:   []string{"SE"},
	})
	s.addCurrency(model.Currency{
		Code:        "DKK",
		NumericCode: "208",
		Name:        "Danish Krone",
		Symbol:      "kr",
		Countries:   []string{"DK", "FO", "GL"},
	})
	s.addCurrency(model.Currency{
		Code:        "NZD",
		NumericCode: "554",
		Name:        "New Zealand Dollar",
		Symbol:      "$",
		Countries:   []string{"NZ", "CK", "NU", "PN", "TK"},
	})
	s.addCurrency(model.Currency{
		Code:        "TRY",
		NumericCode: "949",
		Name:        "Turkish Lira",
		Symbol:      "₺",
		Countries:   []string{"TR"},
	})
}

func (s *CurrencyService) addCurrency(currency model.Currency) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currencies[currency.Code] = currency
}

// GetCurrency returns a currency by its ISO code
func (s *CurrencyService) GetCurrency(code string) (model.Currency, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	currency, exists := s.currencies[code]
	return currency, exists
}

// GetAllCurrencies returns all available currencies
func (s *CurrencyService) GetAllCurrencies() []model.Currency {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	currencies := make([]model.Currency, 0, len(s.currencies))
	for _, currency := range s.currencies {
		currencies = append(currencies, currency)
	}
	return currencies
}

// GetCurrencyCount returns the total number of currencies loaded
func (s *CurrencyService) GetCurrencyCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.currencies)
}
