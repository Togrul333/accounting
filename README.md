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

## Referans kullanıcıları (рефереры)

Люди, которые приводят клиентов. Хранится минимум: имя, фамилия, телефон.
Управление — вкладка **Referanslar** на странице `/settings`.

```bash
# Миграции
mysql -u root -p accounting < migrations/038_create_referans_users.sql
mysql -u root -p accounting < migrations/039_create_referans_user_orders.sql
```

Детальная страница — `/referans-users/:id` (кнопка с глазом рядом с карандашом в списке):
данные реферера, заказы-кандидаты, поиск по всем заказам и список подтверждённых рефералов.

**Кандидаты** подбираются автоматически: `clients.reference_name` содержит имя или фамилию
реферера (`LIKE`, части короче 3 символов игнорируются). По умолчанию — последние 10.
Само подтверждение ручное и хранится в таблице `referans_user_orders`.

Эндпоинты (`/api/referans-users`):

| Метод | Путь | Описание |
|---|---|---|
| `GET` / `POST` | `/` | список / создать |
| `GET` / `PUT` / `DELETE` | `/:id` | карточка реферера |
| `GET` | `/:id/candidates?limit=10` | заказы-кандидаты (`limit=all` — все) |
| `GET` | `/:id/order-search?q=&limit=30` | поиск по всем заказам (клиент, референс, телефон, код тура, номер заказа) |
| `GET` | `/:id/orders` | подтверждённые рефералы |
| `POST` | `/:id/orders` | подтвердить (`{"order_id": 9}`), идемпотентно |
| `DELETE` | `/:id/orders/:order_id` | снять подтверждение |

Поля JSON реферера: `first_name` (обязательное), `last_name`, `phone`.

## Hocalar (персонал агентства)

Сотрудники агентства: имя, фамилия, телефон. На них назначаются задачи.
Управление — вкладка **Hocalar** на странице `/settings`.

```bash
# Миграция (создаёт hoca_users и добавляет tasks.hoca_user_id)
mysql -u root -p accounting < migrations/040_create_hoca_users.sql
```

Эндпоинты `/api/hoca-users`: `GET` список, `POST` создать, `GET /:id`, `PUT /:id`, `DELETE /:id`.
Поля JSON: `first_name` (обязательное), `last_name`, `phone`.

**Задачи:** в модалке создания/редактирования на `/tasks` есть селекты «Atanan Kişi» и «Tur».
В запросах `POST`/`PUT /api/tasks` за это отвечают `hoca_user_id` и `tour_id` (`null` — не выбрано),
в ответе дополнительно приходят `hoca_user_name` и `tour_code`, оба выводятся бейджами на карточке.
При удалении сотрудника задачи сохраняются, назначение обнуляется (`ON DELETE SET NULL`);
при удалении тура его задачи удаляются вместе с ним (`ON DELETE CASCADE`).

## Varsayılan görevler (шаблоны задач тура)

Задачи, которые автоматически создаются при создании тура и привязываются к нему.
Управление — вкладка **Varsayılan Görevler** на странице `/settings`.

```bash
# Миграция (создаёт default_tasks и добавляет tasks.tour_id)
mysql -u root -p accounting < migrations/041_create_default_tasks.sql
```

Поля шаблона: `title` (обязательное), `description`, `days_before_start`, `hoca_user_id`.
При `POST /api/tours` из каждого шаблона создаётся задача со статусом `todo`,
`tour_id` нового тура, унаследованным исполнителем и сроком
**`due_date = start_date − days_before_start`** (0 — день старта тура).
Ошибка при создании задач пишется в лог и не отменяет создание тура.

Эндпоинты `/api/default-tasks`: `GET` список, `POST` создать, `GET /:id`, `PUT /:id`, `DELETE /:id`.
Удаление шаблона не трогает уже созданные задачи.

## Görevler — доска задач (`/tasks`)

Канбан-доска в стиле Jira: три колонки (`todo` / `in_progress` / `done`),
drag & drop между колонками, поиск, фильтры по исполнителю, приоритету и просрочке.
Доска рисуется на клиенте из JSON, отданного шаблоном, — после каждой правки
перерисовывается без перезагрузки страницы.

```bash
# Миграция (приоритет, связи с заказом/клиентом, таблица комментариев)
mysql -u root -p accounting < migrations/042_task_links_and_comments.sql
```

**Связи задачи.** Кроме исполнителя (`hoca_user_id`) и тура (`tour_id`) задачу
можно привязать к заказу (`order_id`) и клиенту (`client_id`) — `null` означает
«не выбрано». В ответе дополнительно приходят подписи `hoca_user_name`, `tour_code`,
`order_label` (`#9 Имя Фамилия`), `client_name` и счётчик `comment_count`;
всё это выводится чипами на карточке. При удалении заказа или клиента задача
сохраняется, связь обнуляется (`ON DELETE SET NULL`).

**Приоритет** `priority`: `urgent` / `high` / `medium` / `low` (по умолчанию `medium`,
некорректное значение приводится к нему же). Сортировка на доске — по приоритету,
внутри приоритета по сроку.

**Попап задачи** открывается кликом по карточке на весь экран: слева заголовок,
описание и лента комментариев, справа поля (статус, приоритет, исполнитель, срок,
тур, заказ, клиент). Любое изменение поля сразу уходит в `PUT /api/tasks/:id`.

> Внимание: `PUT /api/tasks/:id` — полная замена задачи, а не частичное обновление.
> Для переноса между колонками используйте `PUT /api/tasks/:id/status` — он меняет
> только статус.

### Комментарии к задаче

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/api/tasks/:id/comments` | лента комментариев, старые сверху |
| `POST` | `/api/tasks/:id/comments` | добавить (`{"body": "текст"}`) |
| `PUT` | `/api/task-comments/:id` | изменить текст |
| `DELETE` | `/api/task-comments/:id` | удалить |

Автор берётся из сессии; в ответе приходит `user_name`. Пустой `body` — `400`.
Комментарии удаляются вместе с задачей (`ON DELETE CASCADE`).

## Telegram — уведомления о задачах

Когда задаче назначают исполнителя (`hoca_user_id`), он получает сообщение в Telegram.
Настройки — вкладка **Telegram** на странице `/settings`.

```bash
# Миграция (ключи бота в settings, привязка чата в hoca_users)
mysql -u root -p accounting < migrations/043_telegram.sql
```

### Два ограничения Telegram Bot API

1. **Бот не может написать человеку первым.** Пока сотрудник не нажмёт Start,
   `sendMessage` вернёт `403`. Поэтому у каждого hoca есть свой `telegram_chat_id`,
   и его нельзя вписать руками — он появляется только после Start.
2. **Бот не может создать группу.** В Bot API нет такого метода — группы создаёт
   только настоящий аккаунт через MTProto. Для групп по турам придётся создавать
   чат вручную, добавлять туда бота и хранить `chat_id` группы.

### Как подключить

1. `@BotFather` → `/newbot` → получить токен.
2. `/settings` → **Telegram** → вставить токен → **Kaydet**.
   Токен проверяется через `getMe`; при неверном токене настройки не сохраняются.
   Заодно запоминается username бота — он нужен для ссылки-приглашения.
3. Для каждого сотрудника — **Bağlantı linki**: выдаётся разовый код и ссылка
   `https://t.me/<bot>?start=<code>`. Ссылку отправить сотруднику, он жмёт Start.
4. **Bağlantıları eşleştir** — приложение читает `getUpdates`, находит `/start <code>`
   и привязывает `chat_id` к сотруднику; код гасится, сотруднику уходит приветствие.

Webhook не используется (у локального приложения нет внешнего URL), вместо него
`getUpdates` с offset в `settings.telegram_update_offset` — одни и те же сообщения
не разбираются дважды.

### Эндпоинты

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/api/telegram/settings` | состояние (`enabled`, `has_token`, `bot_username`); токен наружу не отдаётся |
| `PUT` | `/api/telegram/settings` | `{"bot_token": "...", "enabled": true}`; пустой токен = «оставить текущий» |
| `POST` | `/api/telegram/sync` | разобрать `/start` и привязать чаты |
| `POST` | `/api/telegram/hoca-users/:id/link-code` | код + ссылка для привязки |
| `POST` | `/api/telegram/hoca-users/:id/test` | тестовое сообщение сотруднику |
| `DELETE` | `/api/telegram/hoca-users/:id` | отвязать чат |

### Поведение уведомлений

- Отправка идёт в отдельной горутине со своим таймаутом 15 с: медленный Telegram
  не задерживает ответ API и **не отменяет** сохранение задачи.
- Уведомление уходит при `POST /api/tasks` с исполнителем и при `PUT /api/tasks/:id`,
  если исполнитель сменился. Снятие исполнителя и смена статуса ничего не шлют.
- Если интеграция выключена, токена нет или сотрудник не нажал Start — тихо пропускаем.
- Задачи из шаблонов (`default_tasks`, создаются при `POST /api/tours`) уведомлений
  **не шлют**: они пишутся напрямую через репозиторий, минуя `TaskService`.


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
