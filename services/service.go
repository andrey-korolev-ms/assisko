package services

import (
	"assisko/models"
	"assisko/store"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	store store.Storage
	queue *TaskQueue
}

func NewService(st store.Storage) *Service {
	return &Service{store: st}
}

func (s *Service) RegisterUser(ctx context.Context, username, password string) (int, error) {
	userID, err := s.store.CreateUser(ctx, username, password)
	if err != nil {
		return 0, fmt.Errorf("не удалось создать пользователя %w", err)
	}

	//кладем задачу в очередь

	payload, _ := json.Marshal(map[string]interface{}{
		"user_id": userID,
		"email":   username,
	})
	task := Task{
		ID:        fmt.Sprintf("task_%d", time.Now().UnixNano()),
		Type:      "send_welcome_email",
		Payload:   payload,
		CreatedAt: time.Now(),
	}
	if err := s.queue.Enqueue(task); err != nil {
		log.Printf("не удалось поставить задачу в очередь %v", err)
	}
	return userID, nil
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

func (s *Service) SetQueue(q *TaskQueue) {
	s.queue = q
}
