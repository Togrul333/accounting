package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockAPI — подставной Bot API: отдаёт заранее заданный ответ и запоминает запрос.
func mockAPI(t *testing.T, response string, status int) (*httptest.Server, *string, *string) {
	t.Helper()
	var gotPath, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		gotBody = string(body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		io.WriteString(w, response)
	}))
	t.Cleanup(srv.Close)

	prev := apiBase
	apiBase = srv.URL
	t.Cleanup(func() { apiBase = prev })

	return srv, &gotPath, &gotBody
}

func TestGetMe(t *testing.T) {
	_, path, _ := mockAPI(t, `{"ok":true,"result":{"id":1,"username":"hisar_bot","first_name":"Hisar"}}`, http.StatusOK)

	bot, err := New("token123").GetMe(context.Background())
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if bot.Username != "hisar_bot" {
		t.Errorf("username = %q, beklenen %q", bot.Username, "hisar_bot")
	}
	if *path != "/bottoken123/getMe" {
		t.Errorf("path = %q", *path)
	}
}

func TestSendMessageBody(t *testing.T) {
	_, path, body := mockAPI(t, `{"ok":true,"result":{}}`, http.StatusOK)

	if err := New("t").SendMessage(context.Background(), "555", "<b>merhaba</b>"); err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if *path != "/bott/sendMessage" {
		t.Errorf("path = %q", *path)
	}

	var sent map[string]any
	if err := json.Unmarshal([]byte(*body), &sent); err != nil {
		t.Fatalf("gövde çözümlenemedi: %v", err)
	}
	if sent["chat_id"] != "555" || sent["text"] != "<b>merhaba</b>" || sent["parse_mode"] != "HTML" {
		t.Errorf("gönderilen gövde beklenenden farklı: %v", sent)
	}
}

// Бот не может писать первым — 403 должен превращаться в ErrChatNotStarted,
// чтобы UI показал понятную подсказку вместо «внутренней ошибки».
func TestSendMessageForbidden(t *testing.T) {
	mockAPI(t, `{"ok":false,"error_code":403,"description":"Forbidden: bot can't initiate conversation with a user"}`, http.StatusForbidden)

	err := New("t").SendMessage(context.Background(), "555", "merhaba")
	if !errors.Is(err, ErrChatNotStarted) {
		t.Fatalf("hata = %v, beklenen ErrChatNotStarted", err)
	}
}

func TestBadTokenIsAPIError(t *testing.T) {
	mockAPI(t, `{"ok":false,"error_code":401,"description":"Unauthorized"}`, http.StatusUnauthorized)

	_, err := New("bad").GetMe(context.Background())
	if err == nil {
		t.Fatal("hata bekleniyordu")
	}
	if !IsAPIError(err) {
		t.Errorf("IsAPIError = false, hata: %v", err)
	}
	if !strings.Contains(err.Error(), "Unauthorized") {
		t.Errorf("hata metni = %q", err.Error())
	}
}

func TestEmptyTokenRejected(t *testing.T) {
	if _, err := New("").GetMe(context.Background()); err == nil {
		t.Fatal("boş token için hata bekleniyordu")
	}
}

func TestGetUpdatesOffset(t *testing.T) {
	_, _, body := mockAPI(t, `{"ok":true,"result":[{"update_id":7,"message":{"message_id":1,"text":"/start abc","chat":{"id":42,"type":"private"},"from":{"id":42,"username":"ali"}}}]}`, http.StatusOK)

	updates, err := New("t").GetUpdates(context.Background(), 5)
	if err != nil {
		t.Fatalf("beklenmeyen hata: %v", err)
	}
	if len(updates) != 1 || updates[0].UpdateID != 7 {
		t.Fatalf("updates = %+v", updates)
	}
	if updates[0].Message.Chat.ID != 42 || updates[0].Message.From.Username != "ali" {
		t.Errorf("mesaj alanları beklenenden farklı: %+v", updates[0].Message)
	}
	if !strings.Contains(*body, `"offset":5`) {
		t.Errorf("offset gönderilmedi: %s", *body)
	}
}
