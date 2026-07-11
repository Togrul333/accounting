# Accounting API

REST API для бухгалтерского учёта. Построен на Go + Gin + MySQL.

1) rengler #00d4ff  #7b2fff  для карточек и попапа задний фон rgb(99, 218.2, 159.8) 

2) css ler hamsi lokaldan sdn yox
 
3) gorm lazimdi sql yox 

4) logaout

5) 

## Стек

- **Go** — язык
- **Gin** — HTTP фреймворк
- **MySQL** — база данных
- **database/sql** — работа с БД без ORM

## Структура проекта

```
accounting/
├── cmd/
│   ├── api/
│   │   └── main.go              # Точка входа: инициализация БД, DI, запуск сервера
│   └── seed/
│       └── main.go              # Сидер: создаёт тестового пользователя
├── internal/
│   ├── model/
│   │   └── account.go           # Структуры: Account, CreateAccountRequest, UpdateAccountRequest
│   ├── repository/
│   │   └── account.go           # Интерфейс + SQL-запросы к MySQL
│   ├── service/
│   │   └── account.go           # Бизнес-логика
│   └── handler/
│       ├── account.go           # Gin-хендлеры (GET, POST, PUT, DELETE)
│       ├── page.go              # Хендлеры HTML-страниц
│       └── router.go            # Регистрация маршрутов
├── migrations/
│   ├── 001_create_accounts.sql  # Таблица accounts
│   └── 002_create_users.sql     # Таблица users
├── web/
│   └── templates/               # HTML шаблоны (Tailwind CDN)
├── go.mod
├── go.sum
└── CLAUDE.md
```

## Установка и запуск

### 1. Клонировать репозиторий

```bash
git clone <repo-url>
cd accounting
```

### 2. Установить зависимости

```bash
go mod tidy
```

### 3. Создать базу данных

```sql
CREATE DATABASE accounting CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. Применить миграции

```bash
# Если MySQL от XAMPP
/Applications/XAMPP/xamppfiles/bin/mysql -u root -p accounting < migrations/001_create_accounts.sql
/Applications/XAMPP/xamppfiles/bin/mysql -u root -p accounting < migrations/002_create_users.sql

# Если MySQL установлен глобально
mysql -u root -p accounting < migrations/001_create_accounts.sql
mysql -u root -p accounting < migrations/002_create_users.sql
```

### 5. Запустить сидер (тестовый пользователь)

```bash
DB_USER=root DB_PASSWORD=secret DB_NAME=accounting go run ./cmd/seed
```

#### Тестовый пользователь

| Поле    | Значение              |
|---------|-----------------------|
| E-posta | admin@hisartour.az    |
| Şifre   | admin123              |

### 6. Запустить сервер

```bash
DB_USER=root DB_PASSWORD=secret DB_NAME=accounting go run ./cmd/api
```

Сервер поднимется на `http://localhost:8080`.

## Переменные окружения

| Переменная    | По умолчанию  | Описание          |
|---------------|---------------|-------------------|
| `DB_HOST`     | `localhost`   | MySQL хост        |
| `DB_PORT`     | `3306`        | MySQL порт        |
| `DB_USER`     | `root`        | MySQL пользователь|
| `DB_PASSWORD` | _(пусто)_     | MySQL пароль      |
| `DB_NAME`     | `accounting`  | Название БД       |
| `PORT`        | `8080`        | Порт HTTP сервера |
| `GOOGLE_CREDENTIALS_PATH` | `./credentials.json` | Путь к JSON-ключу сервисного аккаунта Google (для импорта из Google Sheets) |

## Импорт из Google Sheets

Service Account для чтения таблиц: `sheets-import@liquid-journal-454008-b1.iam.gserviceaccount.com`

Чтобы дать доступ к новой таблице — открыть её в Google Sheets → **Share** → добавить этот email с правами **Viewer**.

JSON-ключ сервисного аккаунта лежит в `credentials.json` в корне проекта (в `.gitignore`, не коммитится). Если нужно перевыпустить ключ — Google Cloud Console → проект `liquid-journal-454008-b1` → APIs & Services → Credentials → Service Accounts → `sheets-import` → Keys → Add Key.

## API — Банковские счета

Base URL: `/api/accounts`

### Получить все счета

```
GET /api/accounts
```

**Ответ `200`:**
```json
[
  {
    "id": 1,
    "name": "Основной счёт",
    "account_number": "AZ12NABZ00000000137010001944",
    "currency": "AZN",
    "balance": 15000.00,
    "description": "Расчётный счёт в Kapital Bank",
    "created_at": "2026-04-29T10:00:00Z",
    "updated_at": "2026-04-29T10:00:00Z"
  }
]
```

---

### Получить счёт по ID

```
GET /api/accounts/:id
```

**Ответ `200`:** объект счёта  
**Ответ `404`:** `{"error": "account not found"}`

---

### Создать счёт

```
POST /api/accounts
Content-Type: application/json
```

**Тело запроса:**
```json
{
  "name": "Основной счёт",
  "account_number": "AZ12NABZ00000000137010001944",
  "currency": "AZN",
  "balance": 15000.00,
  "description": "Расчётный счёт в Kapital Bank"
}
```

**Ответ `201`:** созданный объект счёта

---

### Обновить счёт

```
PUT /api/accounts/:id
Content-Type: application/json
```

**Тело запроса:** те же поля, что и при создании  
**Ответ `200`:** обновлённый объект счёта  
**Ответ `404`:** `{"error": "account not found"}`

---

### Удалить счёт

```
DELETE /api/accounts/:id
```

**Ответ `204`:** нет тела

---

## Модель данных

### Account

| Поле             | Тип            | Описание                          |
|------------------|----------------|-----------------------------------|
| `id`             | `int64`        | Первичный ключ, автоинкремент     |
| `name`           | `string`       | Название счёта                    |
| `account_number` | `string`       | Номер счёта (уникальный)          |
| `currency`       | `string`       | Валюта (по умолчанию `AZN`)       |
| `balance`        | `float64`      | Баланс                            |
| `description`    | `string`       | Описание (необязательное)         |
| `created_at`     | `time.Time`    | Дата создания                     |
| `updated_at`     | `time.Time`    | Дата последнего обновления        |

## Архитектура

Классическая четырёхслойная архитектура:

```
Handler  →  Service  →  Repository  →  MySQL
```

- **Handler** — принимает HTTP-запрос, парсит параметры, возвращает JSON
- **Service** — содержит бизнес-логику, вызывает репозиторий
- **Repository** — все SQL-запросы, реализует интерфейс
- **Model** — чистые структуры без логики

Зависимость идёт через интерфейс `AccountRepository`, что позволяет легко подменять реализацию.

## Сборка бинарника

```bash
go build -o bin/api ./cmd/api
./bin/api
```