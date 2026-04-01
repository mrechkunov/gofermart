package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/cryptoauth"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/repository"
)

func OrdersPost(c chan int64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
			return
		}
		storageUsers := repository.NewUsersStorage(config.DBconn)
		var user model.Users
		// читаем Header Autorization и записываем его в поле token
		user.Bearer = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		err := cryptoauth.ValidateToken(user.Bearer)
		if err != nil {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}
		user = storageUsers.GetUserByToken(r.Context(), user.Bearer)
		//читаем тело запроса
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "body reading error", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
		var incomeOrder model.Orders
		incomeOrder.Number, err = strconv.ParseInt(string(body), 10, 64)
		if err != nil {
			logger.Log.Infoln("error during conversion string number to int64:", err)
			return
		}
		incomeOrder.CreatedBy = user.Login
		incomeOrder.UploadedAt = time.Now().UnixNano()
		incomeOrder.Status = "NEW"
		storageOrders := repository.NewOrdersStorage(config.DBconn)
		orderFromDB := storageOrders.GetByNumber(r.Context(), incomeOrder.Number)
		if !cryptoauth.ValidLuhnOrderNumber(incomeOrder.Number) {
			http.Error(w, "error invalid order number format", http.StatusUnprocessableEntity)
			return
		}
		if orderFromDB.CreatedBy == incomeOrder.CreatedBy {
			http.Error(w, "error order number has already been uploaded by this user.", http.StatusOK)
			return
		}
		if orderFromDB.CreatedBy != incomeOrder.CreatedBy && orderFromDB.Number == incomeOrder.Number {
			http.Error(w, "error order number has already been uploaded by another user", http.StatusConflict)
			return
		}
		c <- incomeOrder.Number
		storageOrders.InsertOrder(r.Context(), incomeOrder)
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
	}
}

func OrdersGet(w http.ResponseWriter, r *http.Request) {
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
	storageOrders := repository.NewOrdersStorage(config.DBconn)
	storageUsers := repository.NewUsersStorage(config.DBconn)
	login := storageUsers.GetUserByToken(r.Context(), authToken).Login
	ordersFromDB := storageOrders.GetByLogin(r.Context(), login)

	if len(ordersFromDB) == 0 {
		http.Error(w, "no data in DB", http.StatusNoContent)
		return
	}
	// формируем батч, отправляем
	var respOrders []model.ResponceOrders
	for _, orderFromDB := range ordersFromDB {
		var respOrder model.ResponceOrders
		respOrder.Number = strconv.FormatInt(orderFromDB.Number, 10)
		respOrder.Status = orderFromDB.Status
		respOrder.Accrual = float64(orderFromDB.Accrual) / 100
		respOrder.UploadedAt = time.Unix(0, orderFromDB.UploadedAt).Format(time.RFC3339)
		respOrders = append(respOrders, respOrder)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// записываем ответ
	err = json.NewEncoder(w).Encode(respOrders)
	if err != nil {
		logger.Log.Warnln("error while encoding json in OrdersGet handler", err)
	}
}
