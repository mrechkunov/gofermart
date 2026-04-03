package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mrechkunov/gofermart/interal/cryptoauth"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/service"
)

const ErrUnauthorized = "401 Unauthorized"

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	var reqdata, user model.Users
	if err := json.NewDecoder(r.Body).Decode(&reqdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// проверяем пару логин пароль
	if service.GetUserByLogin(r.Context(), reqdata.Login).Password != cryptoauth.EncryptPass(reqdata.Password) {
		http.Error(w, ErrUnauthorized, http.StatusUnauthorized)
		return
	}
	// проверяем валидность текущего токена, если не валиден, генерируем новый
	if cryptoauth.ValidateToken(service.GetUserByLogin(r.Context(), reqdata.Login).Bearer) != nil {
		// генерируем token
		var err error
		user.Bearer, err = cryptoauth.GenerateToken(reqdata.Login)
		if err != nil {
			logger.Log.Error("error while generate token", err)
			return
		}
	}
	user.Login = reqdata.Login
	// шифруем переданный пароль
	user.Password = cryptoauth.EncryptPass(reqdata.Password)
	// записываем в БД данные пользователя
	err := service.UpdateUser(r.Context(), &user)
	if err != nil {
		logger.Log.Warnln("error while update user`s data in DB", err)
	}
	// формируем ответ
	w.Header().Set("Authorization", "Bearer "+user.Bearer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
