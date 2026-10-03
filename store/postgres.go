package store

import (
	"assisko/models"
	"context"
	"database/sql"
)

type Storage interface {
	SimpleSaveAtDB(ctx context.Context, human models.Human) error
	GetByID(ctx context.Context, id int) (models.Human, error)
	GetAll(ctx context.Context) ([]models.Human, error)
}
type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(ctx context.Context, conn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) SimpleSaveAtDB(ctx context.Context, human models.Human) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO human () values ()", human.Sex)

	return err
}
