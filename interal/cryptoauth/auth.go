package cryptoauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/service"
)

const secretKey = "secret key"

// generate and sign token
func GenerateToken(uLogin string) (string, error) {
	// Создаем claims (данные токена)
	claims := jwt.MapClaims{
		"username": uLogin,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // Срок действия 2 часа
		"iat":      time.Now().Unix(),
	}
	// Создаем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Подписываем токен
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		logger.Log.Errorln("error while sign token", err)
		return "", err
	}
	return tokenString, nil
}

// validate token signature
func ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		logger.Log.Infoln(tokenString, err)
		return err
	}
	return nil
}

// encrypt the password
func EncryptPass(password string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	_, err := h.Write([]byte(password))
	if err != nil {
		logger.Log.Errorln("error while encrypt password", err)
	}
	encryptedPassword := h.Sum(nil)
	return hex.EncodeToString(encryptedPassword)
}

// проверяет номер заказа по алгоритму Луна
func ValidLuhnOrderNumber(num *int64) bool {
	number := strconv.FormatInt(*num, 10)
	// убираем все пробелы в строке
	number = strings.ReplaceAll(number, " ", "")
	// проверяем что больше 2-х цифр
	if len(number) <= 1 {
		return false
	}
	sum := 0
	// проходим слева направо
	for i := len(number) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false // если не цифра
		}
		// Удваиваем каждую вторую цифру начиная с самой правой -1
		if (len(number)-1-i)%2 == 1 {
			digit *= 2
			if digit > 9 {
				// если удвоение двухзначное вычитаем 9
				digit -= 9
			}
		}
		sum += digit
	}
	// номер заказа валиден если сумма делится без остатка на 10
	return sum%10 == 0
}

// WithAuth добавляет дополнительный код для авторизации пользователя, записывает в контекст структуру авторизированного пользователя
// и возвращает новый http.Handler.
func WithAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		h.ServeHTTP(w, r)
	}
}
