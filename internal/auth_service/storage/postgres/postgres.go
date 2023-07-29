package postgres

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/service"
	"notes-rew/internal/auth_service/storage"
)

type UserStorage struct {
	db *pgx.Conn
}

func (s *UserStorage) SaveUserToDB(ctx context.Context, user service.CreateUser) error {
	sql, args, err := squirrel.Insert("users").
		Columns("id", "username", "email", "password", "created_at", "updated_at").
		Values(user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		logrus.Errorf("error while executing squirrel query: %v", err)
	}

	return nil
}

func (s *UserStorage) GetUserForAuth(ctx context.Context, email string) (models.AuthOutput, error) {
	var user storage.AuthResponse

	sql, args, err := squirrel.Select("id", "username", "email", "password").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
		return models.AuthOutput{}, err
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			logrus.Println("users_service not found")
			return models.AuthOutput{}, err
		}
		logrus.Errorf("error while getting users_service by auth_service: %v", err)
		return models.AuthOutput{}, err
	}

	resp := storage.NewAuthResponse(user.ID, user.Username, user.Email, user.PasswordHash)

	return resp, nil
}

func NewUserStorage(db *pgx.Conn) *UserStorage {
	return &UserStorage{db: db}
}
