package handler

import (
	"encoding/json"
	"net/http"
	"sort"

	"currency-api/internal/model"
)

// CurrencyService defines the interface for currency operations
type CurrencyService interface {
	GetCurrency(code string) (model.Currency, bool)
	GetAllCurrencies() []model.Currency
	GetCurrencyCount() int
	RefreshData() error
}

// CurrencyHandler handles HTTP requests for currency data
type CurrencyHandler struct {
	service CurrencyService
}

// NewCurrencyHandler creates a new currency handler
func NewCurrencyHandler(service CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		service: service,
	}
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GetAllCurrencies handles GET /currencies
// Returns all available ISO currencies
func (h *CurrencyHandler) GetAllCurrencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currencies := h.service.GetAllCurrencies()
	
	// Sort by currency code for consistent ordering
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Code < currencies[j].Code
	})

	h.sendJSON(w, http.StatusOK, currencies)
}

// GetCurrency handles GET /currencies/{code}
// Returns a specific currency by its ISO code
func (h *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract currency code from URL path
	code := r.URL.Query().Get("code")
	if code == "" {
		// Try to extract from path if using path parameter
		code = getPathParam(r.URL.Path, "/currencies/")
	}

	if code == "" {
		h.sendError(w, "Currency code required", http.StatusBadRequest)
		return
	}

	currency, exists := h.service.GetCurrency(code)
	if !exists {
		h.sendError(w, "Currency not found", http.StatusNotFound)
		return
	}

	h.sendJSON(w, http.StatusOK, currency)
}

// Health handles GET /health
// Returns service health status
func (h *CurrencyHandler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := map[string]interface{}{
		"status":          "healthy",
		"currency_count":  h.service.GetCurrencyCount(),
		"service":         "currency-api",
		"version":         "1.0.0",
	}

	h.sendJSON(w, http.StatusOK, status)
}

// sendJSON sends a JSON response
func (h *CurrencyHandler) sendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	w.WriteHeader(statusCode)
	
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// sendError sends an error response
func (h *CurrencyHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := Response{
		Success: false,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// getPathParam extracts path parameter from URL
func getPathParam(path, prefix string) string {
	if len(path) <= len(prefix) {
		return ""
	}
	return path[len(prefix):]
}
