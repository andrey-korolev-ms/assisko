package store

import (
	"assisko/models"
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	SimpleSaveAtDB(ctx context.Context, human models.Human) error
	GetByID(ctx context.Context, id int) (models.Human, error)
	GetAll(ctx context.Context) ([]models.Human, error)
	// методы для пользователей
	CreateUser(ctx context.Context, username, password string) error
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
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
	_, err := s.db.ExecContext(ctx, "INSERT INTO assisko (sex, age, ipk, simpleemployment, isinvalid, invalidgroup, operatorid) values ($1, $2, $3, $4, $5, $6, $7)", human.Sex, human.Age, human.IPK, human.SimpleEmployment, human.IsInvalid, human.InvalidGroup, human.OperatorID)

	return err
}

func (s *PostgresStore) GetByID(ctx context.Context, id int) (models.Human, error) {
	var data models.Human
	err := s.db.QueryRowContext(ctx, "SELECT sex, age, ipk, simpleemployment, isinvalid, invalidgroup, operatorid FROM assisko WHERE id = $1", id).Scan(&data.Sex, &data.Age, &data.IPK, &data.SimpleEmployment, &data.IsInvalid, &data.InvalidGroup, &data.OperatorID)
	if err != nil {
		return models.Human{}, err
	}
	return data, nil
}

func (s *PostgresStore) GetAll(ctx context.Context) ([]models.Human, error) {

	rows, err := s.db.QueryContext(ctx, "SELECT sex, age, ipk, simpleemployment, isinvalid, invalidgroup, operatorid FROM assisko")
	if err != nil {
		return nil, err
	}

	var result []models.Human
	for rows.Next() {
		var data models.Human
		if err := rows.Scan(&data.Sex, &data.Age, &data.IPK, &data.SimpleEmployment, &data.IsInvalid, &data.InvalidGroup, &data.OperatorID); err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	return result, rows.Err()
}

func (s *PostgresStore) CreateUser(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, "INSERT INTO users (username, password_hash) values ($1, $2)", username, hashedPassword)
	return err
}
func (s *PostgresStore) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var user models.User
	err := s.db.QueryRowContext(ctx, "SELECT id, username, password_hash from users where username = $1", username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, errors.New("пользователь не найден")
	}
	return user, err
}

func (s *PostgresStore) Close() error {
	return s.db.Close()
}
