# Excel — импорт банковской выписки

## Статус: реализовано

Кнопка **"Excel'den İçe Aktar"** на странице `/accounts/:id/edit`. Импортирует выписку по счёту
(`.xlsx`/`.xls`) в `incomes`/`expenses`. Библиотека: `github.com/xuri/excelize/v2`.

## Как это работает

1. Клик по кнопке → скрытый `<input type="file">` (`web/templates/account_edit.html`), выбор файла →
   `handleExcelImport` шлёт его на `POST /api/accounts/:id/statement-preview` (multipart, поле `file`).
2. `AccountHandler.ParseStatement` (`internal/handler/account.go`) открывает файл через `excelize`,
   парсит (`parseStatement`) и возвращает превью — **ничего не сохраняется в БД на этом шаге**.
3. Фронт открывает модалку "Banka Ekstresi — Önizleme" с двумя таблицами (Gelirler/Giderler) и двумя
   режимами разбора: построчный визард ("Tək-tək") и массовый ("Hamısını bir dəfəyə"), плюс
   "Tümüne Uygula" — применить категорию/тур/клиента сразу ко всем строкам.
4. Строки, уже присутствующие в БД (по `bank_ref`), помечаются бейджем "Artıq mövcuddur"
   (сервер присылает это в поле `already_imported`).
5. Пользователь расставляет категории и жмёт сохранить → фронт шлёт отдельно:
   - `POST /api/accounts/:id/statement-import-incomes` — тело `{rows: [...]}`, у каждой строки
     свои `income_category_id`, опционально `tour_id`+`client_id`.
   - `POST /api/accounts/:id/statement-import-expenses` — тело `{expense_category_id, rows: [...]}`,
     категория одна на весь запрос (не по строкам).
6. После успеха — тост, затем `location.reload()`.

## Парсинг файла (`parseStatement`, `internal/handler/account.go`)

Вся логика парсинга живёт прямо в handler'е — отдельного "ImportService" нет,
`IncomeService`/`ExpenseService` тут только тонкие прокси к репозиторию (`BulkCreate`,
`GetBankRefsByAccountID`).

1. **IBAN**: сканируются все ячейки первого листа regex'ом `^AZ[A-Z0-9]{20,}$`. Не найден → ошибка
   `"IBAN not found in file"`. Найден, но не совпадает с `account.account_number` счёта → ошибка
   `"IBAN mismatch: file has X, account has Y"` (защита от загрузки выписки чужого счёта).
2. **Строка заголовков**: ищется строка, где есть ячейка с подстрокой `"Əməliyyatın tarixi"`
   (азербайджанский, "Дата операции"). Не найдена / нет колонок Debit+Credit → ошибка
   `"unrecognised Excel format"`.
3. **Маппинг колонок** — по вхождению подстроки в заголовок (регистрозависимо, как в оригинале банка):

   | Подстрока в заголовке | Поле |
   |---|---|
   | `tarixi` | Date |
   | `referens` | Ref (для дедупликации по `bank_ref`) |
   | `hesab` | CP (контрагент) |
   | `Debit` | Debit → строка = расход |
   | `Kredit` | Credit → строка = доход |
   | `yinat` (часть "izahat"/"mahiyyət") | Desc |
   | `VÖEN` | Tax (налоговый ID контрагента) |

4. Числа (`Debit`/`Credit`) парсятся с удалением разделителей тысяч (запятых). Полностью пустые строки
   (`date=="" && debit==0 && credit==0`) пропускаются.
5. `credit > 0` → строка в `Gelirler`, `debit > 0` → в `Giderler`; считаются `TotalCredit`/`TotalDebit`.

## Сохранение — `ImportIncomes` / `ImportExpenses`

- **Дедупликация по `bank_ref`**: сверка и с уже сохранёнными в БД (`GetBankRefsByAccountID`), и с
  дублями внутри самого запроса (map `seen`). Совпавшие строки тихо пропускаются
  (`skipped_duplicates` в ответе), без ошибки всего запроса.
- **Дата** парсится строго форматом `02.01.2006`; если не парсится — строка тоже пропускается как
  skipped (без отдельного сообщения пользователю, только счётчик).
- **Доходы**: `income_category_id` обязателен на каждую строку (400, если 0). Если указаны и
  `tour_id`, и `client_id` — доход привязывается к заказу через
  `OrderService.FindOrCreateByClientAndTour` (кэш `client:tour` внутри запроса, чтобы не плодить
  дубли заказов при массовом импорте). `name` = `Desc`, если пусто — `Ref`.
- **Расходы**: `expense_category_id` обязателен один на весь запрос (400, если 0).
- Вставка — `IncomeService/ExpenseService.BulkCreate` → репозиторий, в транзакции.

## Модель данных

`internal/model/account.go`:
```go
type StatementRow struct {
    Date, Ref, CP string
    Debit, Credit float64
    Desc, Tax     string
    AlreadyImported bool
}
type ImportGelirRow struct { StatementRow; IncomeCategoryID int64; TourID *int64; ClientID *int64 }
type ImportGiderRow struct { StatementRow }
type StatementPreview struct { IBAN string; Gelirler, Giderler []StatementRow; TotalCredit, TotalDebit float64 }
```

Поля-мост между импортом и обычными доменными моделями — `bank_ref`, `counterparty`,
`counterparty_tax_id` в `model.Income`/`model.Expense` (миграция `019_add_bank_fields.sql`,
добавляет эти три колонки в `incomes` и `expenses`). Отдельной таблицы истории импортов нет —
дедупликация целиком опирается на `bank_ref` в самих таблицах `incomes`/`expenses`.

## Валидация и ошибки — сводка

- Файл обязателен (400 `"file is required"`), должен открываться как Excel (400 `"invalid Excel file"`).
- IBAN должен быть найден и совпадать со счётом (422).
- Должны быть распознаваемые заголовки + колонки Debit/Credit (422 `"unrecognised Excel format"`).
- Категория обязательна: по строке для доходов, одна на запрос для расходов (400).
- Дубли по `bank_ref` и строки с непарсящейся датой — пропускаются молча (счётчик `skipped_duplicates`).
- Ошибки БД/сервисов — 500 с текстом ошибки.

## Известные ограничения / куда смотреть при доработке

- Формат заголовков жёстко завязан на азербайджанские подстроки конкретного банка — под другой банк
  потребуется расширять/переписывать маппинг в `parseStatement`.
- Для расходов нет привязки к туру при импорте (`tour_id` не заполняется handler'ом), в отличие от
  доходов.
- Импорт не документирован нигде, кроме этого файла — этим он ранее отличался от
  [Google Sheets импорта](google-sheets-import.md).
