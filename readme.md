# ![RealWorld Example App](logo.png)


[![CI](https://github.com/gothinkster/golang-gin-realworld-example-app/actions/workflows/ci.yml/badge.svg)](https://github.com/gothinkster/golang-gin-realworld-example-app/actions/workflows/ci.yml)
[![Coverage Status](https://coveralls.io/repos/github/gothinkster/golang-gin-realworld-example-app/badge.svg?branch=main)](https://coveralls.io/github/gothinkster/golang-gin-realworld-example-app?branch=main)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/gothinkster/golang-gin-realworld-example-app/blob/main/LICENSE)
[![GoDoc](https://godoc.org/github.com/gothinkster/golang-gin-realworld-example-app?status.svg)](https://godoc.org/github.com/gothinkster/golang-gin-realworld-example-app)

> ### Golang/Gin codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.


This codebase was created to demonstrate a fully fledged fullstack application built with **Golang/Gin** including CRUD operations, authentication, routing, pagination, and more.

## Recent Updates

This project has been modernized with the following updates:
- **Go 1.27+**: Updated to require Go 1.27 or higher
- **GORM v2**: Migrated from deprecated jinzhu/gorm v1 to gorm.io/gorm v2
- **JWT v5**: Updated from deprecated dgrijalva/jwt-go to golang-jwt/jwt/v5 (fixes CVE-2020-26160)
- **Validator v10**: Updated validator tags and package to match gin v1.10.0
- **Latest Dependencies**: All dependencies updated to their 2025 production-stable versions
- **RealWorld API Spec Compliance**:
  - `GET /profiles/:username` now supports optional authentication (anonymous access allowed)
  - `POST /users/login` returns 401 Unauthorized on failure (previously 403)
  - `GET /articles/feed` registered as dedicated authenticated route
  - `DELETE /articles/:slug` and `DELETE /articles/:slug/comments/:id` return empty response body

## Test Coverage

The project maintains high test coverage across all core packages:

| Package | Coverage |
|---------|----------|
| `articles` | 93.4% |
| `users` | 99.5% |
| `common` | 85.7% |
| **Total** | **90.0%** |

To generate a coverage report locally, run:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Dependencies (2025 Stable Versions)

| Package | Version | Release Date | Known Issues |
|---------|---------|--------------|--------------|
| [gin-gonic/gin](https://github.com/gin-gonic/gin) | v1.10.0 | 2024-05 | None; v1.11 has experimental HTTP/3 support |
| [gorm.io/gorm](https://gorm.io/) | v1.25.12 | 2024-08 | None; v1.30+ has breaking changes |
| [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) | v5.2.1 | 2024-06 | None; v5.3 only bumps Go version requirement |
| [go-playground/validator/v10](https://github.com/go-playground/validator) | v10.24.0 | 2024-12 | None; v10.30+ requires Go 1.24 |
| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | v0.32.0 | 2025-01 | None; keep updated for security fixes |
| [gorm.io/driver/sqlite](https://github.com/go-gorm/sqlite) | v1.5.7 | 2024-09 | None; requires cgo; use glebarez/sqlite for pure Go |
| [gosimple/slug](https://github.com/gosimple/slug) | v1.15.0 | 2024-12 | None |
| [stretchr/testify](https://github.com/stretchr/testify) | v1.10.0 | 2024-10 | None; v2 still in development |


# Directory structure
```
.
├── gorm.db
├── hello.go
├── common
│   ├── utils.go        //small tools function
│   └── database.go     //DB connect manager
├── users
|   ├── models.go       //data models define & DB operation
|   ├── serializers.go  //response computing & format
|   ├── routers.go      //business logic & router binding
|   ├── middlewares.go  //put the before & after logic of handle request
|   └── validators.go   //form/json checker
├── ...
...
```

# Getting started

## Install Golang

Make sure you have Go 1.27 or higher installed.

https://golang.org/doc/install

## Environment Config

Environment variables can be set directly in your shell or in a `.env` file in the project root. The application loads `.env` automatically when it starts; variables already set in the process environment take precedence.

Available environment variables:
```bash
PORT=8080                     # Server port (default: 8080)
GIN_MODE=debug               # Gin mode: debug or release
DB_DRIVER=postgres           # sqlite, mysql, mariadb, postgres, or postgresql
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=change-me
DB_NAME=hmachronos
DB_SSLMODE=disable
DB_DSN=postgres connection string (optional)
DATABASE_URL=postgres://user:password@host:5432/database (optional)
DB_PATH=./data/gorm.db       # SQLite path; only used when DB_DRIVER=sqlite
TEST_DB_PATH=./data/test.db  # SQLite database path used for tests
```

Example usage:
```bash
# PostgreSQL
export DB_DRIVER=postgres
export DB_HOST=127.0.0.1
export DB_PORT=5432
export DB_USER=hmachronos
export DB_PASSWORD='change-me'
export DB_NAME=hmachronos
export DB_SSLMODE=require
go run app.go

# SQLite for local development/tests
export PORT=3000
export DB_PATH=./data/myapp.db
go run app.go
```

For PostgreSQL, configure either `DATABASE_URL`, `DB_DSN`, or the individual
`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, and `DB_NAME` variables.
`DATABASE_URL` takes precedence, followed by `DB_DSN`.

See `.env.example` for a complete template.


## Install Dependencies
From the project root, run:
```
go build ./...
go test ./...
go mod tidy
```

## Run the Server
```bash
# Using default port 8080
go run .

# Using custom port
PORT=3000 go run .
```

## Run with Docker

Start the API against PostgreSQL using the database variables in the local
`.env` file. Docker Compose passes these variables to the container:
```bash
docker compose up --build -d
```

Check that the server is responding:
```bash
curl http://localhost:18080/api/ping/
```

The configured database must already contain the `surgical_notes` table; the
application does not create or migrate it automatically. Stop the container
with `docker compose down`.

## Surgical Notes API

The surgical notes endpoints use the existing `surgical_notes` table and do not
create or migrate it at startup. They require a valid authenticated user token
in the `Authorization: Token <token>` header.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/surgical-notes` | List notes (`limit` and `offset` supported) |
| `GET` | `/api/surgical-notes/:id` | Retrieve a note |
| `POST` | `/api/surgical-notes` | Create a note |
| `PUT` | `/api/surgical-notes/:id` | Replace a note's editable fields |
| `DELETE` | `/api/surgical-notes/:id` | Soft-delete a note |

Request and response fields use the `surgical_notes` column names in
`snake_case`, such as `patient_name` and `surgery_datetime`.

## Testing
From the project root, run:
```
go test ./...
```
or
```
go test ./... -cover
```
or
```
go test -v ./... -cover
```
depending on whether you want to see test coverage and how verbose the output you want.

## Todo
- More elegance config
- ProtoBuf support
- Code structure optimize (I think some place can use interface)
- Continuous integration (done)

## Test Coverage

Current test coverage (2026):
- **Total**: 89.2%
- **articles**: 92.1%
- **users**: 99.5%
- **common**: 85.7%

Run coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```
