package cryptoauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrechkunov/gofermart/interal/logger"
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
	// Создаем токен с методом подписи HMAC (HS256)
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
		logger.Log.Infoln("tocken is not valid", err)
		return err
	}
	return nil
}

// encrypt the password
func EncryptPass(password string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	_, err := h.Write([]byte(password))
	if err != nil {
		logger.Log.Infoln("error while encrypt password", err)
	}
	encryptedPassword := h.Sum(nil)
	return hex.EncodeToString(encryptedPassword)
}
