package middleware

import (
	"log"
	"net/http"
)

// AuthMiddleware Middleware функция, которая оборачивает обработчики
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логика перед вызовом основного обработчика

		userId := r.Header.Get("X-User-Id")
		if userId == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user id is empty"))
			log.Printf("user id is empty")
			return
		}

		// Вызов основного обработчика
		next(w, r)

	}
}
