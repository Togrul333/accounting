# Google Sheets — импорт данных

## Что будет сделано

Пользователь вставляет ссылку на Google Sheet → приложение читает строки через Google Sheets API → данные сохраняются в MySQL.

## Настройка (один раз)

1. Зайти на **https://console.cloud.google.com**
2. Создать проект → включить **Google Sheets API**
3. Создать **Service Account** → скачать JSON-ключ
4. Поделиться нужным Google Sheet с email service account (`xxx@project.iam.gserviceaccount.com`)
5. Положить JSON-файл в корень проекта, например `credentials.json`

## Переменная окружения

```bash
GOOGLE_CREDENTIALS_PATH=./credentials.json
```

## Как будет работать в UI

1. Пользователь вставляет ссылку вида:
   ```
   https://docs.google.com/spreadsheets/d/<SPREADSHEET_ID>/edit...
   ```
2. Бэкенд извлекает `SPREADSHEET_ID` из ссылки
3. Читает строки через Sheets API
4. Парсит данные и сохраняет в MySQL

## Зависимости Go

```bash
go get google.golang.org/api/sheets/v4
go get golang.org/x/oauth2/google
```
