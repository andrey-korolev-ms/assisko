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

	//поля для чтения
	getHuman models.Human // заготовка: что вернет мок
	getErr   error        // заготовка: какую ошибку вернет мок
	getIDArg int          // запомнили, какой ID запросили у мока
}

func (m *mockStore) SimpleSaveAtDB(ctx context.Context, p models.Human) (int, error) {
	m.savedHuman = p
	m.saveCalled = true
	return m.saveID, m.saveErr
}
func (m *mockStore) GetByID(ctx context.Context, id int) (models.Human, error) {
	m.getIDArg = id
	return m.getHuman, m.getErr
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
		name           string
		input          models.Human
		wantStoreHuman models.Human
		mockErr        error
		mockID         int
		wantErr        bool
		wantResult     models.Right
	}{{
		name:           "saved success",
		input:          models.Human{Sex: "male", Age: 54},
		wantStoreHuman: models.Human{Sex: "Мой пол_male", Age: 55},
		mockErr:        nil,
		mockID:         1,
		wantErr:        false,
		wantResult: models.Right{
			Status: true,
			ID:     1,
			Data:   models.Human{Sex: "Мой пол_male", Age: 55},
		},
	}, {
		name:           "ошибка базы пробрасывается",
		input:          models.Human{Sex: "female", Age: 34},
		wantStoreHuman: models.Human{Sex: "Мой пол_female", Age: 35},
		mockErr:        errors.New("База недоступна"),
		mockID:         0,
		wantErr:        true,
		wantResult: models.Right{
			Status: true,
			ID:     1,
			Data:   models.Human{Sex: "female", Age: 35},
		},
	}, {
		name:           "пустые данные тоже сохраняются",
		input:          models.Human{},
		wantStoreHuman: models.Human{Sex: "Мой пол_", Age: 1},
		mockErr:        nil,
		mockID:         1,
		wantErr:        false,
		wantResult: models.Right{
			Status: true,
			ID:     1,
			Data:   models.Human{Sex: "Мой пол_", Age: 1},
		},
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
				if result != tt.wantResult {
					t.Errorf("%v,%v ", result, tt.wantResult)
				}
			}
			if !ms.saveCalled && tt.wantErr {
				t.Errorf("store.SimpleSaveData не был вызван")
			}

			if !tt.wantErr && ms.savedHuman != tt.wantStoreHuman {
				t.Errorf("в стор ушло %+v, ожидалось %+v", ms.savedHuman, tt.wantStoreHuman)
			}
		})
	}
}
func TestGetData(t *testing.T) {
	tests := []struct {
		name      string
		inputID   int
		mockHuman models.Human
		mockErr   error
		wantErr   bool
		wantHuman models.Human
	}{{
		name:    "успешное получение данных",
		inputID: 1,
		mockHuman: models.Human{
			Sex: "male",
			Age: 30,
		},
		mockErr: nil,
		wantErr: false,
		wantHuman: models.Human{
			Sex: "male",
			Age: 30,
		},
	}, {
		name:      "Пользователь не найден",
		inputID:   999,
		mockHuman: models.Human{},
		mockErr:   errors.New("пользователь не найден"),
		wantErr:   true,
		wantHuman: models.Human{},
	}, {
		name:      "ошибка базы данных",
		inputID:   5,
		mockHuman: models.Human{},
		mockErr:   errors.New("DB connect failed"),
		wantErr:   true,
		wantHuman: models.Human{},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := &mockStore{
				getHuman: tt.mockHuman,
				getErr:   tt.mockErr,
			}
			svc := NewService(ms)

			result, err := svc.GetData(context.Background(), tt.inputID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if result != tt.wantHuman {
					t.Errorf("GetData() error = %+v, want %+v ", result, tt.wantHuman)
				}
			}
			if ms.getIDArg != tt.inputID {
				t.Errorf("в стор ушел id = %d, ожидалось %d", ms.getIDArg, tt.inputID)
			}
		})
	}
}
