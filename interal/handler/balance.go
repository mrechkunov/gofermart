package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/cryptoauth"
	"github.com/mrechkunov/gofermart/interal/repository"
)

func Balance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	// читаем Header Autorization и записываем его в поле token
	authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	err := cryptoauth.ValidateToken(authToken)
	if err != nil {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}
	storageUsers := repository.NewUsersStorage(config.DBconn)
	login := storageUsers.GetByToken(authToken).Login
	storageBalance := repository.NewBalanceStorage(config.DBconn)
	balance := storageBalance.GetByLogin(login)

	fmt.Println("----------debug balance handler--------")
	fmt.Println("CurrentBalance", balance.CurrentBalance)
	fmt.Println("Updated_at", balance.Updated_at)
	fmt.Println("UserID", balance.UserID)
	fmt.Println("Withdrawn_balance", balance.Withdrawn_balance)
	fmt.Println("----------------------")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// записываем ответ
	json.NewEncoder(w).Encode(balance)
}
