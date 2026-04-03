package cryptoauth

import (
	"context"
	"net/http"
	"strings"

	"github.com/mrechkunov/gofermart/interal/service"
)

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый http.Handler.
func WithAuth(h http.HandlerFunc) http.HandlerFunc {
	AuthFn := func(w http.ResponseWriter, r *http.Request) {
		authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		err := ValidateToken(authToken)
		if err != nil {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}
		user := service.GetUserByToken(r.Context(), authToken)
		if user.Login == "" {
			http.Error(w, "wrong token", http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "user", user))
	}
	return http.HandlerFunc(AuthFn)
}
