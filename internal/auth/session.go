package auth

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const sessionTTL = 7 * 24 * time.Hour

type session struct {
	userID  int64
	expires time.Time
}

// Store — bellek içi oturum deposu (tekil sunucu instance'ı için yeterli).
type Store struct {
	mu       sync.Mutex
	sessions map[string]session
}

func NewStore() *Store {
	return &Store{sessions: make(map[string]session)}
}

func (s *Store) Create(userID int64) string {
	token := generateToken()
	s.mu.Lock()
	s.sessions[token] = session{userID: userID, expires: time.Now().Add(sessionTTL)}
	s.mu.Unlock()
	return token
}

func (s *Store) UserID(token string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[token]
	if !ok || time.Now().After(sess.expires) {
		delete(s.sessions, token)
		return 0, false
	}
	return sess.userID, true
}

func (s *Store) Delete(token string) {
	s.mu.Lock()
	delete(s.sessions, token)
	s.mu.Unlock()
}

func generateToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}
