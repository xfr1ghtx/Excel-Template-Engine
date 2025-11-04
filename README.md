# Excel Template Engine — Acts Service

A small Go service that generates Excel documents (acts) from a template via a simple REST API.

## Quick Start (Docker)

```bash
# Start services
docker-compose up -d --build

# View logs
docker-compose logs -f app

# Stop
docker-compose down
```

Service: http://localhost:8080  |  Health: http://localhost:8080/health  |  DB UI: http://localhost:8081 (admin/admin)

## Usage (API)

- Health
```bash
curl http://localhost:8080/health
```

- Create Act (save the returned id)
```bash
curl -s -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{
    "bigAct": {
      "changed": true,
      "textFields": {
        "contractNumber": "DEMO-001",
        "contractDate": "04.11.2025",
        "customer": "Customer",
        "contractor": "Contractor",
        "objectName": "Project"
      }
    },
    "positions": [ { "currentPeriodCost": 1000000.00 } ]
  }'
```

- Generate Act (replace YOUR_ACT_ID)
```bash
curl -s "http://localhost:8080/api/act/generate?id=YOUR_ACT_ID"
```

- Download File (replace FILENAME)
```bash
curl -O "http://localhost:8080/api/act/download/FILENAME.xlsx"
```

## Local Run (optional)

Requirements: Go 1.24+, MongoDB.

```bash
# Install deps
go mod download

# Run
go run cmd/server/main.go
```

Configuration via environment variables (see `docker-compose.yml`), key vars:
- SERVER_PORT (default 8080)
- MONGODB_URI (e.g., mongodb://mongodb:27017)
- TEMPLATE_PATH (default ./templates/act_template.xlsx)
- GENERATED_PATH (default ./generated)
