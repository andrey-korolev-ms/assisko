package services

import "assisko/models"

type mockStore struct{
	saveError error
	saveHuman models.Human
}

func (m *mockStore) 