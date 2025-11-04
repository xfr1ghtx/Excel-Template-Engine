# Quick Reference Card

Essential commands and information for Excel Template Engine.

## 🚀 Getting Started (30 seconds)

```bash
# Start everything
docker-compose up -d

# Test it works
curl http://localhost:8080/health

# Run complete test
./test_workflow.sh
```

## 📋 Common Commands

### Docker
```bash
make docker-up        # Start services
make docker-down      # Stop services
make docker-logs      # View logs
make docker-restart   # Restart services
```

### Development
```bash
make run             # Run locally
make build           # Build binary
make test            # Run tests
make test-coverage   # Run tests with coverage
make fmt             # Format code
```

### Templates
```bash
make template        # Generate Excel template
```

## 🔗 API Quick Reference

### Health Check
```bash
curl http://localhost:8080/health
```

### Create Act
```bash
curl -X POST http://localhost:8080/api/act/create \
  -H "Content-Type: application/json" \
  -d '{"bigAct":{"changed":true,"textFields":{"contractNumber":"TEST-001"}},"positions":[{"currentPeriodCost":1000000}]}'
```

### Generate Act
```bash
curl "http://localhost:8080/api/act/generate?id=YOUR_ACT_ID"
```

### Download File
```bash
curl -O "http://localhost:8080/api/act/download/FILENAME.xlsx"
```

## 📂 Project Structure

```
Excel-Template-Engine/
├── cmd/server/          → Main application
├── internal/            → Internal packages
│   ├── models/         → Data models
│   ├── services/       → Business logic
│   ├── repository/     → Database
│   ├── handlers/       → HTTP handlers
│   ├── config/         → Configuration
│   └── utils/          → Utilities
├── templates/          → Excel templates
└── generated/          → Generated files
```

## 🌐 Service URLs

- **API**: http://localhost:8080
- **Health**: http://localhost:8080/health
- **Mongo Express**: http://localhost:8081 (admin/admin)

## 🔧 Environment Variables

```env
SERVER_PORT=8080
MONGODB_URI=mongodb://mongodb:27017
MONGODB_DATABASE=acts_db
TEMPLATE_PATH=./templates/act_template.xlsx
GENERATED_PATH=./generated
```

## 📖 Documentation Files

- `README.md` - Full documentation
- `QUICKSTART.md` - 5-minute guide
- `API_EXAMPLES.md` - API examples
- `CONTRIBUTING.md` - Development guide
- `DEPLOYMENT.md` - Deployment guide

## 🐛 Troubleshooting

### Service won't start
```bash
# Check what's using port 8080
lsof -i :8080

# View logs
docker-compose logs app
```

### MongoDB issues
```bash
# Restart MongoDB
docker-compose restart mongodb

# Check MongoDB logs
docker-compose logs mongodb
```

### Build errors
```bash
# Clean and rebuild
make clean
go mod tidy
make build
```

## 🧪 Testing

### Quick test
```bash
./test_workflow.sh
```

### Unit tests
```bash
make test
```

### Manual test
```bash
# 1. Create act (save the ID)
curl -X POST http://localhost:8080/api/act/create -H "Content-Type: application/json" -d '...'

# 2. Generate (use ID from step 1)
curl "http://localhost:8080/api/act/generate?id=YOUR_ID"

# 3. Download (use filename from step 2)
curl -O "http://localhost:8080/api/act/download/FILENAME"
```

## 📊 Key Concepts

### Act Structure
```json
{
  "bigAct": {
    "changed": true,
    "textFields": {
      "contractNumber": "...",
      "contractDate": "...",
      "customer": "...",
      "contractor": "...",
      "objectName": "..."
    }
  },
  "positions": [
    {
      "currentPeriodCost": 1000000.00,
      "currentPeriodCostInspection": 50000.00,
      "currentPeriodCostConsiderations": 25000.00
    }
  ]
}
```

### Excel Template
Placeholders: `{{key}}` gets replaced with values
- `{{contractNumber}}`
- `{{totalCost}}` (auto-formatted: 1,234,567.89)
- `{{createdAt}}`

## 💡 Tips

1. **First time?** → Read `QUICKSTART.md`
2. **API details?** → Check `API_EXAMPLES.md`
3. **Contributing?** → See `CONTRIBUTING.md`
4. **Deploying?** → Read `DEPLOYMENT.md`
5. **Problems?** → Check logs with `docker-compose logs`

## 🎯 One-Liner Examples

```bash
# Complete workflow in one command
ACT_ID=$(curl -s -X POST localhost:8080/api/act/create -H "Content-Type: application/json" -d '{"bigAct":{"changed":true},"positions":[{"currentPeriodCost":1000000}]}' | grep -o '"id":"[^"]*' | cut -d'"' -f4) && curl -s "localhost:8080/api/act/generate?id=$ACT_ID"

# Health check with formatted output
curl -s localhost:8080/health | python -m json.tool

# Count acts in database
docker-compose exec mongodb mongo acts_db --eval "db.acts.count()"
```

---

**Need help?** Check the documentation or open an issue on GitHub!

