package storage

import (
	"context"
	"github.com/sirupsen/logrus"
	"notes-rew/user/service"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"notes-rew/user/models"
)

type PSQLUserStorage struct {
	db *pgx.Conn
}

func (s *PSQLUserStorage) CreateUserByID(ctx context.Context, user service.CreateUser) error {
	sql, args, err := squirrel.Insert("users").
		Columns("id", "username", "email", "password", "salt_key", "created_at", "updated_at").
		Values(user.ID, user.Username, user.Email, user.Password, user.SaltKey, user.CreatedAt, user.UpdatedAt).
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

func (s *PSQLUserStorage) GetUserByID(ctx context.Context, id uuid.UUID) (models.UserOutput, error) {
	user := new(UserResponse)

	sql, args, err := squirrel.Select("id", "username", "email", "created_at", "updated_at").
		From("users").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models.UserOutput{}, err
	}

	return models.UserOutput(*user), nil
}

func (s *PSQLUserStorage) UpdateUserByID(ctx context.Context, id uuid.UUID, user service.UpdateUser) error {

	sql, args, err := squirrel.Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("salt_key", user.SaltKey).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"id": id}).
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

func (s *PSQLUserStorage) DeleteUserByID(ctx context.Context, id uuid.UUID) error {

	sql, args, err := squirrel.Delete("users").
		Where(squirrel.Eq{"id": id}).
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

func (s *PSQLUserStorage) GetUserForAuth(ctx context.Context, email string) (models.AuthOutput, error) {
	user := new(AuthResponse)

	sql, args, err := squirrel.Select("id", "username", "email", "password", "salt_key").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
		return models.AuthOutput{}, err
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.SaltKey)
	if err != nil {
		if err == pgx.ErrNoRows {
			logrus.Println("user not found")
			return models.AuthOutput{}, err
		}
		logrus.Errorf("error while getting user by auth: %v", err)
		return models.AuthOutput{}, err
	}

	return models.AuthOutput(*user), nil
}

func (s *PSQLUserStorage) CheckerByEmail(ctx context.Context, email string) (bool, error) {
	var count int

	sql, args, err := squirrel.Select("count(*)").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		logrus.Errorf("error while building squirrel query: %v", err)
		return false, err
	}

	err = s.db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		logrus.Errorf("error while getting count by email: %v", err)
		return false, err
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func NewPSQLUserStorage(db *pgx.Conn) *PSQLUserStorage {
	return &PSQLUserStorage{
		db: db,
	}
}
