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

Миграции применяются строго по порядку номеров. Для существующей базы:

```bash
for f in migrations/0*.sql; do
  /Applications/XAMPP/xamppfiles/bin/mysql -u root accounting < "$f" || break
done
```

#### Комнаты у туров (миграции 032–036)

Тур связан с комнатами через пивот `tour_rooms(tour_id, room_id, price)`, а не одной
колонкой `tours.room_id`. Что важно знать:

- **Цена номера хранится в пивоте, а не в `rooms`.** `rooms` — справочник типов
  (`code` + `beds_count`); один и тот же тип может стоить по-разному в разных турах.
- **`orders.room_id`** — заказ выбирает одну комнату из комнат своего тура. NULL допустим:
  заказы из банковской выписки создаются без комнаты, тогда цена номера равна 0.
- **Миграция 034 схлопывает дубли туров.** Раньше импорт из Sheets создавал отдельный тур
  на каждую пару (комната × категория), поэтому в базе лежали туры с одинаковым кодом.
  Они сливаются в один по ключу (код, категория, даты); заказы, доходы и расходы
  перевешиваются на выживший тур. Перед прогоном сделайте дамп:

  ```bash
  /Applications/XAMPP/xamppfiles/bin/mysqldump -u root accounting > accounting_backup.sql
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
| `META_ACCESS_TOKEN` | _(пусто)_ | System user token Meta Ads (fallback, если у кабинета не задан свой) |
| `META_AD_ACCOUNT_ID` | _(пусто)_ | Рекламный кабинет Meta по умолчанию, формат `act_1234567890` |
| `META_API_VERSION` | `v23.0` | Версия Graph API |

## Импорт из Google Sheets

Service Account для чтения таблиц: `sheets-import@liquid-journal-454008-b1.iam.gserviceaccount.com`

Чтобы дать доступ к новой таблице — открыть её в Google Sheets → **Share** → добавить этот email с правами **Viewer**.

JSON-ключ сервисного аккаунта лежит в `credentials.json` в корне проекта (в `.gitignore`, не коммитится). Если нужно перевыпустить ключ — Google Cloud Console → проект `liquid-journal-454008-b1` → APIs & Services → Credentials → Service Accounts → `sheets-import` → Keys → Add Key.

## Meta Ads (Facebook / Instagram)

Чтение расходов по рекламным кампаниям и автоматическое создание записей в `expenses`.
Страница — `/meta-ads` (в меню `Araçlar → Meta Reklamlar`).

Полная инструкция (как создать приложение, system user и бессрочный токен, схема
таблиц, список эндпоинтов): **[docs/meta-ads-integration.md](docs/meta-ads-integration.md)**

Коротко:

```bash
# 1. Миграция
mysql -u root -p accounting < migrations/037_create_meta_ads.sql

# 2. Проверить токен (до настройки в UI)
curl -G "https://graph.facebook.com/v23.0/act_ВАШ_ID/insights" \
  -d "fields=campaign_name,spend" -d "date_preset=last_30d" \
  -d "level=campaign" -d "access_token=ВАШ_ТОКЕН"
```

Токен — **system user token** из Business Manager с правами `ads_read` + `read_insights`
и сроком действия **«Никогда»**. App Review не нужен, пока читаем свой кабинет.
Хранится либо в `META_ACCESS_TOKEN`, либо в БД (таблица `meta_ad_accounts`, вводится через UI).

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
