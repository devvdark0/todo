package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/devvdark0/todo/internal/service"
	"github.com/google/uuid"
)

func AuthMiddleware(authService service.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			userIDStr, err := claims.GetSubject()
			if err != nil {
				http.Error(w, "Invalid token claims", http.StatusUnauthorized)
				return
			}

			userId, err := uuid.Parse(userIDStr)
			if err != nil {
				http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", userId.String())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

func GetUserID(r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value("userId").(uuid.UUID)
	return userID, ok
}
