package services

import (
	"assisko/models"
	"context"
	"testing"
)

// mockStore - это наша тестовая заглушка (мок);
// она реализует интефейс store.Storage но не лезет в реальную базу
type mockStore struct {
	saveErr   error        // ошибку, которую мы ходим симулировать
	saveHuman models.Human // сюда мы заполмним, что именно пришло в метод
}

func (m *mockStore) SimpleSaveAtDB(ctx context.Context, p models.Human) error {
	m.saveHuman = p  // запоминаем входные данные для проверки
	return m.saveErr // возвращаем заранее заданную ошибку (или nil)
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

// --- табличные тесты ---

func TestSaveData(t *testing.T) {
	tests := []struct {
		name        string
		input       models.Human
		mockSaveErr error
		wantErr     bool
		wantStatus  bool
	}{
		{

			name: "успешное сохранение",
			input: models.Human{
				Sex:              "male",
				Age:              18,
				IPK:              10.1,
				SimpleEmployment: 10,
				IsInvalid:        true,
				InvalidGroup:     2,
				OperatorID:       1,
			},
			mockSaveErr: nil,
			wantErr:     false,
			wantStatus:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Создаём мок с нужным поведением для этого сценария
			ms := &mockStore{
				saveErr: tt.mockSaveErr,
			}

			// 2. Создаем сервис, передаем ему мок вместо реальной базы
			svr := NewService(ms)

			// 3. Вызываю тестируемый метод
			err := svr.SaveData(context.Background(), tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("SaveData() error  = %v, wantErr %v", err, tt.wantErr)
			}
			// если ошибки не ожидалось, проверяем результат
			if !tt.wantErr {
				// if result.Status != tt.wantStatus {
				// 	t.Errorf("SaveData() status %v, want %v", result.Status, tt.wantStatus)
				// }
				// if result.Data != tt.input {
				// 	t.Errorf("SaveData() data %v, want %v", result.Data, tt.input)
				// }
				// дополнительная проверка: убедимся, что сервис действительно вызвал стор с теми же даными, которые мы ему передали.
				if ms.saveHuman != tt.input {
					t.Errorf("store.SimpleSaveAtData был вызван с данными  %+v, ожидалось %+v", ms.saveHuman, tt.input)
				}
			}
		})
	}
}
