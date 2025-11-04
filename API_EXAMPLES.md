# API Examples

This file contains practical examples for testing the Acts Service API.

## Prerequisites

Make sure the service is running:
```bash
# With Docker Compose
docker-compose up

# Or locally
go run cmd/server/main.go
```

## Examples

### 1. Health Check

Check if the service is running:

```bash
curl http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "ok",
  "service": "acts-service"
}
```

---

### 2. Create Act - Simple Example

Create a basic act with one position:

```bash
curl -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{
    "bigAct": {
      "changed": true,
      "textFields": {
        "contractNumber": "TEST-001",
        "contractDate": "04.11.2025",
        "customer": "Test Customer Inc.",
        "contractor": "Test Contractor LLC",
        "objectName": "Test Construction Project"
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

**Expected Response:**
```json
{
  "id": "673902f2a1b2c3d4e5f67890"
}
```

---

### 3. Create Act - Multiple Positions

Create an act with multiple positions:

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
      },
      {
        "currentPeriodCost": 2500000.75,
        "currentPeriodCostInspection": 125000.50,
        "currentPeriodCostConsiderations": 62500.25
      },
      {
        "currentPeriodCost": 3000000.00,
        "currentPeriodCostInspection": 150000.00,
        "currentPeriodCostConsiderations": 75000.00
      }
    ]
  }'
```

---

### 4. Create Act - With Accumulated Cost

Create an act using accumulated costs (fallback):

```bash
curl -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{
    "bigAct": {
      "changed": true,
      "textFields": {
        "contractNumber": "ACC-2025/001",
        "contractDate": "20.01.2025",
        "customer": "Customer Company",
        "contractor": "Contractor Company",
        "objectName": "Project with Accumulated Costs"
      }
    },
    "positions": [
      {
        "accumulatedCost": 5000000.00
      },
      {
        "accumulatedCost": 3000000.50
      }
    ]
  }'
```

---

### 5. Generate Act

Generate Excel file for an act (replace `{ACT_ID}` with actual ID from create response):

```bash
# Replace {ACT_ID} with the actual ID from the create response
ACT_ID="673902f2a1b2c3d4e5f67890"

curl -X GET "http://localhost:8080/api/act/generate?id=${ACT_ID}"
```

**Expected Response:**
```json
{
  "downloadLink": "/api/act/download/act_673902f2a1b2c3d4e5f67890_1730739600.xlsx"
}
```

---

### 6. Download Generated File

Download the generated Excel file:

```bash
# Replace {FILENAME} with the filename from the generate response
FILENAME="act_673902f2a1b2c3d4e5f67890_1730739600.xlsx"

curl -O "http://localhost:8080/api/act/download/${FILENAME}"
```

This will download the file to your current directory.

---

## Complete Workflow Example

Here's a complete workflow from creating an act to downloading the file:

```bash
#!/bin/bash

# Step 1: Create an act
echo "Step 1: Creating act..."
RESPONSE=$(curl -s -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{
    "bigAct": {
      "changed": true,
      "textFields": {
        "contractNumber": "DEMO-2025/001",
        "contractDate": "04.11.2025",
        "customer": "Demo Customer Corporation",
        "contractor": "Demo Contractor LLC",
        "objectName": "Demo Construction Project"
      }
    },
    "positions": [
      {
        "currentPeriodCost": 1234567.89,
        "currentPeriodCostInspection": 123456.78,
        "currentPeriodCostConsiderations": 12345.67
      }
    ]
  }')

# Extract ID
ACT_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Act created with ID: $ACT_ID"

# Step 2: Generate the act
echo "Step 2: Generating Excel file..."
GEN_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/act/generate?id=${ACT_ID}")
echo "Generate response: $GEN_RESPONSE"

# Extract download link
DOWNLOAD_LINK=$(echo $GEN_RESPONSE | grep -o '"downloadLink":"[^"]*' | cut -d'"' -f4)
echo "Download link: $DOWNLOAD_LINK"

# Step 3: Download the file
FILENAME=$(basename "$DOWNLOAD_LINK")
echo "Step 3: Downloading file: $FILENAME"
curl -O "http://localhost:8080${DOWNLOAD_LINK}"

echo "Workflow complete! File downloaded: $FILENAME"
```

Save this as `test_workflow.sh`, make it executable (`chmod +x test_workflow.sh`), and run it:

```bash
./test_workflow.sh
```

---

## Error Cases

### Missing ID Parameter

```bash
curl -X GET "http://localhost:8080/api/act/generate"
```

**Expected Response:**
```json
{
  "error": "Bad Request",
  "message": "ID parameter is required",
  "code": 400
}
```

### Invalid ID

```bash
curl -X GET "http://localhost:8080/api/act/generate?id=invalid_id"
```

**Expected Response:**
```json
{
  "error": "Internal Server Error",
  "message": "Failed to generate act",
  "code": 500
}
```

### File Not Found

```bash
curl -X GET "http://localhost:8080/api/act/download/nonexistent_file.xlsx"
```

**Expected Response:**
```json
{
  "error": "Not Found",
  "message": "File not found",
  "code": 404
}
```

---

## Testing with HTTPie

If you have [HTTPie](https://httpie.io/) installed, you can use these cleaner commands:

```bash
# Create act
http POST localhost:8080/api/act/create \
  bigAct:='{"changed":true,"textFields":{"contractNumber":"TEST-001"}}' \
  positions:='[{"currentPeriodCost":1000000}]'

# Generate act
http GET localhost:8080/api/act/generate id==673902f2a1b2c3d4e5f67890

# Download file
http --download localhost:8080/api/act/download/act_673902f2a1b2c3d4e5f67890_1730739600.xlsx
```

---

## Testing with Postman

Import this collection to Postman:

1. Create Act:
   - Method: POST
   - URL: `http://localhost:8080/api/act/create`
   - Headers: `Content-Type: application/json`
   - Body: Raw JSON (see examples above)

2. Generate Act:
   - Method: GET
   - URL: `http://localhost:8080/api/act/generate`
   - Params: `id` = `{ACT_ID}`

3. Download Act:
   - Method: GET
   - URL: `http://localhost:8080/api/act/download/{filename}`

---

## Notes

- All monetary values in the response will be formatted with thousand separators
  - Example: 1234567.89 → "1,234,567.89"
- The Excel template placeholders ({{key}}) will be replaced with actual values
- Generated files are stored in the `./generated` directory
- MongoDB data is persisted in Docker volumes

