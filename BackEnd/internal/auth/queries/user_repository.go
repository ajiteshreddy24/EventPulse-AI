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
	err := r.DB.QueryRow(`
		INSERT INTO users (name, email, password_hash, interests)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`,
		user.Name,
		user.Email,
		user.PasswordHash,
		serializeInterests(user.Interests),
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil && strings.Contains(strings.ToLower(err.Error()), "unique") {
		return ErrEmailAlreadyUsed
	}

	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	var interests string

	err := r.DB.QueryRow(`
		SELECT id, name, email, password_hash, interests, created_at
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&interests,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	user.Interests = parseInterests(interests)
	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	var interests string

	err := r.DB.QueryRow(`
		SELECT id, name, email, password_hash, interests, created_at
		FROM users
		WHERE id = $1
	`, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&interests,
		&user.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	user.Interests = parseInterests(interests)
	return &user, nil
}

func (r *UserRepository) UpdateInterests(userID int, interests []string) error {
	result, err := r.DB.Exec(`
		UPDATE users
		SET interests = $1
		WHERE id = $2
	`, serializeInterests(interests), userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

func parseInterests(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")
	interests := make([]string, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(strings.ToLower(part))
		if value != "" {
			interests = append(interests, value)
		}
	}

	return interests
}

func serializeInterests(interests []string) string {
	if len(interests) == 0 {
		return ""
	}

	return strings.Join(interests, ",")
}
