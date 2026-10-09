package services

import (
	"assisko/models"
	"assisko/store"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store store.Storage
}

func NewService(st store.Storage) *Service {
	return &Service{store: st}
}

func (s *Service) RegisterUser(ctx context.Context, username, password string) error {
	return s.store.CreateUser(ctx, username, password)
}

func (s *Service) AuthenticateUser(ctx context.Context, username, password string) (models.User, error) {
	user, err := s.store.GetUserByUsername(ctx, username)
	if err != nil {
		return models.User{}, err
	}
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return models.User{}, errors.New("неверный пароль")
	}

	return user, err
}

func (s *Service) SaveData(ctx context.Context, data models.Human) error {
	//допишу валидацию и логику позднее
	return s.store.SimpleSaveAtDB(ctx, data)
}

func (s *Service) GetData(ctx context.Context, id int) (models.Human, error) {
	return s.store.GetByID(ctx, id)
}
