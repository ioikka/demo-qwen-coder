package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"currency-api/internal/handler"
	"currency-api/internal/service"
)

func main() {
	// Initialize service with in-memory cache
	currencyService := service.NewCurrencyService()
	
	// Initialize handler
	currencyHandler := handler.NewCurrencyHandler(currencyService)

	// Setup routes
	http.HandleFunc("/health", currencyHandler.Health)
	http.HandleFunc("/currencies", currencyHandler.GetAllCurrencies)
	http.HandleFunc("/currencies/", currencyHandler.GetCurrency)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Server configuration for high load
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      loggingMiddleware(http.DefaultServeMux),
		ReadTimeout:  5 * time.Second,  // Fast failure for slow clients
		WriteTimeout: 10 * time.Second, // Reasonable time to send response
		IdleTimeout:  120 * time.Second, // Keep-alive connections
	}

	log.Printf("Starting Currency API server on port %s", port)
	log.Printf("Endpoints:")
	log.Printf("  GET /health - Service health check")
	log.Printf("  GET /currencies - List all currencies")
	log.Printf("  GET /currencies/{code} - Get specific currency")
	
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// loggingMiddleware logs incoming requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log request (in production, use structured logging)
		log.Printf("[%s] %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		
		// Enable CORS for all origins (configure appropriately for production)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
