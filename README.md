# Golang Shop RESTful API

Простое RESTful API для интернет-магазина, написанное на Go с использованием фреймворка Gin.

## Описание

Этот проект предоставляет базовый функционал для управления продуктами и пользователями в интернет-магазине:

- CRUD операции для продуктов
- Регистрация и аутентификация пользователей
- JWT аутентификация
- Пагинация для списка продуктов

## Требования

- Go 1.20+
- PostgreSQL
- Git

## Установка

1. Клонируйте репозиторий:

```bash
git clone https://github.com/kavlan-dev/golang-shop-restful.git
cd golang-shop-restful
```

2. Установите зависимости:

```bash
go mod download
```

## Конфигурация

1. **Шаблон конфигурации**: В проекте доступен шаблон конфигурационного файла `config/config.example.yaml`. Скопируйте его в `config/config.yaml` и настройте параметры:

```bash
cp config/config.example.yaml config/config.yaml
```

2. **Структура конфигурационного файла**:

```yaml
# Server Configuration
server:
  host: localhost
  port: 8080

# Database Configuration
database:
  host: localhost
  user: myuser
  password: pass
  name: mydb
  port: 5432

# JWT Configuration
jwt:
  secret: your-very-secure-secret-key
```

3. **Важно**: Файл `config/config.yaml` добавлен в `.gitignore`, чтобы избежать коммита чувствительных данных (паролей, секретных ключей) в репозиторий.

Приложение будет автоматически загружать конфигурацию из файла `config/config.yaml` при запуске.

## Запуск

```bash
go run main.go
```

Сервер будет запущен на порту `8080` по умолчанию.

## API Эндпоинты

### Аутентификация

- **POST** `/api/auth/register` - Регистрация нового пользователя
- **POST** `/api/auth/login` - Аутентификация пользователя и получение JWT токена

### Продукты

Все эндпоинты для продуктов требуют JWT аутентификации (передавайте токен в заголовке `Authorization: Bearer <token>`):

- **GET** `/api/products` - Получение списка продуктов (с поддержкой пагинации)
  - Параметры: `limit` (по умолчанию: 100), `offset` (по умолчанию: 0)
- **POST** `/api/products` - Создание нового продукта
- **GET** `/api/products/:id` - Получение продукта по ID
- **PUT** `/api/products/:id` - Обновление продукта
- **DELETE** `/api/products/:id` - Удаление продукта

## Примеры использования

### Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass",
    "email": "test@example.com"
  }'
```

### Аутентификация пользователя

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "testpass"
  }'
```

### Создание продукта

```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Новый продукт",
    "price": 1000
  }'
```

### Получение списка продуктов

```bash
curl -X GET "http://localhost:8080/api/products?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Структура проекта

```
.
├── main.go                  # Точка входа
├── internal/                # Конфигурации проекта
├── internal/
│   ├── config/              # Применение конфигураций
│   ├── database/            # Подключение к базе данных
│   ├── handlers/            # Обработчики HTTP запросов
│   ├── middleware/          # Middleware (аутентификация)
│   ├── models/              # Модели данных
│   ├── services/            # Бизнес-логика
│   └── utils/               # Утилиты (JWT)
├── go.mod                   # Модуль Go
└── go.sum                   # Контрольные суммы зависимостей
```

## Технологии

- **Фреймворк**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Логирование**: [Zap](https://github.com/uber-go/zap)
- **JWT**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- **Конфигурация**: [Viper](https://github.com/spf13/viper)
- **База данных**: PostgreSQL

## Лицензия

Этот проект лицензирован по лицензии MIT. См. файл [LICENSE](LICENSE) для получения дополнительной информации.
