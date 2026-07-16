# Task Management API

A backend REST API for managing tasks with role-based authorization, built in Go.

This project was developed as a take-home assessment for a Software Engineer II position.

---

# Features

- REST API
- Role-based authorization
- Forward-only task lifecycle
- In-memory storage
- Notification logging and persistence
- Unit tests
- Clean project structure

---

# Tech Stack

- Go
- Gorilla Mux
- Standard Library
- In-memory repository

---

# Project Structure

```
task-management-api/

cmd/
    api/
        main.go

internal/
    handlers/
    middleware/
    models/
    repository/
    services/
    storage/

tests/

go.mod
README.md
```

---

# Setup

## Clone

```bash
git clone <repository-url>

cd task-management-api
```

## Install dependencies

```bash
go mod tidy
```

## Run

```bash
go run ./cmd/api
```

The API will start on

```
http://localhost:8080
```

---

# Seed Users

| ID | Name | Role |
|----|------|------|
|1|Supervisor|SUPERVISOR|
|2|Worker One|WORKER|
|3|Worker Two|WORKER|

Use the following header to identify the caller:

```
X-User-Id: 1
```

Example:

```
X-User-Id: 2
```

---

# Task Lifecycle

```
CREATED
    ↓
ASSIGNED
    ↓
IN_PROGRESS
    ↓
COMPLETED
```

Status transitions are forward-only.

---

# API Endpoints

## Health

```
GET /health
```

---

## Create Task

```
POST /api/tasks
```

Supervisor only.

Example body:

```json
{
    "title": "Write documentation",
    "description": "Finish README"
}
```

---

## Assign Task

```
POST /api/tasks/{id}/assign
```

Supervisor only.

Example:

```json
{
    "worker_id": 2
}
```

---

## Update Task Status

```
PATCH /api/tasks/{id}/status
```

Worker only.

Example:

```json
{
    "status": "IN_PROGRESS"
}
```

---

## List Tasks

```
GET /api/tasks
```

Supervisor:
- returns every task

Worker:
- returns assigned tasks only

---

## Get Single Task

```
GET /api/tasks/{id}
```

---

## Get Notifications

```
GET /api/notifications
```

Worker only.

---

# Running Tests

Run all tests:

```bash
go test ./...
```

---

# Design Decisions

The project is organized into several layers.

### Handlers

Responsible for:

- parsing HTTP requests
- validating input
- returning HTTP responses

### Services

Responsible for:

- business rules
- authorization
- task lifecycle validation
- notification creation

### Repository

The service depends on a repository interface instead of a concrete implementation.

This allows the storage layer to be replaced without changing the business logic.

### Storage

The assessment uses an in-memory implementation.

This can later be replaced with PostgreSQL or another database by implementing the repository interface.

---

# Future Improvements

Given more time, the following improvements could be added:

- PostgreSQL persistence
- Database migrations
- JWT authentication
- Docker support
- OpenAPI / Swagger documentation
- Structured logging
- Configuration using environment variables
- Pagination
- Request logging middleware
- Integration tests
- CI/CD pipeline

---

# Author

<Your Name>