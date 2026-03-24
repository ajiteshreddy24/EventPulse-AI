package queries

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/models"
)

var ErrUserNotFound = errors.New("user not found")
var ErrEmailAlreadyUsed = errors.New("email already exists")

type UserRepository struct {
	DB *sql.DB
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
	INSERT INTO users (name, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	err := r.DB.QueryRow(query, user.Name, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt)
	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return ErrEmailAlreadyUsed
	}

	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
	SELECT id, name, email, password_hash, created_at
	FROM users
	WHERE email = $1`

	var user models.User
	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
	SELECT id, name, email, password_hash, created_at
	FROM users
	WHERE id = $1`

	var user models.User
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
