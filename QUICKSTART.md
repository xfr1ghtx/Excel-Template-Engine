# Quick Start Guide

Get the Acts Service up and running in 5 minutes!

## 🚀 Quick Start with Docker (Recommended)

### 1. Start the Services

```bash
docker-compose up -d
```

Wait for the services to start (about 30 seconds).

### 2. Verify Services are Running

```bash
# Check health
curl http://localhost:8080/health

# Expected output:
# {"status":"ok","service":"acts-service"}
```

### 3. Create Your First Act

```bash
curl -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{
    "bigAct": {
      "changed": true,
      "textFields": {
        "contractNumber": "DEMO-001",
        "contractDate": "04.11.2025",
        "customer": "My Company",
        "contractor": "Contractor LLC",
        "objectName": "My First Project"
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

**Save the returned ID!** You'll need it for the next step.

Example response:
```json
{"id":"673902f2a1b2c3d4e5f67890"}
```

### 4. Generate Excel File

Replace `YOUR_ACT_ID` with the ID from step 3:

```bash
curl -X GET "http://localhost:8080/api/act/generate?id=YOUR_ACT_ID"
```

Example:
```bash
curl -X GET "http://localhost:8080/api/act/generate?id=673902f2a1b2c3d4e5f67890"
```

Example response:
```json
{"downloadLink":"/api/act/download/act_673902f2a1b2c3d4e5f67890_1730739600.xlsx"}
```

### 5. Download the File

Replace `FILENAME` with the filename from step 4:

```bash
curl -O "http://localhost:8080/api/act/download/FILENAME.xlsx"
```

Example:
```bash
curl -O "http://localhost:8080/api/act/download/act_673902f2a1b2c3d4e5f67890_1730739600.xlsx"
```

### 6. Open the Excel File

Open the downloaded `.xlsx` file in Excel, LibreOffice, or Google Sheets to see your generated act!

---

## 🎯 All-in-One Script

Save this as `test.sh`:

```bash
#!/bin/bash
set -e

echo "🚀 Testing Acts Service"
echo "======================="
echo ""

# Create act
echo "📝 Step 1: Creating act..."
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{
    "bigAct": {
      "changed": true,
      "textFields": {
        "contractNumber": "DEMO-001",
        "contractDate": "04.11.2025",
        "customer": "My Company",
        "contractor": "Contractor LLC",
        "objectName": "My First Project"
      }
    },
    "positions": [
      {
        "currentPeriodCost": 1000000.00,
        "currentPeriodCostInspection": 50000.00,
        "currentPeriodCostConsiderations": 25000.00
      }
    ]
  }')

ACT_ID=$(echo $CREATE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "✅ Act created with ID: $ACT_ID"
echo ""

# Generate Excel
echo "📊 Step 2: Generating Excel file..."
GEN_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/act/generate?id=${ACT_ID}")
DOWNLOAD_LINK=$(echo $GEN_RESPONSE | grep -o '"downloadLink":"[^"]*' | cut -d'"' -f4)
echo "✅ File generated: $DOWNLOAD_LINK"
echo ""

# Download file
FILENAME=$(basename "$DOWNLOAD_LINK")
echo "⬇️  Step 3: Downloading file..."
curl -s -O "http://localhost:8080${DOWNLOAD_LINK}"
echo "✅ File downloaded: $FILENAME"
echo ""

echo "🎉 Success! Open $FILENAME to view your act!"
```

Make it executable and run:

```bash
chmod +x test.sh
./test.sh
```

---

## 🔍 View Database (Optional)

Access Mongo Express to view your data:

1. Open browser: http://localhost:8081
2. Login: `admin` / `admin`
3. Navigate to: `acts_db` → `acts`

---

## 🛑 Stop Services

```bash
docker-compose down
```

To also remove the database:

```bash
docker-compose down -v
```

---

## 💻 Local Development (Without Docker)

### Prerequisites

- Go 1.21+
- MongoDB running on localhost:27017

### Steps

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Start MongoDB locally** (if not already running)

3. **Run the service:**
   ```bash
   go run cmd/server/main.go
   ```

4. **Follow steps 2-6 from the Docker guide above**

---

## 📚 Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check [API_EXAMPLES.md](API_EXAMPLES.md) for more examples
- Explore the [PLAN.md](PLAN.md) to understand the architecture

---

## 🐛 Troubleshooting

### Service won't start

```bash
# Check if ports are already in use
lsof -i :8080
lsof -i :27017

# View logs
docker-compose logs app
```

### Can't connect to MongoDB

```bash
# Restart MongoDB
docker-compose restart mongodb

# Check MongoDB logs
docker-compose logs mongodb
```

### Generated files not found

Check the `generated/` directory in the project root.

---

## 🎓 Understanding the Output

When you open the generated Excel file, you'll see:

- **Contract Number**: From your request (`contractNumber`)
- **Contract Date**: From your request (`contractDate`)
- **Customer**: From your request (`customer`)
- **Contractor**: From your request (`contractor`)
- **Object Name**: From your request (`objectName`)
- **Total Cost**: Sum of all position costs (formatted: `1,234,567.89`)
- **Total Inspection**: Sum of inspection costs (formatted)
- **Total Considerations**: Sum of consideration costs (formatted)
- **Position IDs**: Comma-separated IDs of all positions
- **Created At**: When the act was created

All numbers are automatically formatted with thousand separators!

---

**Happy generating! 🎉**

