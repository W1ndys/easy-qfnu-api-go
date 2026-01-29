# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run

- **Install Dependencies**: `go mod tidy`
- **Run Locally**: `go run main.go`
- **Build**: `go build -o easy-qfnu-api-go`
- **Test**: `go test ./...`
- **Lint**: `go vet ./...`

## Architecture & Code Structure

This is a Go-based API gateway and scraping service for QFNU campus services, built with the Gin framework. It aggregates data from educational systems (grades, schedules) and other sources.

### Core Structure
- **Entry Point**: `main.go` initializes the logger, router, and embeds static assets.
- **Routing**: `router/router.go` defines API groups (`/api/v1`) and HTML rendering.
- **API Handlers**: Located in `api/` (e.g., `api/v1/zhjw` for educational system, `api/v1/questions` for question bank).
- **Services**: `services/` contains business logic, particularly the scraping and HTML parsing logic using `go-resty` and `goquery`.
- **Middleware**: `middleware/` handles logging, CORS, and authentication (`AuthRequired`).
- **Frontend**: `web/` contains static assets and HTML templates which are embedded into the Go binary using `go:embed`.
- **Common**: `common/` holds shared utilities like logging (`logger`) and standardized API responses (`response`).

### Key Concepts
- **Single Binary**: The project uses `embed` to package HTML/CSS/JS with the binary for easy deployment.
- **Scraping**: Data is primarily fetched by simulating HTTP requests to school servers and parsing HTML responses.
- **Configuration**: Environment variables are loaded via `.env` (using `godotenv`). Key variables: `PORT`, `GIN_MODE`.
- **Database**: Uses SQLite (`modernc.org/sqlite`) for local data storage where necessary.

### Development Conventions
- **Routing**: API routes are versioned (`/api/v1`).
- **Logging**: Uses a custom logger setup in `common/logger`.
- **Response Format**: Standardized JSON responses via `common/response` package.
