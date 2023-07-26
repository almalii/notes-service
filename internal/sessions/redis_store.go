package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type SessionStore struct {
	client *redis.Client
}

func (s *SessionStore) Save(ctx context.Context, session *Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal sessions data: %w", err)
	}

	err = s.client.Set(ctx, session.ID, data, session.ExpiresAt).Err()
	if err != nil {
		return fmt.Errorf("failed to set sessions in Redis: %w", err)
	}

	return nil
}

func (s *SessionStore) Get(ctx context.Context, sessionID string) (*Session, error) {
	if s.client == nil {
		return nil, fmt.Errorf("Redis client is nil")
	}

	data, err := s.client.Get(ctx, sessionID).Bytes()
	if err != nil {
		if err == redis.Nil {
			logrus.Println("session not found")
			return nil, err
		}
		return nil, fmt.Errorf("failed to get session from Redis: %w", err)
	}

	var session *Session
	err = json.Unmarshal(data, &session) // Обратите внимание на использование &session
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal session data: %w", err)
	}

	return session, nil
}

func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	err := s.client.Del(ctx, sessionID).Err()
	if err != nil {
		return fmt.Errorf("failed to delete sessions from Redis: %w", err)
	}
	return nil
}

func NewRedisSessionStore(addr, password string, db int) *SessionStore {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &SessionStore{
		client: client,
	}
}
