package handler

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/cryptoauth"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/repository"
)

func Register(ctx context.Context) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
			return
		}
		storageUsers := repository.NewUsersStorage(config.DBconn)
		// читаем тело запроса
		var reqdata, user model.Users
		if err := json.NewDecoder(r.Body).Decode(&reqdata); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// проверяем свободный ли логин если занят вернуть ошибку
		if storageUsers.GetUserByLogin(ctx, reqdata.Login).Login == reqdata.Login {
			http.Error(w, "login "+reqdata.Login+" is exist in DB", http.StatusConflict)
			return
		}
		user.Login = reqdata.Login
		// шифруем переданный пароль
		user.Password = cryptoauth.EncryptPass(reqdata.Password)
		// генерируем token
		tokenString, err := cryptoauth.GenerateToken(reqdata.Login)
		if err != nil {
			logger.Log.Error("error while generate token", err)
			return
		}
		user.Bearer = tokenString
		// записываем в БД данные пользователя
		storageUsers.InsertUser(ctx, user)
		storageBalance := repository.NewBalanceStorage(config.DBconn)
		storageBalance.AddUserBalance(ctx, user.Login)
		// формируем ответ
		w.Header().Set("Authorization", "Bearer "+tokenString)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
