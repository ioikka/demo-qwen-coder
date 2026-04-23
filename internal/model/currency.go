package model

// Currency represents an ISO 4217 currency code
type Currency struct {
	Code        string `json:"code"`         // e.g., "USD"
	NumericCode string `json:"numeric_code"` // e.g., "840"
	Name        string `json:"name"`         // e.g., "US Dollar"
	Symbol      string `json:"symbol"`       // e.g., "$"
	Countries   []string `json:"countries"`  // List of country codes using this currency
}
