package cryptoauth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/mrechkunov/gofermart/interal/logger"
)

var secretKey = "secret key"

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
func ValidateTokenSign(token string) error {
	h := hmac.New(sha256.New, []byte(secretKey))
	data, err := hex.DecodeString(token)
	if err != nil {
		logger.Log.Warnln("error while decoding incoming token", err)
		return err
	}
	h.Write([]byte(data[:4]))
	sign := h.Sum(nil)
	if hmac.Equal(sign, data[4:]) {
		return nil
	}
	return errors.New("not valid token signature")
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
