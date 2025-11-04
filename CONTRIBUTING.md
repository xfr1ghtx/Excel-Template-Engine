# Contributing to Excel Template Engine

Thank you for your interest in contributing to the Excel Template Engine project!

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Git
- A code editor (VS Code, GoLand, etc.)

### Getting Started

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/stepanpotapov/Excel-Template-Engine.git
   cd Excel-Template-Engine
   ```

2. **Install dependencies:**
   ```bash
   make install-deps
   # or
   go mod download
   ```

3. **Generate the Excel template:**
   ```bash
   make template
   # or
   go run scripts/simple_template.go
   ```

4. **Run the service locally:**
   ```bash
   make run
   # or
   go run cmd/server/main.go
   ```

## Project Structure

```
Excel-Template-Engine/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/         # Configuration management
│   ├── handlers/       # HTTP request handlers
│   ├── models/         # Data models
│   ├── repository/     # Database operations
│   ├── services/       # Business logic
│   └── utils/          # Utility functions
├── templates/          # Excel templates
├── generated/          # Generated Excel files
└── scripts/            # Helper scripts
```

## Code Style

### Go Code Guidelines

1. **Follow standard Go conventions:**
   - Use `gofmt` for formatting
   - Follow [Effective Go](https://golang.org/doc/effective_go)
   - Use meaningful variable names

2. **Format your code:**
   ```bash
   make fmt
   # or
   go fmt ./...
   ```

3. **Run the linter:**
   ```bash
   make lint
   # or
   golangci-lint run
   ```

### Comments

- Add comments for exported functions and types
- Use godoc-style comments
- Explain complex logic

Example:
```go
// FormatNumber formats a number with thousand separators and 2 decimal places.
// Example: 1234567.89 -> "1,234,567.89"
func FormatNumber(value float64) string {
    // implementation
}
```

## Testing

### Writing Tests

1. **Create test files:**
   - Name: `*_test.go`
   - Location: Same directory as the code being tested

2. **Test structure:**
   ```go
   func TestFunctionName(t *testing.T) {
       tests := []struct {
           name     string
           input    interface{}
           expected interface{}
       }{
           {
               name:     "test case 1",
               input:    value1,
               expected: expected1,
           },
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // test implementation
           })
       }
   }
   ```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
go test ./internal/utils -v
```

## Making Changes

### Workflow

1. **Create a new branch:**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes**

3. **Write tests for your changes**

4. **Run tests:**
   ```bash
   make test
   ```

5. **Format code:**
   ```bash
   make fmt
   ```

6. **Commit your changes:**
   ```bash
   git add .
   git commit -m "Add feature: description"
   ```

7. **Push to your fork:**
   ```bash
   git push origin feature/my-feature
   ```

8. **Create a Pull Request**

### Commit Message Guidelines

Use clear and descriptive commit messages:

- `Add feature: user authentication`
- `Fix bug: null pointer in act generation`
- `Update docs: API examples`
- `Refactor: improve number formatting`
- `Test: add tests for Excel service`

## Adding New Features

### Example: Adding a New API Endpoint

1. **Define the handler in `internal/handlers/`:**
   ```go
   func (h *ActHandler) NewEndpoint(c *gin.Context) {
       // implementation
   }
   ```

2. **Add business logic in `internal/services/`:**
   ```go
   func (s *actService) NewMethod() error {
       // implementation
   }
   ```

3. **Register the route in `cmd/server/main.go`:**
   ```go
   act.GET("/new-endpoint", actHandler.NewEndpoint)
   ```

4. **Add tests**

5. **Update documentation**

### Example: Adding a New Model Field

1. **Update the model in `internal/models/`:**
   ```go
   type Act struct {
       // existing fields
       NewField string `json:"newField" bson:"newField"`
   }
   ```

2. **Update database operations if needed**

3. **Update Excel template processing if needed**

4. **Add tests**

5. **Update API documentation**

## Documentation

### Update Documentation When:

- Adding new features
- Changing API behavior
- Modifying configuration options
- Adding new dependencies

### Documentation Files:

- `README.md` - Main documentation
- `API_EXAMPLES.md` - API usage examples
- `QUICKSTART.md` - Quick start guide
- `PLAN.md` - Architecture and design

## Common Tasks

### Adding a New Dependency

```bash
# Add the dependency
go get github.com/some/package

# Tidy up
go mod tidy

# Commit changes
git add go.mod go.sum
git commit -m "Add dependency: package-name"
```

### Updating the Excel Template

1. **Modify `scripts/simple_template.go`**

2. **Regenerate the template:**
   ```bash
   make template
   ```

3. **Test with the new template**

4. **Commit both the script and template**

### Debugging

#### Local Debugging

```bash
# Run with verbose logging
LOG_LEVEL=debug go run cmd/server/main.go

# Check MongoDB connection
docker-compose logs mongodb

# View application logs
docker-compose logs app
```

#### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Run with debugger
dlv debug cmd/server/main.go
```

## Pull Request Process

1. **Ensure all tests pass**
2. **Update documentation**
3. **Add a clear description of changes**
4. **Reference any related issues**
5. **Wait for review**
6. **Address review comments**

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] All tests pass
- [ ] No linter errors
- [ ] Commit messages are clear

## Need Help?

- Check existing issues and documentation
- Ask questions in pull request comments
- Review the [PLAN.md](PLAN.md) for architecture details

## Code of Conduct

- Be respectful and constructive
- Welcome newcomers
- Focus on what's best for the project
- Show empathy towards other contributors

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing!** 🎉

