# Hello Gin Project

## Quick Start

```bash
go run .
```

## Configuration

- `HOST` (default: empty, binds on all interfaces)
- `PORT` (default: `8080`)

Example:

```bash
HOST=127.0.0.1 PORT=8080 go run .
```

## Endpoints

- `GET /hello`
- `GET /healthz`

Example:

```bash
curl -sS http://localhost:8080/hello
curl -sS http://localhost:8080/healthz
```