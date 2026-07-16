package middleware

import (
	"context"
	"net/http"
	"strconv"

	"task-management-api/internal/models"
	"task-management-api/internal/storage"
)

type contextKey string

const UserContextKey contextKey = "user"

func UserMiddleware(store *storage.MemoryStore) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			header := r.Header.Get("X-User-Id")

			if header == "" {
				http.Error(w, "missing X-User-Id header", http.StatusBadRequest)
				return
			}

			userID, err := strconv.Atoi(header)
			if err != nil {
				http.Error(w, "invalid X-User-Id", http.StatusBadRequest)
				return
			}

			user, ok := store.GetUser(userID)
			if !ok {
				http.Error(w, "user not found", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(r *http.Request) *models.User {

	user, ok := r.Context().Value(UserContextKey).(*models.User)

	if !ok {
		return nil
	}

	return user
}
