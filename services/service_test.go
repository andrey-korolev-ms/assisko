package services

import (
	"assisko/models"
	"context"
	"errors"
	"testing"
)

type mockStore struct {
	saveErr    error
	saveID     int
	savedHuman models.Human
	saveCalled bool
}

func (m *mockStore) SimpleSaveAtDB(ctx context.Context, p models.Human) (int, error) {
	m.savedHuman = p
	m.saveCalled = true
	return m.saveID, m.saveErr
}
func (m *mockStore) GetByID(ctx context.Context, id int) (models.Human, error) {
	return models.Human{}, nil
}
func (m *mockStore) GetAll(ctx context.Context) ([]models.Human, error) {
	return []models.Human{}, nil
}
func (m *mockStore) CreateUser(ctx context.Context, username, password string) (int, error) {
	return 0, nil
}
func (m *mockStore) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return models.User{}, nil
}

func TestSaveData(t *testing.T) {
	tests := []struct {
		name       string
		input      models.Human
		mockErr    error
		mockID     int
		wantErr    bool
		wantStatus bool
		wantID     int
	}{{
		name:       "saved success",
		input:      models.Human{Sex: "male", Age: 54},
		mockErr:    nil,
		mockID:     1,
		wantErr:    false,
		wantStatus: true,
		wantID:     1,
	}, {
		name:       "ошибка базы пробрасывается",
		input:      models.Human{Sex: "female", Age: 34},
		mockErr:    errors.New("База недоступна"),
		mockID:     0,
		wantErr:    true,
		wantStatus: false,
		wantID:     0,
	}, {
		name:       "пустые данные тоже сохраняются",
		input:      models.Human{},
		mockErr:    nil,
		mockID:     0,
		wantErr:    false,
		wantStatus: true,
		wantID:     0,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &mockStore{
				saveErr: tt.mockErr,
				saveID:  tt.mockID,
			}
			svc := NewService(ms)
			result, err := svc.SaveData(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveData() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if result.ID != tt.wantID {
					t.Errorf("%v,%v ", result.ID, tt.wantID)
				}
			}
			if !ms.saveCalled && tt.wantErr {
				t.Errorf("store.SimpleSaveData не был вызван")
			}
			if !tt.wantErr && ms.savedHuman != tt.input {
				t.Errorf("в стор ушло %+v, ожидалось %+v", ms.savedHuman, tt.input)
			}
		})
	}
}

func isInt(v any) bool {
	// Проверяет, лежит ли в интерфейсе тип int
	_, ok := v.(int)
	return ok
}
