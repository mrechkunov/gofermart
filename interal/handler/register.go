package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"net/http"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/cryptoauth"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/repository"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	storage := repository.NewUsersStorage(config.DBconn)
	// читаем Header Autorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
		return
	}
	// Ожидаем формат "Bearer <token>" проверяем на валидность, если не валиден юзеру надо авторизироваться
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	err := cryptoauth.ValidateToken(tokenString)
	if err != nil {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	// читаем тело запроса
	var reqdata model.Users
	if err := json.NewDecoder(r.Body).Decode(&reqdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// проверяем свободный ли логин если занят вернуть ошибку
	dbData := storage.GetByLogin(reqdata.Login)
	if dbData.Login == reqdata.Login {
		http.Error(w, "login is exist", http.StatusConflict)
		return
	}
	// генерируем token и шифруем пароль

	fmt.Println("Bearer:", tokenString)
	fmt.Println("dbData:", dbData)
	fmt.Println("reqdataLogin:", reqdata.Login, "reqdataPassword", reqdata.Password)

	// storageUsers := repository.NewUsersStorage(config.DBconn)
	// storageUsers.InsertRow()
	w.Header().Set("Authorization", "Bearer FSDfrefsdFSDferf")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//записываем ответ
	json.NewEncoder(w).Encode("")

}
