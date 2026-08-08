# Meta Ads интеграция (Facebook / Instagram)

Чтение расходов по рекламным кампаниям из Meta Marketing API и автоматическое
создание записей в `expenses`.

**Только чтение.** Права `ads_read` + `read_insights`. Создавать, менять или
останавливать кампании интеграция не умеет — и не должна, поэтому App Review
со стороны Meta не требуется.

---

## 1. Что нужно получить в Meta (делается один раз, вручную)

### Шаг 1. Business Manager

**business.facebook.com** → бизнес-аккаунт должен существовать, и рекламный
кабинет должен быть привязан к нему: `Настройки бизнеса → Аккаунты → Рекламные аккаунты`.

> Приложение и рекламный кабинет обязаны находиться в **одном** Business Manager,
> иначе токен не увидит данные.

### Шаг 2. Ad Account ID

`Настройки бизнеса → Рекламные аккаунты` → номер кабинета.
Формат: `act_1234567890` (префикс `act_` система добавляет сама, если его не ввели).

### Шаг 3. Приложение

**developers.facebook.com** → `My Apps` → **Create App**

- Тип: **Business** (в новом интерфейсе: Use case → **Other** → **Business**)
- Название: `Hisar Tour Accounting`
- Привязать к своему Business Account

`App settings → Basic` → сохранить **App ID** и **App Secret**.

### Шаг 4. Продукт Marketing API

В приложении: `Add Product` → **Marketing API** → **Set up**.

### Шаг 5. System User и бессрочный токен

Обычный пользовательский токен живёт 1–2 часа (или 60 дней после обмена) и
умрёт на сервере. Нужен **system user token — он бессрочный**.

`business.facebook.com → Настройки бизнеса → Пользователи → Системные пользователи`

1. **Добавить** → имя `accounting-api`, роль **Employee**
2. **Добавить активы** → `Рекламные аккаунты` → свой кабинет → права **Просмотр эффективности** (View Performance)
3. **Добавить активы** → `Приложения` → приложение из шага 3 → **Разработка приложения**
4. **Создать новый токен**:
   - Приложение: из шага 3
   - Срок действия: **Никогда**
   - Права: **`ads_read`**, **`read_insights`**
5. Токен показывается **один раз** — скопировать сразу

### Шаг 6. Проверка токена

```bash
curl -G "https://graph.facebook.com/v23.0/act_ВАШ_ID/insights" \
  -d "fields=campaign_name,spend,impressions,clicks" \
  -d "date_preset=last_30d" \
  -d "level=campaign" \
  -d "access_token=ВАШ_ТОКЕН"
```

| Ответ | Что значит |
|---|---|
| JSON с кампаниями | всё работает |
| `(#200) Requires ads_read permission` | права не выданы — шаг 5.4 |
| `Unsupported get request` | неверный Ad Account ID или system user не получил доступ к кабинету — шаг 5.2 |
| `Error validating access token` | токен истёк или отозван — сгенерировать заново |

---

## 2. Установка

### Миграция

```bash
mysql -u root -p accounting < migrations/037_create_meta_ads.sql
```

Создаёт `meta_ad_accounts`, `meta_ad_spend` и категорию расходов `Reklam`.

### Переменные окружения (опционально)

Токен можно хранить либо в БД (через UI), либо в `.env` — тогда он используется
как общий для всех кабинетов, у которых свой токен не задан.

| Переменная | По умолчанию | Описание |
|---|---|---|
| `META_ACCESS_TOKEN` | — | system user token; fallback, если у кабинета не задан свой |
| `META_AD_ACCOUNT_ID` | — | кабинет по умолчанию, `act_...` |
| `META_API_VERSION` | `v23.0` | версия Graph API |

---

## 3. Как пользоваться

Страница **`/meta-ads`** (в меню: `Araçlar → Meta Reklamlar`).

1. **Reklam Hesabı Ekle** — ввести Ad Account ID, токен, выбрать категорию
   расхода (`Reklam`) и счёт списания.
2. **Test Et** — проверка токена; при успехе имя и валюта кабинета подтягиваются из Meta.
3. Выбрать период → **Meta'dan Çek** — тянет дневную статистику по кампаниям.
4. Если у кабинета включён **«Harcamaları otomatik olarak giderlere aktar»**,
   для каждой строки создаётся запись в `expenses`.

### Повторный запуск безопасен

- Статистика пишется через `ON DUPLICATE KEY UPDATE` по ключу
  `(ad_account_id, campaign_id, date)` — дублей строк не будет.
- Расход создаётся только там, где `meta_ad_spend.expense_id` пуст, и
  помечается `bank_ref = meta:<кабинет>:<кампания>:<дата>` — тот же механизм
  защиты от дублей, что и при импорте банковских выписок.

---

## 4. Архитектура

Стандартная схема проекта: **handler → service → repository → MySQL**,
плюс отдельный пакет-клиент внешнего API (как `internal/googlesheets`).

```
internal/meta/client.go          — HTTP-клиент Graph API (insights, account info, пагинация)
internal/model/meta_ads.go       — MetaAdAccount, MetaAdSpend, MetaSyncResult
internal/repository/meta_ads.go  — MetaAdAccountRepository, MetaAdSpendRepository (upsert)
internal/service/meta_ads.go     — синхронизация, создание расходов
internal/handler/meta_ads.go     — gin-хендлеры
web/templates/meta_ads.html      — страница
migrations/037_create_meta_ads.sql
```

### Таблицы

**`meta_ad_accounts`** — настройки кабинета

| Колонка | Описание |
|---|---|
| `ad_account_id` | `act_1234567890`, уникален |
| `access_token` | system user token; пусто → берётся `META_ACCESS_TOKEN` |
| `expense_category_id` | категория для создаваемых расходов |
| `account_id` | счёт списания |
| `auto_expense` | создавать ли расходы автоматически |
| `last_synced_at` | время последней синхронизации |

**`meta_ad_spend`** — расход по кампании за день

| Колонка | Описание |
|---|---|
| `campaign_id`, `campaign_name`, `date` | ключ строки вместе с `ad_account_id` |
| `spend`, `impressions`, `clicks`, `reach`, `cpc`, `ctr` | метрики из Meta |
| `tour_id` | ручная привязка кампании к туру |
| `expense_id` | созданная запись расхода (NULL = ещё не списано) |

### API

| Метод | Путь | Описание |
|---|---|---|
| `GET` | `/api/meta-ads/accounts` | список кабинетов |
| `POST` | `/api/meta-ads/accounts` | добавить кабинет |
| `GET` | `/api/meta-ads/accounts/:id` | один кабинет |
| `PUT` | `/api/meta-ads/accounts/:id` | изменить (пустой `access_token` = не менять) |
| `DELETE` | `/api/meta-ads/accounts/:id` | удалить |
| `POST` | `/api/meta-ads/accounts/:id/verify` | проверить токен и доступ |
| `POST` | `/api/meta-ads/accounts/:id/sync` | `{since, until}` — синхронизация |
| `GET` | `/api/meta-ads/spend` | `?ad_account_id=&since=&until=` — дневные строки |
| `GET` | `/api/meta-ads/summary` | агрегат по кампаниям |
| `PUT` | `/api/meta-ads/spend/:id/tour` | `{tour_id}` — привязать кампанию к туру |

Период по умолчанию, если `since`/`until` не переданы — последние 30 дней.

---

## 5. Ограничения и что дальше

- **Чужие рекламные кабинеты.** Чтобы читать кабинет клиента, он должен добавить
  наш system user в свой Business Manager (просто, без App Review). Полный OAuth
  для произвольных клиентов потребует Advanced Access `ads_read` от Meta.
- **Валюта.** Расход пишется в валюте кабинета как есть, без конвертации.
  Если кабинет в USD, а учёт в TRY — нужна конвертация по курсу из `settings`.
- **Синхронизация ручная.** Автоматический ежедневный запуск (cron / горутина с
  тикером) пока не сделан — при необходимости вызывать
  `POST /api/meta-ads/accounts/:id/sync` по расписанию.
- **Атрибуция.** Meta пересчитывает `spend` последних дней — период стоит
  перетягивать с запасом (например, последние 7 дней ежедневно); upsert обновит
  цифры, а созданный ранее расход **не пересоздастся и не изменится**.
