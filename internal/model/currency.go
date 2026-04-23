package model

// Currency represents an ISO 4217 currency code
type Currency struct {
	Code          string   `json:"code"`           // e.g., "USD"
	NumericCode   string   `json:"numeric_code"`   // e.g., "840"
	Name          string   `json:"name"`           // e.g., "US Dollar"
	Countries     []string `json:"countries"`      // List of country codes using this currency
	MinorUnits    int      `json:"minor_units"`    // Number of minor units (decimal places)
}

// ISO4217Response represents the parsed XML response from SIX Group
type ISO4217Response struct {
	Published string     `xml:"Pblshd,attr"`
	CurrencyTable []CurrencyEntry `xml:"CcyTbl>CcyNtry"`
}

// CurrencyEntry represents a single currency entry in the XML
type CurrencyEntry struct {
	CountryName string `xml:"CtryNm"`
	CurrencyName string `xml:"CcyNm"`
	CurrencyCode string `xml:"Ccy"`
	CurrencyNumber string `xml:"CcyNbr"`
	CurrencyMinorUnits string `xml:"CcyMnrUnts"`
}
