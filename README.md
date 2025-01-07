## Overview

A RESTful API template project for Go projects. I created this project to help bootstrap some of the boilerplate code required to create a RESTful API in Go. This template starts users with a strong separation of concerns across the project structure, electing to split up the project into the following packages:

- `cmd/api`: The executables for the API server project
- `internal/config`: Environment variable-based configuration for the project (e.g. database, server, and application configuration)
- `internal/database`: The database utilities used by the project including transaction management and connection setup
- `internal/domain`: The domain models for the project
- `internal/handler`: The HTTP handlers for the API server (essentially controllers if you come from an MVC background)
- `internal/logger`: The logger utilities for the project
- `internal/middleware`: The middleware utilities for the project
- `internal/repository`: The database repository layer for the project 
- `internal/service`: The service layer for the project with utilities for managing business logic

This project includes the `User` resource as an example. You can expand on this base or replace it with your own resources, mirroring the instances where `User` is used in the project. You will need to create domain models, handlers, repositories, and services for your additional resources.

## Dependencies

* github.com/gorilla/mux - A powerful HTTP router and URL matcher for building Go web servers with a focus on high performance and configurability
* github.com/lib/pq v1.10.9 - A pure Go Postgres driver for Go's database/sql package
* go.uber.org/zap v1.27.0 - Logging library for Go
* golang.org/x/crypto - Go's cryptography packages, used for password hashing and verification in this project
* golang.org/x/time - Go's time packages, used for rate limiting in this project

## Getting Started

1. Clone the repository
2. Install the dependencies
3. Review the `internal/config/config.go` file and update the environment variables to match your local setup. Keep in mind that you should not store sensitive information in this file. Instead, use environment variables to store sensitive information. The config file defaults are solely for local development.
4. Ensure that you have a PostgreSQL database with a `users` table as described in the code. If you are not using the default resource, you will need to create a table for your resource and update the project code to reflect the new resource.
5. Run the project with `go run cmd/api/main.go`