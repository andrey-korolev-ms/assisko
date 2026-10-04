package services

import (
	"assisko/models"
	"assisko/store"
	"context"
)

type Service struct {
	store store.Storage
}

func NewService(st store.Storage) *Service {
	return &Service{store: st}
}

func (s *Service) SaveData(ctx context.Context, data models.Human) error {
	//допишу валидацию и логику позднее
	return s.store.SimpleSaveAtDB(ctx, data)
}

func (s *Service) GetData(ctx context.Context, id int) (models.Human, error) {
	return s.store.GetByID(ctx, id)
}
