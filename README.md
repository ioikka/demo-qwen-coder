# Currency API

A high-performance REST API service for ISO 4217 currency reference data.

## Tech Stack

- **Language**: Go (Golang) - chosen for excellent concurrency support, low memory footprint, and high performance under load
- **Architecture**: Clean architecture with separation of concerns (handler, service, model layers)
- **Caching**: In-memory caching with thread-safe access using `sync.RWMutex`
- **HTTP Server**: Standard library `net/http` with optimized timeouts for high-load scenarios

## Features

- ✅ Preloaded ISO 4217 currency codes (21 major currencies)
- ✅ Thread-safe in-memory storage
- ✅ RESTful endpoints with proper HTTP status codes
- ✅ CORS support
- ✅ Request logging middleware
- ✅ Health check endpoint
- ✅ Optimized server configuration for high load
- ✅ JSON responses with consistent format

## API Endpoints

### GET /health
Returns service health status and metadata.

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "currency_count": 21,
  "service": "currency-api",
  "status": "healthy",
  "version": "1.0.0"
}
```

### GET /currencies
Returns all available ISO currencies.

```bash
curl http://localhost:8080/currencies
```

Response:
```json
[
  {
    "code": "USD",
    "numeric_code": "840",
    "name": "US Dollar",
    "symbol": "$",
    "countries": ["US", "EC", "SV", "ZW", "TL"]
  },
  ...
]
```

### GET /currencies/{code}
Returns a specific currency by its ISO code.

```bash
curl http://localhost:8080/currencies/EUR
```

Response:
```json
{
  "code": "EUR",
  "numeric_code": "978",
  "name": "Euro",
  "symbol": "€",
  "countries": ["AD", "AT", "BE", "CY", "EE", "FI", "FR", "DE", "GR", "IE", "IT", "LV", "LT", "LU", "MT", "MC", "ME", "NL", "PT", "SM", "SK", "SI", "ES", "VA"]
}
```

Error response (404):
```json
{
  "success": false,
  "error": "Currency not found"
}
```

## Project Structure

```
currency-api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── handler/
│   │   └── currency_handler.go  # HTTP request handlers
│   ├── model/
│   │   └── currency.go      # Data models
│   └── service/
│       └── currency_service.go  # Business logic and caching
├── pkg/
│   └── cache/               # Reserved for future cache implementations
├── bin/
│   └── api                  # Compiled binary
├── go.mod
└── README.md
```

## Building and Running

### Prerequisites
- Go 1.21 or later

### Build
```bash
cd /workspace
go mod tidy
go build -o bin/api ./cmd/api
```

### Run
```bash
./bin/api
```

The server starts on port 8080 by default. Configure via environment variable:
```bash
PORT=3000 ./bin/api
```

## Performance Considerations

This implementation is optimized for high-load scenarios:

1. **In-Memory Storage**: All currency data is loaded into memory at startup for O(1) lookups
2. **Read-Write Mutex**: Allows concurrent reads while ensuring thread-safe writes
3. **Connection Timeouts**: Configured to prevent resource exhaustion from slow clients
4. **Keep-Alive**: Idle timeout allows connection reuse
5. **No External Dependencies**: Minimal overhead from third-party libraries
6. **Cache Headers**: Responses include Cache-Control headers for client-side caching

## Future Enhancements

- Add Redis/Memcached integration for distributed caching
- Implement rate limiting
- Add metrics endpoint (Prometheus)
- Support for currency updates/refresh from external sources
- GraphQL endpoint alternative
- OpenAPI/Swagger documentation

## License

MIT
