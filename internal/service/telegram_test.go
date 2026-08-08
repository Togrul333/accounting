package service

import (
	"strings"
	"testing"
	"time"

	"accounting/internal/model"
)

func TestStartCode(t *testing.T) {
	cases := []struct {
		text string
		want string
		ok   bool
	}{
		{"/start abc123", "abc123", true},
		{"  /start abc123  ", "abc123", true},
		{"/START abc123", "abc123", true},
		{"/start@hisar_bot abc123", "abc123", true},
		{"/start", "", false},
		{"/start a b", "", false},
		{"merhaba", "", false},
		{"", "", false},
	}
	for _, c := range cases {
		got, ok := startCode(c.text)
		if ok != c.ok || got != c.want {
			t.Errorf("startCode(%q) = (%q, %v), beklenen (%q, %v)", c.text, got, ok, c.want, c.ok)
		}
	}
}

// Заголовок и описание задачи попадают в HTML-сообщение, поэтому их нужно экранировать.
func TestTaskMessageEscapesHTML(t *testing.T) {
	due := time.Date(2026, 8, 20, 0, 0, 0, 0, time.UTC)
	msg := taskMessage(model.Task{
		ID:          14,
		Title:       "<script>alert(1)</script>",
		Description: "R&D & test",
		Priority:    "urgent",
		DueDate:     &due,
		TourCode:    "T-1",
		OrderLabel:  "#9 Ali Veli",
		ClientName:  "Ayşe Yılmaz",
	})

	if strings.Contains(msg, "<script>") {
		t.Errorf("başlık kaçırılmamış: %s", msg)
	}
	if !strings.Contains(msg, "&lt;script&gt;") {
		t.Errorf("beklenen kaçış yok: %s", msg)
	}
	if !strings.Contains(msg, "R&amp;D") {
		t.Errorf("açıklama kaçırılmamış: %s", msg)
	}
	for _, want := range []string{"GRV-14", "Acil", "20.08.2026", "T-1", "#9 Ali Veli", "Ayşe Yılmaz"} {
		if !strings.Contains(msg, want) {
			t.Errorf("mesajda %q yok:\n%s", want, msg)
		}
	}
}

func TestTaskMessageMinimal(t *testing.T) {
	msg := taskMessage(model.Task{ID: 1, Title: "Basit görev"})
	if !strings.Contains(msg, "GRV-1") || !strings.Contains(msg, "Basit görev") {
		t.Errorf("beklenmeyen mesaj: %s", msg)
	}
	if strings.Contains(msg, "Son tarih") || strings.Contains(msg, "Tur:") {
		t.Errorf("boş alanlar mesaja girmiş: %s", msg)
	}
}

func TestAssigneeChanged(t *testing.T) {
	id2, id3 := int64(2), int64(3)
	cases := []struct {
		name          string
		before, after *int64
		want          bool
	}{
		{"atama yok", nil, nil, false},
		{"yeni atama", nil, &id2, true},
		{"başkasına devredildi", &id2, &id3, true},
		{"aynı kişi", &id2, &id2, false},
		{"atama kaldırıldı", &id2, nil, false},
	}
	for _, c := range cases {
		got := assigneeChanged(&model.Task{HocaUserID: c.before}, &model.Task{HocaUserID: c.after})
		if got != c.want {
			t.Errorf("%s: assigneeChanged = %v, beklenen %v", c.name, got, c.want)
		}
	}
}
