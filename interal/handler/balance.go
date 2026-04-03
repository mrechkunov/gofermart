package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/mrechkunov/gofermart/interal/cryptoauth"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/service"
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
	balance := service.GetBalanceByToken(r.Context(), &authToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// записываем ответ
	err = json.NewEncoder(w).Encode(balance)
	if err != nil {
		logger.Log.Warnln("error while encoding json in balance handler", err)
	}
}

func Withdraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
		return
	}
	// читаем Header Autorization и записываем его в поле token
	authToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	err := cryptoauth.ValidateToken(authToken)
	if err != nil {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}
	user := service.GetUserByToken(r.Context(), &authToken)
	// читаем тело запроса
	var withdrawOrder model.WithdrawOrder
	if err := json.NewDecoder(r.Body).Decode(&withdrawOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	incomeOrderNumber, err := strconv.ParseInt(string(withdrawOrder.Order), 10, 64)
	if !cryptoauth.ValidLuhnOrderNumber(&incomeOrderNumber) {
		http.Error(w, "error invalid order number format", http.StatusUnprocessableEntity)
		return
	}
	userCurrentBalance := service.GetBalanceByToken(r.Context(), &user.Bearer).CurrentBalance
	if withdrawOrder.Sum > userCurrentBalance {
		http.Error(w, "error insufficient funds in the account", http.StatusPaymentRequired)
		return
	}
	amountInt := int64(withdrawOrder.Sum * -10000)
	orderInt, err := strconv.ParseInt(withdrawOrder.Order, 10, 64)
	if err != nil {
		logger.Log.Errorln("error while convert order number fron string to int64 (withdraw)")
	}

	err = service.TransactionAdd(r.Context(), &user.Login, &amountInt, &orderInt)
	if err != nil {
		logger.Log.Errorln("error while transaction add at withdraw", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func Withdrawals(w http.ResponseWriter, r *http.Request) {
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
	withdrawals := service.GetTransactionsByToken(r.Context(), &authToken)
	if len(withdrawals) == 0 {
		http.Error(w, "no withdrawals in DB", http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// записываем ответ
	err = json.NewEncoder(w).Encode(withdrawals)
	if err != nil {
		logger.Log.Warnln("error while encoding json in withdrawals handler", err)
	}
}
