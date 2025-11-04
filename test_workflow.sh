#!/bin/bash

set -e

echo "🎯 Excel Template Engine - Test Workflow"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if service is running
echo -e "${BLUE}Checking service health...${NC}"
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ Service is not running!"
    echo "Please start the service first:"
    echo "  docker-compose up -d"
    echo "  OR"
    echo "  go run cmd/server/main.go"
    exit 1
fi
echo -e "${GREEN}✅ Service is healthy${NC}"
echo ""

# Step 1: Create act
echo -e "${BLUE}📝 Step 1: Creating act...${NC}"
CREATE_RESPONSE=$(curl -s -X POST http://localhost:8080/api/act/create \
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
      },
      {
        "currentPeriodCost": 2500000.00,
        "currentPeriodCostInspection": 250000.00,
        "currentPeriodCostConsiderations": 25000.00
      }
    ]
  }')

# Extract ID
ACT_ID=$(echo $CREATE_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)

if [ -z "$ACT_ID" ]; then
    echo "❌ Failed to create act"
    echo "Response: $CREATE_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Act created with ID: $ACT_ID${NC}"
echo ""

# Step 2: Generate the act
echo -e "${BLUE}📊 Step 2: Generating Excel file...${NC}"
GEN_RESPONSE=$(curl -s -X GET "http://localhost:8080/api/act/generate?id=${ACT_ID}")

# Extract download link
DOWNLOAD_LINK=$(echo $GEN_RESPONSE | grep -o '"downloadLink":"[^"]*' | cut -d'"' -f4)

if [ -z "$DOWNLOAD_LINK" ]; then
    echo "❌ Failed to generate act"
    echo "Response: $GEN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ File generated: $DOWNLOAD_LINK${NC}"
echo ""

# Step 3: Download the file
FILENAME=$(basename "$DOWNLOAD_LINK")
echo -e "${BLUE}⬇️  Step 3: Downloading file...${NC}"

if curl -s -f -O "http://localhost:8080${DOWNLOAD_LINK}"; then
    echo -e "${GREEN}✅ File downloaded successfully: $FILENAME${NC}"
else
    echo "❌ Failed to download file"
    exit 1
fi
echo ""

# Summary
echo "=========================================="
echo -e "${GREEN}🎉 Test completed successfully!${NC}"
echo ""
echo "Summary:"
echo "  - Act ID: $ACT_ID"
echo "  - File: $FILENAME"
echo "  - Total Cost: 1,234,567.89 + 2,500,000.00 = 3,734,567.89"
echo "  - Total Inspection: 123,456.78 + 250,000.00 = 373,456.78"
echo "  - Total Considerations: 12,345.67 + 25,000.00 = 37,345.67"
echo ""
echo "📂 Open $FILENAME in Excel to view the generated act!"
echo ""
