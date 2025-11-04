# Excel Template Engine - Acts Service

A Go-based service for generating acts (documents) using Excel templates. The service accepts data via REST API, processes it according to business logic, and generates Excel files with populated data.

**Project completed as part of the "Legacy Projects" course in the HITs TSU Master's program.**

## 🚀 Features

- **Excel Template Processing**: Automatically replaces placeholders in Excel templates with actual data
- **MongoDB Storage**: Stores acts and positions in MongoDB for persistence
- **RESTful API**: Simple HTTP API for creating acts and generating documents
- **Number Formatting**: Automatic formatting of numbers with thousand separators (e.g., 1,234,567.89)
- **Business Logic**: Complex logic for calculating totals and selecting positions
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **Graceful Shutdown**: Proper cleanup of resources on shutdown

## 📋 Table of Contents

- [Architecture](#architecture)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Service](#running-the-service)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Development](#development)
- [Testing](#testing)

## 🏗 Architecture

```
┌─────────────┐     HTTP      ┌──────────────┐
│   Client    │ ────────────> │   Handlers   │
└─────────────┘               └──────┬───────┘
                                     │
                              ┌──────▼───────┐
                              │  ActService  │
                              └──────┬───────┘
                                     │
                    ┌────────────────┼────────────────┐
                    │                │                │
            ┌───────▼────────┐ ┌────▼─────┐ ┌───────▼────────┐
            │ ExcelService   │ │MongoDB   │ │ Utilities      │
            │ (Excelize)     │ │Repository│ │ (Formatters)   │
            └────────────────┘ └──────────┘ └────────────────┘
```

## 🛠 Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: [Gin](https://github.com/gin-gonic/gin)
- **Database**: MongoDB 7.0
- **MongoDB Driver**: [mongo-driver](https://github.com/mongodb/mongo-go-driver)
- **Excel Library**: [Excelize v2](https://github.com/xuri/excelize)
- **Configuration**: [godotenv](https://github.com/joho/godotenv)
- **Validation**: [validator/v10](https://github.com/go-playground/validator)
- **Containerization**: Docker, Docker Compose

## ✅ Prerequisites

- Go 1.21 or higher (for local development)
- Docker and Docker Compose (for containerized deployment)
- MongoDB 7.0 (if running locally without Docker)

## 📥 Installation

### Option 1: Using Docker (Recommended)

1. Clone the repository:
```bash
git clone https://github.com/stepanpotapov/Excel-Template-Engine.git
cd Excel-Template-Engine
```

2. Build and run with Docker Compose:
```bash
docker-compose up --build
```

The service will be available at `http://localhost:8080`

### Option 2: Local Development

1. Clone the repository:
```bash
git clone https://github.com/stepanpotapov/Excel-Template-Engine.git
cd Excel-Template-Engine
```

2. Install dependencies:
```bash
go mod download
```

3. Make sure MongoDB is running locally or update the `.env` file with your MongoDB URI

4. Run the service:
```bash
go run cmd/server/main.go
```

## ⚙️ Configuration

Create a `.env` file in the root directory (or copy from `.env.example`):

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
BASE_URL=http://localhost:8080

# MongoDB Configuration
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=acts_db
MONGODB_COLLECTION=acts
MONGODB_TIMEOUT=10s

# File Paths
TEMPLATE_PATH=./templates/act_template.xlsx
GENERATED_PATH=./generated

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

## 🚀 Running the Service

### With Docker Compose

```bash
# Start all services
docker-compose up

# Start in detached mode
docker-compose up -d

# View logs
docker-compose logs -f app

# Stop services
docker-compose down

# Stop and remove volumes
docker-compose down -v
```

### Local Development

```bash
# Run the service
go run cmd/server/main.go

# Build binary
go build -o server ./cmd/server

# Run binary
./server
```

### Accessing Services

- **Acts Service**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Mongo Express** (DB UI): http://localhost:8081 (admin/admin)

## 📚 API Documentation

### 1. Health Check

Check if the service is running.

```bash
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "service": "acts-service"
}
```

### 2. Create Act

Create a new act with positions and metadata.

```bash
POST /api/act/create
Content-Type: application/json
```

**Request Body:**
```json
{
  "bigAct": {
    "changed": true,
    "textFields": {
      "contractNumber": "ДГ-2025/001",
      "contractDate": "15.01.2025",
      "customer": "ООО Строительная компания",
      "contractor": "ООО Подрядчик",
      "objectName": "Строительство жилого дома"
    }
  },
  "positions": [
    {
      "currentPeriodCost": 1500000.50,
      "currentPeriodCostInspection": 75000.25,
      "currentPeriodCostConsiderations": 37500.10
    },
    {
      "currentPeriodCost": 2500000.75,
      "currentPeriodCostInspection": 125000.50,
      "currentPeriodCostConsiderations": 62500.25
    }
  ]
}
```

**Response (201 Created):**
```json
{
  "id": "673902f2a1b2c3d4e5f67890"
}
```

**cURL Example:**
```bash
curl -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{
    "bigAct": {
      "changed": true,
      "textFields": {
        "contractNumber": "ДГ-2025/001",
        "contractDate": "15.01.2025",
        "customer": "ООО Строительная компания",
        "contractor": "ООО Подрядчик",
        "objectName": "Строительство жилого дома"
      }
    },
    "positions": [
      {
        "currentPeriodCost": 1500000.50,
        "currentPeriodCostInspection": 75000.25,
        "currentPeriodCostConsiderations": 37500.10
      }
    ]
  }'
```

### 3. Generate Act

Generate an Excel file for an existing act.

```bash
GET /api/act/generate?id={actId}
```

**Parameters:**
- `id` (required): The ID of the act to generate

**Response (200 OK):**
```json
{
  "downloadLink": "/api/act/download/act_673902f2a1b2c3d4e5f67890_1730739600.xlsx"
}
```

**cURL Example:**
```bash
curl -X GET "http://localhost:8080/api/act/generate?id=673902f2a1b2c3d4e5f67890"
```

### 4. Download Act

Download a generated Excel file.

```bash
GET /api/act/download/{filename}
```

**Parameters:**
- `filename` (required): The filename returned from the generate endpoint

**Response:**
- Excel file download (application/vnd.openxmlformats-officedocument.spreadsheetml.sheet)

**cURL Example:**
```bash
curl -O "http://localhost:8080/api/act/download/act_673902f2a1b2c3d4e5f67890_1730739600.xlsx"
```

## 📁 Project Structure

```
Excel-Template-Engine/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── models/
│   │   ├── act.go                 # Act model
│   │   ├── position.go            # Position model
│   │   └── big_act.go             # BigAct model
│   ├── handlers/
│   │   └── act_handler.go         # HTTP handlers
│   ├── services/
│   │   ├── act_service.go         # Business logic
│   │   └── excel_service.go       # Excel generation
│   ├── repository/
│   │   ├── act_repository.go      # MongoDB operations
│   │   └── mongodb.go             # MongoDB connection
│   └── utils/
│       ├── number_formatter.go    # Number formatting utilities
│       └── response.go            # HTTP response utilities
├── templates/
│   └── act_template.xlsx          # Excel template
├── generated/                     # Generated Excel files
├── scripts/
│   └── simple_template.go         # Template generator script
├── docker-compose.yml             # Docker Compose config
├── Dockerfile                     # Docker image config
├── go.mod                         # Go dependencies
├── go.sum                         # Go dependencies lock
├── .env.example                   # Environment variables example
├── .gitignore                     # Git ignore rules
└── README.md                      # This file
```

## 👨‍💻 Development

### Building the Project

```bash
# Build for current platform
go build -o server ./cmd/server

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o server-linux ./cmd/server

# Build with Docker
docker build -t acts-service .
```

### Running Tests

```bash
# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -cover

# Generate coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Regenerating the Template

```bash
go run scripts/simple_template.go
```

## 🧪 Testing

### Manual Testing with cURL

1. **Create an act:**
```bash
curl -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{
    "bigAct": {
      "changed": true,
      "textFields": {
        "contractNumber": "TEST-001",
        "contractDate": "04.11.2025",
        "customer": "Test Customer",
        "contractor": "Test Contractor",
        "objectName": "Test Object"
      }
    },
    "positions": [
      {
        "currentPeriodCost": 1000000.00,
        "currentPeriodCostInspection": 50000.00,
        "currentPeriodCostConsiderations": 25000.00
      }
    ]
  }'
```

2. **Generate the act (use the ID from step 1):**
```bash
curl -X GET "http://localhost:8080/api/act/generate?id=YOUR_ACT_ID"
```

3. **Download the file:**
```bash
curl -O "http://localhost:8080/api/act/download/FILENAME_FROM_STEP_2.xlsx"
```

## 🔧 Business Logic

The service implements the following business logic when generating acts:

1. **If `bigAct.changed` is `true`:**
   - Find positions with current period costs
   - If found, calculate totals from these positions
   - If not found, fallback to positions with accumulated costs
   - Update `BigAct` with calculated totals
   - Generate Excel file
   - Save `bigActLink` to database
   - Set `changed` to `false`

2. **If `bigAct.changed` is `false`:**
   - Return existing `bigActLink`

### Number Formatting

Numbers are automatically formatted with:
- Thousand separators (commas)
- Two decimal places
- Example: `1234567.89` → `"1,234,567.89"`

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Author

**Stepan Potapov**
- GitHub: [@stepanpotapov](https://github.com/stepanpotapov)

## 🙏 Acknowledgments

- HITs TSU Master's Program - "Legacy Projects" course
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [Excelize](https://github.com/xuri/excelize) - Excel library
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)

---

**Note**: This is an educational project created as part of a university course.
