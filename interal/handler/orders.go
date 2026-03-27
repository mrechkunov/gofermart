package handler

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/cryptoauth"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/repository"
)

func Orders(w http.ResponseWriter, r *http.Request) {
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
	user = storageUsers.GetByToken(user.Bearer)
	//читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Body reading error", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var incomeOrder model.Orders
	incomeOrder.Number, err = strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		fmt.Printf("Error during conversion: %v\n", err)
		return
	}
	incomeOrder.CreatedBy = user.Login
	incomeOrder.UploadedAt = time.Now().Truncate(time.Second).Unix()
	incomeOrder.Status = "NEW"
	storageOrders := repository.NewOrdersStorage(config.DBconn)
	orderFromDB := storageOrders.GetByNumber(incomeOrder.Number)
	if !cryptoauth.ValidLuhnOrderNumber(incomeOrder.Number) {
		http.Error(w, "неверный формат номера заказа", http.StatusUnprocessableEntity)
		return
	}
	if orderFromDB.CreatedBy == incomeOrder.CreatedBy {
		http.Error(w, "номер заказа уже был загружен этим пользователем", http.StatusOK)
		return
	}
	if orderFromDB.CreatedBy != incomeOrder.CreatedBy && orderFromDB.Number == incomeOrder.Number {
		http.Error(w, "номер заказа уже был загружен другим пользователем", http.StatusConflict)
		return
	}
	storageOrders.InsertOrder(incomeOrder)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
}
