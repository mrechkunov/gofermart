package cryptoauth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrechkunov/gofermart/interal/logger"
)

const secretKey = "secret key"

func generateToken(uLogin string) (string, error) {
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
		return "", err
	}
	return tokenString, nil
}

// generate NEW UID, sign it and return cookieString
func GenerateNewToken() string {
	id := make([]byte, 4)
	_, err := rand.Read(id)
	if err != nil {
		logger.Log.Infoln("error while generate UID", err)
	}
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write(id)
	if err != nil {
		logger.Log.Infoln("error while sign UID", err)
	}
	sign := h.Sum(nil)
	result := append(id, sign...)
	return hex.EncodeToString(result)
}

// validate cookie signature
func ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return err
	}
	return nil
}

// get ID from cookie
func GetIDFromCookie(cookie string) (uint32, error) {
	data, err := hex.DecodeString(cookie)
	if err != nil {
		logger.Log.Infoln("error while decoding incoming cookie", err)
		return 0, err
	}
	return binary.BigEndian.Uint32(data[:4]), nil
}
