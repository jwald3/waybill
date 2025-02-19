## Overview

A RESTful API template project for a fleet management system built in Go. This API provides endpoints for managing:

- Drivers: Track driver information, licensing, and employment status with state management (Active, Suspended, Terminated)
- Trucks: Manage fleet vehicles including status tracking (Available, In Transit, Under Maintenance, Retired), maintenance history, and mileage logs
- Facilities: Handle locations for loading, unloading, and fleet services with configurable service availability
- Trips: Coordinate shipments with full lifecycle management (Scheduled, In Transit, Completed, Failed Delivery, Canceled)
- Maintenance Logs: Record vehicle maintenance activities and repairs
- Fuel Logs: Track fuel consumption and costs
- Incident Reports: Document accidents, mechanical failures, and other incidents

The project structure is organized into the following packages:

- `cmd/api`: The main application executable and server initialization
- `internal/config`: Configuration management using environment variables
- `internal/database`: MongoDB connection and transaction management
- `internal/domain`: Domain models and business logic interfaces
- `internal/handler`: HTTP request handlers and routing logic
- `internal/logger`: Logging configuration and utilities
- `internal/middleware`: HTTP middleware components
- `internal/repository`: Data access layer for MongoDB operations
- `internal/service`: Business logic implementation layer

## Features

- State machine implementation for managing resource status transitions
- MongoDB integration with aggregation pipelines for related data
- Pagination and filtering for list endpoints
- CORS and logging middleware
- Structured error handling
- Environment-based configuration

## Dependencies

* [gorilla/mux](https://github.com/gorilla/mux) - Powerful HTTP router and URL matcher
* [mongo-driver](https://github.com/mongodb/mongo-go-driver) - Official MongoDB driver for Go
* [zap](https://github.com/uber-go/zap) - Fast, structured logging
* [lollipop](https://github.com/jwald3/lollipop) - State machine implementation
* [golang.org/x/time](https://golang.org/x/time) - Rate limiting implementation

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Set up your environment variables:
   - Copy `.env.example` to `.env` (if using environment file)
   - Configure your MongoDB connection and application settings
   - See `internal/config/config.go` for available configuration options

4. Set up MongoDB:
   - Create a MongoDB database
   - The application will create collections as needed

5. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

## Project Structure
```
├── cmd/
│   └── api/            # Application entrypoint
├── internal/
│   ├── config/         # Configuration management
│   ├── database/       # Database utilities
│   ├── domain/         # Domain models and interfaces
│   ├── handler/        # HTTP handlers
│   ├── logger/         # Logging setup
│   ├── middleware/     # HTTP middleware
│   ├── repository/     # Data access layer
│   └── service/        # Business logic layer
└── README.md
```

## Development

- The project follows standard Go project layout and coding conventions
- Each resource (Driver, Truck, etc.) has its own set of models, handlers, and services
- State transitions are managed through a state machine pattern
- MongoDB aggregation pipelines are used for efficient data retrieval

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.