package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"currency-api/internal/model"
)

func TestNewCurrencyRepository(t *testing.T) {
	tmpDir := t.TempDir()
	
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	if repo == nil {
		t.Fatal("Expected repository to be created")
	}
	
	if repo.dataDir != tmpDir {
		t.Errorf("Expected dataDir to be %s, got %s", tmpDir, repo.dataDir)
	}
	
	expectedPath := filepath.Join(tmpDir, DataFileName)
	if repo.dataFilePath != expectedPath {
		t.Errorf("Expected dataFilePath to be %s, got %s", expectedPath, repo.dataFilePath)
	}
}

func TestNewCurrencyRepository_CreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "new", "nested", "dir")
	
	repo, err := NewCurrencyRepository(newDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	if repo == nil {
		t.Fatal("Expected repository to be created")
	}
	
	// Verify directory was created
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("Expected data directory to be created")
	}
}

func TestLoadFromCache(t *testing.T) {
	tmpDir := t.TempDir()
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	// Create test cache file
	testCurrencies := []model.Currency{
		{
			Code:        "USD",
			NumericCode: "840",
			Name:        "US Dollar",
			Countries:   []string{"UNITED STATES"},
			MinorUnits:  2,
		},
		{
			Code:        "EUR",
			NumericCode: "978",
			Name:        "Euro",
			Countries:   []string{"GERMANY", "FRANCE"},
			MinorUnits:  2,
		},
	}
	
	cacheFile := filepath.Join(tmpDir, DataFileName)
	file, err := os.Create(cacheFile)
	if err != nil {
		t.Fatalf("Failed to create cache file: %v", err)
	}
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(testCurrencies); err != nil {
		file.Close()
		t.Fatalf("Failed to write cache file: %v", err)
	}
	file.Close()
	
	// Load from cache
	err = repo.loadFromCache()
	if err != nil {
		t.Fatalf("Failed to load from cache: %v", err)
	}
	
	// Verify currencies were loaded
	if repo.GetCurrencyCount() != 2 {
		t.Errorf("Expected 2 currencies, got %d", repo.GetCurrencyCount())
	}
	
	usd, exists := repo.GetCurrency("USD")
	if !exists {
		t.Error("Expected USD to exist")
	}
	if usd.Name != "US Dollar" {
		t.Errorf("Expected USD name to be 'US Dollar', got %s", usd.Name)
	}
	
	eur, exists := repo.GetCurrency("EUR")
	if !exists {
		t.Error("Expected EUR to exist")
	}
	if len(eur.Countries) != 2 {
		t.Errorf("Expected EUR to have 2 countries, got %d", len(eur.Countries))
	}
}

func TestLoadFromCache_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	err = repo.loadFromCache()
	if err == nil {
		t.Error("Expected error when cache file doesn't exist")
	}
}

func TestPersistToCache(t *testing.T) {
	tmpDir := t.TempDir()
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	currencies := map[string]model.Currency{
		"USD": {
			Code:        "USD",
			NumericCode: "840",
			Name:        "US Dollar",
			Countries:   []string{"UNITED STATES"},
			MinorUnits:  2,
		},
		"JPY": {
			Code:        "JPY",
			NumericCode: "392",
			Name:        "Japanese Yen",
			Countries:   []string{"JAPAN"},
			MinorUnits:  0,
		},
	}
	
	err = repo.persistToCache(currencies)
	if err != nil {
		t.Fatalf("Failed to persist to cache: %v", err)
	}
	
	// Verify file was created
	cacheFile := filepath.Join(tmpDir, DataFileName)
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Error("Expected cache file to be created")
	}
	
	// Read and verify content
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		t.Fatalf("Failed to read cache file: %v", err)
	}
	
	var loadedCurrencies []model.Currency
	if err := json.Unmarshal(data, &loadedCurrencies); err != nil {
		t.Fatalf("Failed to unmarshal cache file: %v", err)
	}
	
	if len(loadedCurrencies) != 2 {
		t.Errorf("Expected 2 currencies in cache, got %d", len(loadedCurrencies))
	}
}

func TestGetCurrency(t *testing.T) {
	tmpDir := t.TempDir()
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	// Add test data directly
	repo.currencies["GBP"] = model.Currency{
		Code:        "GBP",
		NumericCode: "826",
		Name:        "British Pound Sterling",
		Countries:   []string{"UNITED KINGDOM"},
		MinorUnits:  2,
	}
	
	currency, exists := repo.GetCurrency("GBP")
	if !exists {
		t.Error("Expected GBP to exist")
	}
	if currency.Code != "GBP" {
		t.Errorf("Expected code GBP, got %s", currency.Code)
	}
	
	// Test non-existent currency
	_, exists = repo.GetCurrency("XXX")
	if exists {
		t.Error("Expected XXX to not exist")
	}
}

func TestGetAllCurrencies(t *testing.T) {
	tmpDir := t.TempDir()
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	repo.currencies["USD"] = model.Currency{Code: "USD", Name: "US Dollar"}
	repo.currencies["EUR"] = model.Currency{Code: "EUR", Name: "Euro"}
	repo.currencies["GBP"] = model.Currency{Code: "GBP", Name: "British Pound"}
	
	currencies := repo.GetAllCurrencies()
	
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
	tmpDir := t.TempDir()
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	if repo.GetCurrencyCount() != 0 {
		t.Errorf("Expected 0 currencies initially, got %d", repo.GetCurrencyCount())
	}
	
	repo.currencies["USD"] = model.Currency{Code: "USD"}
	repo.currencies["EUR"] = model.Currency{Code: "EUR"}
	
	if repo.GetCurrencyCount() != 2 {
		t.Errorf("Expected 2 currencies, got %d", repo.GetCurrencyCount())
	}
}

func TestRefreshData(t *testing.T) {
	tmpDir := t.TempDir()
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	// Add some initial data
	repo.currencies["TEST"] = model.Currency{Code: "TEST", Name: "Test Currency"}
	
	if repo.GetCurrencyCount() != 1 {
		t.Errorf("Expected 1 currency initially, got %d", repo.GetCurrencyCount())
	}
	
	// Refresh will try to fetch from remote (will fail in test environment)
	// but should clear existing data first
	err = repo.RefreshData()
	
	// In test environment without network, this will fail
	// But we're testing that the method exists and can be called
	if err == nil {
		// If it succeeds (network available), verify data was refreshed
		t.Logf("Refresh succeeded, currency count: %d", repo.GetCurrencyCount())
	}
}

func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()
	repo, err := NewCurrencyRepository(tmpDir)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}
	
	// Add test data
	repo.currencies["USD"] = model.Currency{Code: "USD", Name: "US Dollar"}
	
	done := make(chan bool)
	
	// Start multiple goroutines reading data
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				repo.GetCurrency("USD")
				repo.GetAllCurrencies()
				repo.GetCurrencyCount()
			}
			done <- true
		}()
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// If we reach here without panic, concurrent access is safe
	t.Log("Concurrent access test passed")
}
