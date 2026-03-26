package handler

import (
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/mrechkunov/gofermart/interal/config"
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
	Bearer := r.Header.Get("Authorization")
	if Bearer != "" {
		// проверить что токен активный, вернуть что пользователь уже зарегистрирован
	}
	// читаем тело запроса
	var reqdata model.Users
	if err := json.NewDecoder(r.Body).Decode(&reqdata); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// проверяем свободный ли логин
	dbData := storage.GetByLogin(reqdata.Login)

	fmt.Println("Bearer:", Bearer)
	fmt.Println("dbData:", dbData)
	fmt.Println("reqdata:", reqdata)
	// storageUsers := repository.NewUsersStorage(config.DBconn)
	// storageUsers.InsertRow()
	w.Header().Set("Authorization", "Bearer FSDfrefsdFSDferf")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//записываем ответ
	json.NewEncoder(w).Encode("")

}
