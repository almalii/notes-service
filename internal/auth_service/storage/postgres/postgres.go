package postgres

import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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
		return err
	}

	_, err = s.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
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
		return models.AuthOutput{}, err
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
		return models.AuthOutput{}, err
	}

	resp := storage.NewAuthResponse(user.ID, user.Username, user.Email, user.PasswordHash)

	return resp, nil
}

func (s *UserStorage) CheckUserByEmail(ctx context.Context, email string) error {
	var count int

	sql, args, err := squirrel.Select("count(*)").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return err
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("user already exists")
	}

	return nil
}

func NewUserStorage(db *pgx.Conn) *UserStorage {
	return &UserStorage{db: db}
}
