package service

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	authModels "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/models"
	authQueries "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/queries"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrInvalidAuthPayload = errors.New("name, email, and password are required")
var ErrInvalidLoginPayload = errors.New("email and password are required")

type AuthService struct {
	Repo *authQueries.UserRepository
}

func (s *AuthService) Register(req authModels.RegisterRequest) (*authModels.AuthResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return nil, ErrInvalidAuthPayload
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &authModels.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}

	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &authModels.AuthResponse{Token: token, User: sanitizeUser(*user)}, nil
}

func (s *AuthService) Login(req authModels.LoginRequest) (*authModels.AuthResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		return nil, ErrInvalidLoginPayload
	}

	user, err := s.Repo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, authQueries.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := generateToken(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &authModels.AuthResponse{Token: token, User: sanitizeUser(*user)}, nil
}

func (s *AuthService) GetUserByID(id int) (*authModels.User, error) {
	user, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	sanitized := sanitizeUser(*user)
	return &sanitized, nil
}

func (s *AuthService) ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret()), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func generateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   strconv.Itoa(userID),
		"email": email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret()))
}

func jwtSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "eventpulse-dev-secret"
	}
	return secret
}

func sanitizeUser(user authModels.User) authModels.User {
	user.PasswordHash = ""
	return user
}
