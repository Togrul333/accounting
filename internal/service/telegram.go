package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"html"
	"log"
	"strconv"
	"strings"
	"time"

	"accounting/internal/model"
	"accounting/internal/repository"
	"accounting/internal/telegram"
)

var (
	ErrTelegramDisabled = errors.New("telegram entegrasyonu kapalı")
	ErrTelegramNoToken  = errors.New("telegram bot token tanımlı değil")
	ErrTelegramNoChat   = errors.New("personel telegram hesabını bağlamamış")
)

// notifyTimeout — уведомления уходят вне HTTP-запроса, со своим сроком.
const notifyTimeout = 15 * time.Second

type TelegramService struct {
	settings repository.SettingRepository
	hocaRepo repository.HocaUserRepository
}

func NewTelegramService(settings repository.SettingRepository, hocaRepo repository.HocaUserRepository) *TelegramService {
	return &TelegramService{settings: settings, hocaRepo: hocaRepo}
}

// config — сырые настройки бота из таблицы settings.
type telegramConfig struct {
	token    string
	username string
	enabled  bool
	offset   int64
}

func (s *TelegramService) config(ctx context.Context) (telegramConfig, error) {
	vals, err := s.settings.Get(ctx,
		model.SettingTelegramBotToken,
		model.SettingTelegramBotUsername,
		model.SettingTelegramEnabled,
		model.SettingTelegramOffset,
	)
	if err != nil {
		return telegramConfig{}, err
	}
	offset, _ := strconv.ParseInt(vals[model.SettingTelegramOffset], 10, 64)
	return telegramConfig{
		token:    vals[model.SettingTelegramBotToken],
		username: vals[model.SettingTelegramBotUsername],
		enabled:  vals[model.SettingTelegramEnabled] == "1",
		offset:   offset,
	}, nil
}

func (s *TelegramService) GetSettings(ctx context.Context) (model.TelegramSettings, error) {
	cfg, err := s.config(ctx)
	if err != nil {
		return model.TelegramSettings{}, err
	}
	return model.TelegramSettings{
		Enabled:     cfg.enabled,
		HasToken:    cfg.token != "",
		BotUsername: cfg.username,
	}, nil
}

// UpdateSettings — сохраняет токен, попутно проверяя его через getMe
// и запоминая username бота (он нужен для ссылки t.me/<bot>?start=<code>).
// Пустой токен в запросе означает «оставить текущий».
func (s *TelegramService) UpdateSettings(ctx context.Context, req model.UpdateTelegramSettingsRequest) (model.TelegramSettings, error) {
	cfg, err := s.config(ctx)
	if err != nil {
		return model.TelegramSettings{}, err
	}

	token := strings.TrimSpace(req.BotToken)
	if token == "" {
		token = cfg.token
	}
	if req.Enabled && token == "" {
		return model.TelegramSettings{}, ErrTelegramNoToken
	}

	values := map[string]string{
		model.SettingTelegramEnabled: boolToSetting(req.Enabled),
	}

	username := cfg.username
	if token != "" && token != cfg.token {
		bot, err := telegram.New(token).GetMe(ctx)
		if err != nil {
			return model.TelegramSettings{}, err
		}
		username = bot.Username
		values[model.SettingTelegramBotToken] = token
		values[model.SettingTelegramBotUsername] = username
	}

	if err := s.settings.Set(ctx, values); err != nil {
		return model.TelegramSettings{}, err
	}
	return model.TelegramSettings{Enabled: req.Enabled, HasToken: token != "", BotUsername: username}, nil
}

// client — готовый клиент, если интеграция включена и токен на месте.
func (s *TelegramService) client(ctx context.Context) (*telegram.Client, telegramConfig, error) {
	cfg, err := s.config(ctx)
	if err != nil {
		return nil, cfg, err
	}
	if cfg.token == "" {
		return nil, cfg, ErrTelegramNoToken
	}
	if !cfg.enabled {
		return nil, cfg, ErrTelegramDisabled
	}
	return telegram.New(cfg.token), cfg, nil
}

// LinkCode — выдаёт сотруднику разовый код и ссылку на бота.
// Код живёт до тех пор, пока сотрудник не нажмёт Start.
func (s *TelegramService) LinkCode(ctx context.Context, hocaID int64) (string, string, error) {
	cfg, err := s.config(ctx)
	if err != nil {
		return "", "", err
	}
	if cfg.token == "" {
		return "", "", ErrTelegramNoToken
	}

	hoca, err := s.hocaRepo.GetByID(ctx, hocaID)
	if err != nil {
		return "", "", err
	}

	code := hoca.TelegramLinkCode
	if code == "" {
		if code, err = randomCode(); err != nil {
			return "", "", err
		}
		if err := s.hocaRepo.SetLinkCode(ctx, hocaID, code); err != nil {
			return "", "", err
		}
	}

	link := ""
	if cfg.username != "" {
		link = fmt.Sprintf("https://t.me/%s?start=%s", cfg.username, code)
	}
	return code, link, nil
}

func (s *TelegramService) Unlink(ctx context.Context, hocaID int64) error {
	return s.hocaRepo.ClearTelegram(ctx, hocaID)
}

// SyncLinks — вычитывает новые сообщения бота и привязывает чаты к сотрудникам
// по коду из «/start <code>». Webhook не используется: приложение локальное.
// Возвращает имена привязанных сотрудников.
func (s *TelegramService) SyncLinks(ctx context.Context) ([]string, error) {
	client, cfg, err := s.client(ctx)
	if err != nil {
		return nil, err
	}

	updates, err := client.GetUpdates(ctx, cfg.offset)
	if err != nil {
		return nil, err
	}

	var linked []string
	maxID := cfg.offset
	for _, u := range updates {
		if u.UpdateID >= maxID {
			maxID = u.UpdateID + 1
		}
		if u.Message == nil {
			continue
		}
		code, ok := startCode(u.Message.Text)
		if !ok {
			continue
		}
		hoca, err := s.hocaRepo.GetByLinkCode(ctx, code)
		if err != nil {
			continue // код чужой или уже использован
		}
		chatID := strconv.FormatInt(u.Message.Chat.ID, 10)
		if err := s.hocaRepo.SetTelegram(ctx, hoca.ID, chatID, u.Message.From.Username); err != nil {
			log.Printf("telegram bağlantısı kaydedilemedi (hoca %d): %v", hoca.ID, err)
			continue
		}
		linked = append(linked, hoca.FullName())

		greeting := fmt.Sprintf("✅ <b>Bağlantı tamam</b>\nMerhaba %s, bundan sonra görevleriniz buraya gelecek.",
			html.EscapeString(hoca.FullName()))
		if err := client.SendMessage(ctx, chatID, greeting); err != nil {
			log.Printf("telegram karşılama mesajı gönderilemedi (hoca %d): %v", hoca.ID, err)
		}
	}

	// Сдвигаем offset, чтобы те же сообщения не разбирались повторно.
	if maxID != cfg.offset {
		if err := s.settings.Set(ctx, map[string]string{
			model.SettingTelegramOffset: strconv.FormatInt(maxID, 10),
		}); err != nil {
			log.Printf("telegram offset kaydedilemedi: %v", err)
		}
	}
	return linked, nil
}

// SendTest — проверочное сообщение конкретному сотруднику.
func (s *TelegramService) SendTest(ctx context.Context, hocaID int64) error {
	client, _, err := s.client(ctx)
	if err != nil {
		return err
	}
	hoca, err := s.hocaRepo.GetByID(ctx, hocaID)
	if err != nil {
		return err
	}
	if hoca.TelegramChatID == "" {
		return ErrTelegramNoChat
	}
	return client.SendMessage(ctx, hoca.TelegramChatID,
		"🔔 <b>Hisar Tour</b>\nTelegram bildirimleri çalışıyor.")
}

// NotifyTaskAssigned — уведомление исполнителю о назначенной задаче.
// Вызывается из TaskService в отдельной горутине: письмо не должно
// задерживать ответ на HTTP-запрос и не должно ронять сохранение задачи.
func (s *TelegramService) NotifyTaskAssigned(task model.Task) {
	if task.HocaUserID == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), notifyTimeout)
	defer cancel()

	client, _, err := s.client(ctx)
	if err != nil {
		if !errors.Is(err, ErrTelegramDisabled) && !errors.Is(err, ErrTelegramNoToken) {
			log.Printf("telegram bildirimi atlandı (görev %d): %v", task.ID, err)
		}
		return
	}

	hoca, err := s.hocaRepo.GetByID(ctx, *task.HocaUserID)
	if err != nil {
		log.Printf("telegram bildirimi: personel bulunamadı (görev %d): %v", task.ID, err)
		return
	}
	if hoca.TelegramChatID == "" {
		return // сотрудник ещё не нажал Start — это не ошибка
	}

	if err := client.SendMessage(ctx, hoca.TelegramChatID, taskMessage(task)); err != nil {
		log.Printf("telegram bildirimi gönderilemedi (görev %d, hoca %d): %v", task.ID, hoca.ID, err)
	}
}

// taskMessage — текст уведомления в разметке HTML.
func taskMessage(t model.Task) string {
	var b strings.Builder
	b.WriteString("🔔 <b>Yeni görev atandı</b>\n")
	fmt.Fprintf(&b, "<b>GRV-%d</b> — %s\n", t.ID, html.EscapeString(t.Title))
	fmt.Fprintf(&b, "Öncelik: %s\n", priorityLabel(t.Priority))

	if t.DueDate != nil {
		fmt.Fprintf(&b, "Son tarih: %s\n", t.DueDate.Format("02.01.2006"))
	}
	if t.TourCode != "" {
		fmt.Fprintf(&b, "Tur: %s\n", html.EscapeString(t.TourCode))
	}
	if t.OrderLabel != "" {
		fmt.Fprintf(&b, "Sipariş: %s\n", html.EscapeString(t.OrderLabel))
	}
	if t.ClientName != "" {
		fmt.Fprintf(&b, "Müşteri: %s\n", html.EscapeString(t.ClientName))
	}
	if t.Description != "" {
		fmt.Fprintf(&b, "\n%s", html.EscapeString(t.Description))
	}
	return b.String()
}

func priorityLabel(p string) string {
	switch p {
	case "urgent":
		return "🔴 Acil"
	case "high":
		return "🟠 Yüksek"
	case "low":
		return "⚪️ Düşük"
	default:
		return "🟡 Orta"
	}
}

// startCode — код привязки из команды «/start <code>».
func startCode(text string) (string, bool) {
	fields := strings.Fields(strings.TrimSpace(text))
	if len(fields) != 2 {
		return "", false
	}
	cmd := strings.ToLower(fields[0])
	if cmd != "/start" && !strings.HasPrefix(cmd, "/start@") {
		return "", false
	}
	return fields[1], true
}

func randomCode() (string, error) {
	buf := make([]byte, 6)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func boolToSetting(v bool) string {
	if v {
		return "1"
	}
	return "0"
}
