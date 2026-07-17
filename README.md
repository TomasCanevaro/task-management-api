# Task Management API

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
go test ./internal/services
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

These are some features that could be added in the future:

- Database persistence
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