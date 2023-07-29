package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"notes-rew/internal/auth_service/models"
	"notes-rew/internal/auth_service/service"
)

type UserStorage struct {
	db *mongo.Client
}

func (s *UserStorage) SaveUserToDB(ctx context.Context, user service.CreateUser) error {

	return nil
}
func (s *UserStorage) GetUserForAuth(ctx context.Context, email string) (models.AuthOutput, error) {
	return models.AuthOutput{}, nil
}

func NewUserStorage(db *mongo.Client) *UserStorage {
	return &UserStorage{db: db}
}
