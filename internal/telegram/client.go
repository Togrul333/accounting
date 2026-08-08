// Package telegram — тонкий клиент Telegram Bot API.
// Знает только про HTTP: настройки и бизнес-логика живут в service.
package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// apiBase — переменная, а не константа, чтобы тесты могли подменить её моком.
var apiBase = "https://api.telegram.org"

// ErrChatNotStarted — бот не может писать первым: пользователь не нажал Start
// или заблокировал бота.
var ErrChatNotStarted = errors.New("kullanıcı botu başlatmamış")

type Client struct {
	token string
	http  *http.Client
}

func New(token string) *Client {
	return &Client{token: token, http: &http.Client{Timeout: 10 * time.Second}}
}

// Bot — ответ getMe.
type Bot struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"first_name"`
}

// Chat — чат, из которого пришло сообщение.
type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Username string `json:"username"`
	First    string `json:"first_name"`
	Last     string `json:"last_name"`
	Title    string `json:"title"`
}

type Message struct {
	MessageID int64 `json:"message_id"`
	Chat      Chat  `json:"chat"`
	From      struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		First    string `json:"first_name"`
		Last     string `json:"last_name"`
	} `json:"from"`
	Text string `json:"text"`
}

type Update struct {
	UpdateID int64    `json:"update_id"`
	Message  *Message `json:"message"`
}

// apiError — ошибка уровня Telegram (ok:false), а не транспорта.
type apiError struct {
	Code        int
	Description string
}

func (e *apiError) Error() string {
	return fmt.Sprintf("telegram: %s (%d)", e.Description, e.Code)
}

// IsAPIError — ошибка пришла от самого Telegram (ok:false), а не от сети.
// Такие ошибки чинит пользователь (неверный токен, чужой chat_id), поэтому
// наружу они отдаются как 400.
func IsAPIError(err error) bool {
	var e *apiError
	return errors.As(err, &e)
}

// call — общий вызов метода Bot API с JSON-телом.
func (c *Client) call(ctx context.Context, method string, payload any, out any) error {
	if c.token == "" {
		return errors.New("telegram: bot token boş")
	}

	var body []byte
	if payload != nil {
		var err error
		if body, err = json.Marshal(payload); err != nil {
			return err
		}
	}

	url := fmt.Sprintf("%s/bot%s/%s", apiBase, c.token, method)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var envelope struct {
		OK          bool            `json:"ok"`
		Result      json.RawMessage `json:"result"`
		ErrorCode   int             `json:"error_code"`
		Description string          `json:"description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return fmt.Errorf("telegram: cevap okunamadı: %w", err)
	}
	if !envelope.OK {
		if envelope.ErrorCode == http.StatusForbidden {
			return ErrChatNotStarted
		}
		return &apiError{Code: envelope.ErrorCode, Description: envelope.Description}
	}
	if out != nil {
		return json.Unmarshal(envelope.Result, out)
	}
	return nil
}

// GetMe — проверка токена; возвращает данные бота.
func (c *Client) GetMe(ctx context.Context) (*Bot, error) {
	var bot Bot
	if err := c.call(ctx, "getMe", nil, &bot); err != nil {
		return nil, err
	}
	return &bot, nil
}

// SendMessage — отправка сообщения в чат. text в разметке HTML.
func (c *Client) SendMessage(ctx context.Context, chatID, text string) error {
	return c.call(ctx, "sendMessage", map[string]any{
		"chat_id":                  chatID,
		"text":                     text,
		"parse_mode":               "HTML",
		"disable_web_page_preview": true,
	}, nil)
}

// GetUpdates — накопившиеся обновления начиная с offset.
// Используется вместо webhook: приложение локальное и внешнего URL у него нет.
func (c *Client) GetUpdates(ctx context.Context, offset int64) ([]Update, error) {
	var updates []Update
	payload := map[string]any{
		"limit":           100,
		"timeout":         0,
		"allowed_updates": []string{"message"},
	}
	if offset > 0 {
		payload["offset"] = offset
	}
	if err := c.call(ctx, "getUpdates", payload, &updates); err != nil {
		return nil, err
	}
	return updates, nil
}
