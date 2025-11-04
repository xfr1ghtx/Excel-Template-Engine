# План реализации сервиса генерации актов на GO

## Обзор проекта

Сервис для создания и генерации актов с использованием Excel-шаблонов. Сервис будет работать с MongoDB для хранения данных и использовать библиотеку Excelize для работы с Excel файлами.

## Архитектура проекта

```
YBNFs/
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа приложения
├── internal/
│   ├── config/
│   │   └── config.go              # Конфигурация приложения
│   ├── models/
│   │   ├── act.go                 # Модель акта
│   │   ├── position.go            # Модель позиции
│   │   └── big_act.go             # Модель большого акта
│   ├── handlers/
│   │   └── act_handler.go         # HTTP handlers для эндпоинтов
│   ├── services/
│   │   ├── act_service.go         # Бизнес-логика работы с актами
│   │   └── excel_service.go       # Работа с Excel (Excelize)
│   ├── repository/
│   │   ├── act_repository.go      # Работа с MongoDB для актов
│   │   └── mongodb.go             # Подключение к MongoDB
│   └── utils/
│       ├── number_formatter.go    # Форматирование чисел
│       └── response.go            # Утилиты для HTTP ответов
├── templates/
│   └── act_template.xlsx          # Шаблон Excel файла
├── generated/
│   └── .gitkeep                   # Директория для сгенерированных файлов
├── docker-compose.yml             # Docker Compose конфигурация
├── Dockerfile                     # Dockerfile для GO сервиса
├── go.mod                         # Go модули
├── go.sum                         # Go зависимости
├── .env.example                   # Пример переменных окружения
├── .gitignore                     # Git ignore файл
└── README.md                      # Документация проекта
```

## Технологический стек

- **Язык**: Go 1.21+
- **Web Framework**: Gin (легковесный и производительный)
- **База данных**: MongoDB
- **MongoDB Driver**: go.mongodb.org/mongo-driver
- **Excel библиотека**: github.com/xuri/excelize/v2
- **Конфигурация**: github.com/joho/godotenv
- **Валидация**: github.com/go-playground/validator/v10
- **Контейнеризация**: Docker + Docker Compose

## Подробный план реализации

### Этап 1: Инициализация проекта и базовая структура

**Задачи:**
1. Инициализировать Go модуль
2. Создать базовую структуру директорий
3. Настроить .gitignore
4. Создать Dockerfile для GO приложения
5. Создать docker-compose.yml с сервисами:
   - GO приложение
   - MongoDB
   - Mongo Express (опционально, для удобства разработки)

**Файлы:**
- `go.mod`, `go.sum`
- `.gitignore`
- `Dockerfile`
- `docker-compose.yml`
- `.env.example`

### Этап 2: Модели данных

**Модель Act (internal/models/act.go):**
```go
type Act struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    BigAct    *BigAct           `json:"bigAct,omitempty" bson:"bigAct,omitempty"`
    Positions []Position        `json:"positions,omitempty" bson:"positions,omitempty"`
    // Общие поля акта
    CreatedAt time.Time         `json:"createdAt" bson:"createdAt"`
    UpdatedAt time.Time         `json:"updatedAt" bson:"updatedAt"`
}
```

**Модель BigAct (internal/models/big_act.go):**
```go
type BigAct struct {
    Changed                bool       `json:"changed" bson:"changed"`
    TotalCost             float64    `json:"totalCost,omitempty" bson:"totalCost,omitempty"`
    TotalCostInspection   float64    `json:"totalCostInspection,omitempty" bson:"totalCostInspection,omitempty"`
    TotalCostConsiderations float64  `json:"totalCostConsiderations,omitempty" bson:"totalCostConsiderations,omitempty"`
    PositionIDs           string     `json:"positionIds,omitempty" bson:"positionIds,omitempty"`
    BigActLink            string     `json:"bigActLink,omitempty" bson:"bigActLink,omitempty"`
    // Текстовые значения для подстановки в шаблон
    TextFields            map[string]interface{} `json:"textFields,omitempty" bson:"textFields,omitempty"`
}
```

**Модель Position (internal/models/position.go):**
```go
type Position struct {
    ID                              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    CurrentPeriodCost               *float64          `json:"currentPeriodCost,omitempty" bson:"currentPeriodCost,omitempty"`
    CurrentPeriodCostInspection     *float64          `json:"currentPeriodCostInspection,omitempty" bson:"currentPeriodCostInspection,omitempty"`
    CurrentPeriodCostConsiderations *float64          `json:"currentPeriodCostConsiderations,omitempty" bson:"currentPeriodCostConsiderations,omitempty"`
    AccumulatedCost                 *float64          `json:"accumulatedCost,omitempty" bson:"accumulatedCost,omitempty"`
    // Дополнительные поля позиции
}
```

### Этап 3: Конфигурация приложения

**Файл internal/config/config.go:**
- Чтение переменных окружения
- Параметры подключения к MongoDB
- Порт сервера
- Пути к директориям (templates, generated)
- URL для скачивания файлов

**Переменные окружения (.env):**
```
SERVER_PORT=8080
MONGODB_URI=mongodb://mongodb:27017
MONGODB_DATABASE=acts_db
MONGODB_COLLECTION=acts
TEMPLATE_PATH=./templates/act_template.xlsx
GENERATED_PATH=./generated
BASE_URL=http://localhost:8080
```

### Этап 4: Работа с MongoDB

**Файл internal/repository/mongodb.go:**
- Функция подключения к MongoDB
- Управление соединением
- Ping для проверки соединения

**Файл internal/repository/act_repository.go:**
```go
type ActRepository interface {
    Create(ctx context.Context, act *models.Act) (string, error)
    FindByID(ctx context.Context, id string) (*models.Act, error)
    Update(ctx context.Context, id string, act *models.Act) error
}
```

**Реализация:**
- `Create`: Вставка нового акта в MongoDB
- `FindByID`: Поиск акта по ID
- `Update`: Обновление акта (для bigActLink)

### Этап 5: Форматирование чисел

**Файл internal/utils/number_formatter.go:**
```go
func FormatNumber(value float64) string {
    // Форматирование в формат: 123,456,789.00
    // Логика:
    // 1. Разделить на целую и дробную части
    // 2. Добавить запятые каждые 3 цифры в целой части
    // 3. Оставить 2 знака после точки
}
```

**Примеры:**
- `1234567.89` → `1,234,567.89`
- `1000` → `1,000.00`
- `500.5` → `500.50`

### Этап 6: Сервис работы с Excel

**Файл internal/services/excel_service.go:**

**Основные функции:**
```go
type ExcelService interface {
    GenerateAct(act *models.Act, outputPath string) error
}
```

**Логика генерации:**
1. Открыть шаблон Excel файла с помощью Excelize
2. Найти все ячейки с паттерном `{{key}}`
3. Для каждой найденной ячейки:
   - Извлечь ключ
   - Найти соответствующее значение в модели Act
   - Если значение - число, отформатировать его
   - Заменить `{{key}}` на значение
4. Сохранить файл в директорию `generated/`

**Поиск и замена ключей:**
- Итерация по всем листам (sheets)
- Итерация по всем ячейкам
- Regex поиск паттерна `\{\{([^}]+)\}\}`
- Замена значений

**Обработка чисел:**
- Определение типа данных (число/строка)
- Применение форматирования для чисел

### Этап 7: Бизнес-логика (Act Service)

**Файл internal/services/act_service.go:**

**Интерфейс:**
```go
type ActService interface {
    CreateAct(ctx context.Context, act *models.Act) (string, error)
    GenerateAct(ctx context.Context, actID string) (string, error)
}
```

**Метод CreateAct:**
1. Валидация входных данных
2. Установка timestamp (createdAt, updatedAt)
3. Сохранение в MongoDB через repository
4. Возврат ID созданного акта

**Метод GenerateAct (основной алгоритм):**

```
1. Получить акт из БД по ID
2. Проверить bigAct.changed:
   
   ЕСЛИ bigActChanged == true:
   
   a) Найти позиции с currentPeriodCost*, != null:
      - Фильтровать positions где хотя бы одно из полей не null:
        * currentPeriodCost
        * currentPeriodCostInspection
        * currentPeriodCostConsiderations
   
   b) ЕСЛИ нашли позиции:
      - Конкатенировать их ID (например, через запятую)
      - Перейти к шагу d)
   
   c) ИНАЧЕ (не нашли позиции):
      - Проверить наличие накопительной стоимости (accumulatedCost)
      - ЕСЛИ есть позиции с accumulatedCost != null:
        * Конкатенировать их ID
      - ИНАЧЕ:
        * Пропустить конкатенацию
   
   d) Найти общую сумму найденных позиций:
      - Суммировать currentPeriodCost
      - Суммировать currentPeriodCostInspection
      - Суммировать currentPeriodCostConsiderations
   
   e) Записать общие суммы в bigAct:
      - bigAct.totalCost
      - bigAct.totalCostInspection
      - bigAct.totalCostConsiderations
      - bigAct.positionIDs
   
   f) Проставить текстовые значения в bigAct
   
   g) Обновить акт в MongoDB
   
   h) ЕСЛИ успешно сохранили:
      - Создать Excel-документ через ExcelService
      - Сгенерировать имя файла (например: act_{id}_{timestamp}.xlsx)
      - Сохранить файл в generated/
      - Обновить bigAct.bigActLink = "/api/act/download/{filename}"
      - Обновить акт в БД
      - Вернуть bigActLink
   
   i) ИНАЧЕ:
      - Вернуть ошибку
   
   ИНАЧЕ (bigActChanged == false):
   
   j) Вернуть существующий bigAct.bigActLink
```

**Вспомогательные функции:**
- `findPositionsWithCurrentPeriod()`: Фильтрация позиций
- `findPositionsWithAccumulated()`: Фильтрация по накопительной стоимости
- `concatenatePositionIDs()`: Конкатенация ID
- `calculateTotals()`: Суммирование стоимостей
- `buildTemplateData()`: Построение map для подстановки в шаблон

### Этап 8: HTTP Handlers

**Файл internal/handlers/act_handler.go:**

**1. POST /api/act/create**
```go
func (h *ActHandler) CreateAct(c *gin.Context) {
    // 1. Парсинг JSON из request body
    // 2. Валидация данных
    // 3. Вызов actService.CreateAct()
    // 4. Возврат JSON ответа:
    //    Success: {"id": "..."}
    //    Error: {"error": "..."}
}
```

**Request body пример:**
```json
{
  "bigAct": {
    "changed": true,
    "textFields": {
      "contractNumber": "12345",
      "contractDate": "01.01.2025",
      "customer": "ООО Рога и Копыта"
    }
  },
  "positions": [
    {
      "currentPeriodCost": 1000000.50,
      "currentPeriodCostInspection": 50000.25,
      "currentPeriodCostConsiderations": 25000.10
    }
  ]
}
```

**Response пример:**
```json
{
  "id": "507f1f77bcf86cd799439011"
}
```

**2. GET /api/act/generate**
```go
func (h *ActHandler) GenerateAct(c *gin.Context) {
    // 1. Получить ID из query параметра (?id=...)
    // 2. Валидация ID
    // 3. Вызов actService.GenerateAct(id)
    // 4. Возврат JSON ответа:
    //    Success: {"downloadLink": "/api/act/download/..."}
    //    Error: {"error": "..."}
}
```

**Request:**
```
GET /api/act/generate?id=507f1f77bcf86cd799439011
```

**Response пример:**
```json
{
  "downloadLink": "/api/act/download/act_507f1f77bcf86cd799439011_1730739600.xlsx"
}
```

**3. GET /api/act/download/:filename** (дополнительный endpoint)
```go
func (h *ActHandler) DownloadAct(c *gin.Context) {
    // 1. Получить filename из URL параметра
    // 2. Проверить существование файла
    // 3. Установить headers для скачивания
    // 4. Отдать файл
}
```

### Этап 9: Главный файл приложения

**Файл cmd/server/main.go:**

```go
func main() {
    // 1. Загрузка конфигурации
    cfg := config.Load()
    
    // 2. Подключение к MongoDB
    mongoClient := repository.ConnectMongoDB(cfg)
    defer mongoClient.Disconnect()
    
    // 3. Инициализация репозиториев
    actRepo := repository.NewActRepository(mongoClient, cfg)
    
    // 4. Инициализация сервисов
    excelService := services.NewExcelService(cfg)
    actService := services.NewActService(actRepo, excelService, cfg)
    
    // 5. Инициализация handlers
    actHandler := handlers.NewActHandler(actService)
    
    // 6. Настройка Gin router
    router := gin.Default()
    
    // 7. Регистрация маршрутов
    api := router.Group("/api")
    {
        act := api.Group("/act")
        {
            act.POST("/create", actHandler.CreateAct)
            act.GET("/generate", actHandler.GenerateAct)
            act.GET("/download/:filename", actHandler.DownloadAct)
        }
    }
    
    // 8. Статические файлы (опционально)
    router.Static("/generated", cfg.GeneratedPath)
    
    // 9. Запуск сервера
    router.Run(":" + cfg.ServerPort)
}
```

### Этап 10: Docker конфигурация

**Dockerfile:**
```dockerfile
# Многоступенчатая сборка для оптимизации размера образа

# Stage 1: Build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Stage 2: Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
RUN mkdir -p ./generated
EXPOSE 8080
CMD ["./server"]
```

**docker-compose.yml:**
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - SERVER_PORT=8080
      - MONGODB_URI=mongodb://mongodb:27017
      - MONGODB_DATABASE=acts_db
      - MONGODB_COLLECTION=acts
      - TEMPLATE_PATH=./templates/act_template.xlsx
      - GENERATED_PATH=./generated
      - BASE_URL=http://localhost:8080
    volumes:
      - ./generated:/root/generated
      - ./templates:/root/templates
    depends_on:
      - mongodb
    restart: unless-stopped

  mongodb:
    image: mongo:7.0
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    restart: unless-stopped

  mongo-express:
    image: mongo-express:latest
    ports:
      - "8081:8081"
    environment:
      - ME_CONFIG_MONGODB_URL=mongodb://mongodb:27017
      - ME_CONFIG_BASICAUTH_USERNAME=admin
      - ME_CONFIG_BASICAUTH_PASSWORD=admin
    depends_on:
      - mongodb
    restart: unless-stopped

volumes:
  mongodb_data:
```

### Этап 11: Excel Template

**Создание шаблона templates/act_template.xlsx:**

Шаблон должен содержать ячейки с ключами в формате `{{key}}`:

Примеры ключей:
- `{{contractNumber}}` - Номер договора
- `{{contractDate}}` - Дата договора
- `{{customer}}` - Заказчик
- `{{totalCost}}` - Общая стоимость
- `{{totalCostInspection}}` - Стоимость инспекции
- `{{totalCostConsiderations}}` - Стоимость рассмотрения
- `{{positionIds}}` - ID позиций
- `{{currentDate}}` - Текущая дата
- И другие по необходимости

**Структура шаблона:**
- Может содержать несколько листов
- Ключи размещаются в любых ячейках
- Поддержка форматирования Excel (шрифты, границы, и т.д.)
- Числовые значения будут отформатированы автоматически

### Этап 12: Обработка ошибок и валидация

**Типы ошибок:**
- Ошибки валидации входных данных
- Ошибки подключения к БД
- Ошибки при работе с файлами
- Ошибки при парсинге Excel

**Структура ответов об ошибках:**
```go
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
    Code    int    `json:"code"`
}
```

**HTTP коды:**
- 200: Успешная операция
- 201: Создан новый ресурс
- 400: Неверные входные данные
- 404: Ресурс не найден
- 500: Внутренняя ошибка сервера

### Этап 13: Логирование

**Использование стандартного log пакета или zerolog:**
- Логирование запросов
- Логирование ошибок
- Логирование операций с БД
- Логирование генерации файлов

**Уровни логирования:**
- INFO: Информационные сообщения
- WARN: Предупреждения
- ERROR: Ошибки
- DEBUG: Отладочная информация

### Этап 14: Тестирование

**Unit тесты:**
- `number_formatter_test.go`: Тестирование форматирования чисел
- `act_service_test.go`: Тестирование бизнес-логики
- `excel_service_test.go`: Тестирование генерации Excel

**Integration тесты:**
- Тестирование API endpoints
- Тестирование работы с MongoDB (с использованием testcontainers)

**Запуск тестов:**
```bash
go test ./... -v
go test ./... -cover
```

### Этап 15: Документация

**README.md:**
- Описание проекта
- Требования
- Установка и запуск
- API документация
- Примеры использования
- Структура проекта

**API документация:**
- Swagger/OpenAPI спецификация (опционально)
- Примеры curl запросов
- Описание моделей данных

## Порядок разработки

### Фаза 1: Инфраструктура (1-2 дня)
1. ✅ Инициализация Go проекта
2. ✅ Создание структуры директорий
3. ✅ Настройка Docker и Docker Compose
4. ✅ Создание базовой конфигурации

### Фаза 2: Модели и репозиторий (1-2 дня)
5. ✅ Определение моделей данных
6. ✅ Реализация MongoDB подключения
7. ✅ Реализация ActRepository
8. ✅ Тестирование работы с БД

### Фаза 3: Утилиты и вспомогательные сервисы (1 день)
9. ✅ Реализация форматирования чисел
10. ✅ Создание утилит для HTTP ответов
11. ✅ Unit тесты для утилит

### Фаза 4: Excel сервис (2-3 дня)
12. ✅ Создание Excel шаблона
13. ✅ Реализация ExcelService
14. ✅ Парсинг и замена ключей
15. ✅ Обработка числовых значений
16. ✅ Тестирование генерации

### Фаза 5: Бизнес-логика (2-3 дня)
17. ✅ Реализация ActService
18. ✅ Реализация метода CreateAct
19. ✅ Реализация метода GenerateAct
20. ✅ Реализация алгоритма обработки bigAct
21. ✅ Тестирование бизнес-логики

### Фаза 6: HTTP Layer (1-2 дня)
22. ✅ Реализация HTTP handlers
23. ✅ Настройка Gin router
24. ✅ Реализация endpoint /api/act/create
25. ✅ Реализация endpoint /api/act/generate
26. ✅ Реализация endpoint /api/act/download/:filename

### Фаза 7: Интеграция и тестирование (2-3 дня)
27. ✅ Интеграция всех компонентов
28. ✅ End-to-end тестирование
29. ✅ Тестирование в Docker окружении
30. ✅ Исправление багов

### Фаза 8: Документация и финализация (1 день)
31. ✅ Написание README.md
32. ✅ API документация
33. ✅ Создание примеров использования
34. ✅ Финальное тестирование

## Потенциальные проблемы и решения

### Проблема 1: Производительность при работе с большими Excel файлами
**Решение:**
- Использование streaming API Excelize для больших файлов
- Оптимизация поиска ячеек с ключами
- Кэширование шаблонов

### Проблема 2: Конкурентный доступ к файлам
**Решение:**
- Использование уникальных имен файлов (ID + timestamp)
- Mutex для критических секций
- Очистка старых файлов (cron job)

### Проблема 3: Обработка различных типов данных в Excel
**Решение:**
- Определение типа через reflection
- Специальные обработчики для дат, чисел, строк
- Fallback на строковое представление

### Проблема 4: Масштабирование
**Решение:**
- Горизонтальное масштабирование через Docker Swarm/Kubernetes
- Использование MongoDB replica set
- Хранение сгенерированных файлов в S3-совместимом хранилище

## Зависимости (go.mod)

```go
module github.com/yourname/acts-service

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/xuri/excelize/v2 v2.8.0
    go.mongodb.org/mongo-driver v1.13.1
    github.com/joho/godotenv v1.5.1
    github.com/go-playground/validator/v10 v10.15.5
)
```

## Переменные окружения

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
BASE_URL=http://localhost:8080

# MongoDB Configuration
MONGODB_URI=mongodb://mongodb:27017
MONGODB_DATABASE=acts_db
MONGODB_COLLECTION=acts
MONGODB_TIMEOUT=10s

# File Paths
TEMPLATE_PATH=./templates/act_template.xlsx
GENERATED_PATH=./generated

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# File Cleanup (опционально)
CLEANUP_ENABLED=true
CLEANUP_INTERVAL=24h
FILE_RETENTION_DAYS=7
```

## Команды для разработки

**Локальная разработка:**
```bash
# Установка зависимостей
go mod download

# Запуск в dev режиме
go run cmd/server/main.go

# Запуск тестов
go test ./... -v

# Запуск с покрытием
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out
```

**Docker разработка:**
```bash
# Сборка и запуск
docker-compose up --build

# Запуск в фоне
docker-compose up -d

# Просмотр логов
docker-compose logs -f app

# Остановка
docker-compose down

# Полная очистка
docker-compose down -v
```

## API Примеры использования

### 1. Создание акта
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
      }
    ]
  }'
```

**Ответ:**
```json
{
  "id": "507f1f77bcf86cd799439011"
}
```

### 2. Генерация акта
```bash
curl -X GET "http://localhost:8080/api/act/generate?id=507f1f77bcf86cd799439011"
```

**Ответ:**
```json
{
  "downloadLink": "/api/act/download/act_507f1f77bcf86cd799439011_1730739600.xlsx"
}
```

### 3. Скачивание файла
```bash
curl -O "http://localhost:8080/api/act/download/act_507f1f77bcf86cd799439011_1730739600.xlsx"
```

## Метрики успеха

- ✅ Сервис успешно поднимается в Docker
- ✅ API endpoints отвечают корректно
- ✅ Данные корректно сохраняются в MongoDB
- ✅ Excel файлы генерируются правильно
- ✅ Числа форматируются в нужном формате
- ✅ Все unit тесты проходят
- ✅ Покрытие тестами > 70%
- ✅ Документация полная и актуальная

## Следующие шаги (после базовой реализации)

1. **Безопасность:**
   - Добавить аутентификацию (JWT)
   - Добавить авторизацию
   - Rate limiting
   - CORS настройки

2. **Мониторинг:**
   - Prometheus метрики
   - Health check endpoints
   - Grafana дашборды

3. **Оптимизация:**
   - Кэширование
   - Connection pooling
   - Batch operations

4. **Расширенный функционал:**
   - Поддержка нескольких шаблонов
   - Версионирование актов
   - История изменений
   - Email уведомления
   - Webhook интеграции

---

**Общая оценка времени:** 10-15 дней разработки

**Приоритет:** Сначала реализовать core функционал (Фазы 1-6), затем тестирование и документация.

