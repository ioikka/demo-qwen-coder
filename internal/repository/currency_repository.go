package repository

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"currency-api/internal/model"
)

const (
	// ISODataURL is the official SIX Group URL for ISO 4217 data
	ISODataURL = "https://www.six-group.com/dam/download/financial-information/data-center/iso-currrency/lists/list-one.xml"
	// DataFileName is the name of the cached JSON file
	DataFileName = "currencies.json"
	// DefaultTimeout is the default HTTP request timeout
	DefaultTimeout = 30 * time.Second
)

// CurrencyRepository handles currency data persistence and retrieval
type CurrencyRepository struct {
	dataDir      string
	dataFilePath string
	mu           sync.RWMutex
	currencies   map[string]model.Currency
	httpClient   *http.Client
}

// NewCurrencyRepository creates a new currency repository
func NewCurrencyRepository(dataDir string) (*CurrencyRepository, error) {
	repo := &CurrencyRepository{
		dataDir:    dataDir,
		httpClient: &http.Client{Timeout: DefaultTimeout},
		currencies: make(map[string]model.Currency),
	}
	
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}
	
	repo.dataFilePath = filepath.Join(dataDir, DataFileName)
	
	return repo, nil
}

// LoadData loads currency data from cache file or remote source
func (r *CurrencyRepository) LoadData() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Try to load from local cache first
	if err := r.loadFromCache(); err == nil {
		return nil
	}
	
	// Cache miss or error, fetch from remote source
	if err := r.fetchFromRemote(); err != nil {
		// If remote fetch fails, try to load from cache again as backup
		if cacheErr := r.loadFromCache(); cacheErr == nil {
			return nil
		}
		return fmt.Errorf("failed to load data from both remote and cache: %w", err)
	}
	
	return nil
}

// loadFromCache loads currency data from the local JSON cache file
func (r *CurrencyRepository) loadFromCache() error {
	file, err := os.Open(r.dataFilePath)
	if err != nil {
		return fmt.Errorf("failed to open cache file: %w", err)
	}
	defer file.Close()
	
	var currencies []model.Currency
	if err := json.NewDecoder(file).Decode(&currencies); err != nil {
		return fmt.Errorf("failed to decode cache file: %w", err)
	}
	
	for _, curr := range currencies {
		r.currencies[curr.Code] = curr
	}
	
	return nil
}

// fetchFromRemote fetches currency data from the SIX Group XML endpoint
func (r *CurrencyRepository) fetchFromRemote() error {
	resp, err := r.httpClient.Get(ISODataURL)
	if err != nil {
		return fmt.Errorf("failed to fetch remote data: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote server returned status: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	
	// Parse XML
	var isoResponse model.ISO4217Response
	if err := xml.Unmarshal(body, &isoResponse); err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}
	
	// Convert XML entries to Currency model and deduplicate by currency code
	currencyMap := make(map[string]model.Currency)
	countryMap := make(map[string][]string) // Map currency code to countries
	
	// First pass: collect all countries for each currency code
	for _, entry := range isoResponse.CurrencyTable {
		code := entry.CurrencyCode
		countryMap[code] = append(countryMap[code], entry.CountryName)
	}
	
	// Second pass: create unique currency entries with all countries
	seenCodes := make(map[string]bool)
	for _, entry := range isoResponse.CurrencyTable {
		code := entry.CurrencyCode
		
		if seenCodes[code] {
			continue
		}
		seenCodes[code] = true
		
		minorUnits, err := strconv.Atoi(entry.CurrencyMinorUnits)
		if err != nil {
			minorUnits = 2 // Default to 2 if parsing fails
		}
		
		currency := model.Currency{
			Code:        code,
			NumericCode: entry.CurrencyNumber,
			Name:        entry.CurrencyName,
			Countries:   countryMap[code],
			MinorUnits:  minorUnits,
		}
		
		currencyMap[code] = currency
		r.currencies[code] = currency
	}
	
	// Persist to cache file
	if err := r.persistToCache(currencyMap); err != nil {
		return fmt.Errorf("failed to persist to cache: %w", err)
	}
	
	return nil
}

// persistToCache saves currency data to the local JSON cache file
func (r *CurrencyRepository) persistToCache(currencies map[string]model.Currency) error {
	// Convert map to slice for JSON serialization
	currencyList := make([]model.Currency, 0, len(currencies))
	for _, curr := range currencies {
		currencyList = append(currencyList, curr)
	}
	
	// Create temporary file for atomic write
	tmpFile, err := os.CreateTemp(r.dataDir, "currencies-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	
	// Write JSON data
	encoder := json.NewEncoder(tmpFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(currencyList); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to close temp file: %w", err)
	}
	
	// Atomic rename
	if err := os.Rename(tmpPath, r.dataFilePath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename temp file: %w", err)
	}
	
	return nil
}

// GetCurrency returns a currency by its ISO code
func (r *CurrencyRepository) GetCurrency(code string) (model.Currency, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	currency, exists := r.currencies[code]
	return currency, exists
}

// GetAllCurrencies returns all available currencies
func (r *CurrencyRepository) GetAllCurrencies() []model.Currency {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	currencies := make([]model.Currency, 0, len(r.currencies))
	for _, currency := range r.currencies {
		currencies = append(currencies, currency)
	}
	return currencies
}

// GetCurrencyCount returns the total number of currencies loaded
func (r *CurrencyRepository) GetCurrencyCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.currencies)
}

// RefreshData forces a refresh of data from the remote source
func (r *CurrencyRepository) RefreshData() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Clear existing data
	r.currencies = make(map[string]model.Currency)
	
	// Fetch fresh data from remote
	if err := r.fetchFromRemote(); err != nil {
		return err
	}
	
	return nil
}
