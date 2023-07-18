package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type RedisSessionStore struct {
	client *redis.Client
}

func NewRedisSessionStore(addr, password string, db int) *RedisSessionStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisSessionStore{
		client: client,
	}
}

type Session struct {
	SessionID string
	Values    map[string]interface{}
	CreatedAt time.Time
	ExpiresAt time.Time
}

// CreateSession создает новый сеанс.
func CreateSession(duration time.Duration) (*Session, error) {
	return &Session{
		SessionID: GenerateUniqueID(),
		CreatedAt: time.Now().UTC(),
		ExpiresAt: time.Now().UTC().Add(duration),
		Values:    make(map[string]interface{}),
	}, nil
}

// GenerateUniqueID генерирует уникальный идентификатор сеанса.
func GenerateUniqueID() string {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		logrus.Error("error while generating random bytes")
		return err.Error()
	}

	uuIdentification := hex.EncodeToString(bytes)

	return uuIdentification
}

// Set сохраняет сеанс в Redis.
func (s *RedisSessionStore) Set(ctx context.Context, session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session data: %w", err)
	}

	key := fmt.Sprintf("session:%s", session.SessionID)
	duration := session.ExpiresAt.Sub(time.Now())

	err = s.client.Set(ctx, key, data, duration).Err()
	if err != nil {
		return fmt.Errorf("failed to set session in Redis: %w", err)
	}

	return nil
}

// Get возвращает сеанс с указанным идентификатором из Redis.
func (s *RedisSessionStore) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := fmt.Sprintf("session:%s", sessionID)

	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			logrus.Println("session not found")
			return nil, nil // Сеанс не найден
		}
		return nil, fmt.Errorf("failed to get session from Redis: %w", err)
	}

	session := &Session{}
	err = json.Unmarshal(data, session)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return session, nil
}

// Delete удаляет сеанс с указанным идентификатором из Redis.
func (s *RedisSessionStore) Delete(ctx context.Context, sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session from Redis: %w", err)
	}
	return nil
}
