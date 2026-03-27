package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/cryptoauth"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/repository"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	storageUsers := repository.NewUsersStorage(config.DBconn)
	var reqdata, user model.Users
	// читаем Header Autorization и записываем его в поле token
	user.Bearer = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	// if authHeader == "" {
	// 	http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
	//  	return
	//  }
	// // Ожидаем формат "Bearer <token>" проверяем на валидность, если не валиден юзеру надо авторизироваться
	// tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	// err := cryptoauth.ValidateToken(tokenString)
	// if err != nil {
	// 	http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	// читаем тело запроса
	if err := json.NewDecoder(r.Body).Decode(&reqdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// проверяем пару логин пароль
	if storageUsers.GetByLogin(reqdata.Login).Password != cryptoauth.EncryptPass(reqdata.Password) {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}
	// проверяем валидность текущего токена, если не валиден, генерируем новый
	if cryptoauth.ValidateToken(storageUsers.GetByLogin(reqdata.Login).Bearer) != nil {
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
	storageUsers.UpdateUser(user)
	// формируем ответ
	w.Header().Set("Authorization", "Bearer "+user.Bearer)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
