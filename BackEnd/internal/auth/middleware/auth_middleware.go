package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	authService "github.com/ajiteshreddy24/EventPulse-AI/BackEnd/internal/auth/service"
)

type contextKey string

const userIDContextKey contextKey = "userID"

type AuthMiddleware struct {
	Service *authService.AuthService
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := m.Service.ParseToken(tokenString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		subject, ok := claims["sub"].(string)
		if !ok {
			http.Error(w, "invalid token subject", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(subject)
		if err != nil {
			http.Error(w, "invalid token subject", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDContextKey).(int)
	return userID, ok
}
