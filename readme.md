# Finance Tracker — Курсовой проект

## Дисциплина
**Проектирование и разработка клиент-серверных приложений**

---

##  Описание проекта

Finance Tracker — это полнофункциональное клиент-серверное приложение для учета личных финансов. Система позволяет пользователю вести учет доходов и расходов, управлять категориями, устанавливать бюджеты и просматривать аналитическую сводку по финансовым операциям.

Проект реализован как распределённое приложение с разделением на:

- backend (Go API сервис)
- frontend (React SPA)
- database (PostgreSQL)

---

##  Архитектура backend

Backend написан на языке **Go** с применением **гексагональной архитектуры (Hexagonal / Ports & Adapters)**.

### Основная идея архитектуры

Бизнес-логика полностью изолирована от внешних зависимостей (HTTP, БД, внешние сервисы).

### Слои приложения

- Domain layer
  - Сущности: User, Transaction, Category, Budget
  - Бизнес-правила

- Use Cases (Application layer)
  - регистрация пользователя
  - авторизация
  - создание транзакций
  - управление бюджетами

- Adapters
  - HTTP handlers (REST API)
  - PostgreSQL repository

- Infrastructure
  - подключение к базе данных
  - конфигурация окружения

Такой подход позволяет легко тестировать бизнес-логику и заменять инфраструктурные компоненты без изменения ядра системы.

---

##  Авторизация

В системе реализована JWT-аутентификация.

### Механизм

1. Пользователь регистрируется (email + password)
2. Пароль хэшируется (bcrypt)
3. При логине проверяются учетные данные
4. Сервер выдает JWT токен
5. Все защищённые API требуют:

Authorization: Bearer <token>

### Особенности

- Stateless архитектура
- Токен содержит userID
- Middleware проверяет валидность токена

---

##  Работа с базой данных

Используется PostgreSQL.

Миграции выполняются через Goose:

```bash
goose -dir /app/migrations postgres $DB_DSN up

Пример миграции
-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
```

## Frontend

React SPA приложение.

### Функциональность

- авторизация / регистрация
- просмотр финансов
- транзакции
- категории
- бюджеты

### Особенности

- React Hooks
- REST API
- JWT в localStorage
- автообновление данных каждые 30 сек

---

##  API

POST /auth/login
POST /auth/register
GET /api/transactions
POST /api/transactions
GET /api/categories
GET /api/budgets
GET /api/analytics/summary

---

##  Docker

### Сервисы

- PostgreSQL 16
- Backend (Go)
- Frontend (React)

### Особенности

- multi-stage build
- Alpine runtime
- изолированная сеть

---

## Миграции при старте

- Поднимается PostgreSQL
- Backend ждёт готовности БД
- Запускается goose migration
- Стартует сервер

---

##  Сборка

```bash
go build -o main ./cmd/app
npm install
npm run build
docker compose up --build
```

##  ENV

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=finance

DB_DSN=postgres://postgres:postgres@postgres:5432/finance?sslmode=disable

REACT_APP_API_URL=http://localhost:8080
```

##  Запуск

```bash
git clone https://github.com/liftoach/kr_rksp.git
cd finance-tracker
docker compose up --build
```


# Finance Tracker

##  Доступ

- Frontend: http://localhost:3000  
- Backend: http://localhost:8080  
- DB: localhost:5432  

---

##  Особенности

- Clean Architecture (Hexagonal)
- Stateless JWT auth
- Dockerized
- Auto migrations
- REST API

---

##  Итог

Проект демонстрирует:

- современную гексоганальную backend архитектуру на Go
- SPA frontend на React
- контейнеризацию
- миграции БД
- REST API
- безопасную авторизацию
- запуск проекта `docker compose up --build`