package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/mrechkunov/gofermart/internal/config"
	"github.com/mrechkunov/gofermart/internal/cryptoauth"
	"github.com/mrechkunov/gofermart/internal/logger"
	"github.com/mrechkunov/gofermart/internal/model"
	"github.com/mrechkunov/gofermart/internal/service"
)

func OrdersPost(c chan int64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusBadRequest)
			return
		}
		user := r.Context().Value("user").(model.Users)
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
		orderFromDB := service.GetOrderByNumber(r.Context(), incomeOrder.Number)
		if !cryptoauth.ValidLuhnOrderNumber(&incomeOrder.Number) {
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
		service.InsertOrder(r.Context(), &incomeOrder)
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
	}
}

func OrdersGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	user := r.Context().Value("user").(model.Users)
	ordersFromDB := service.GetOrdersSliceByLogin(r.Context(), user.Login)
	if len(ordersFromDB) == 0 {
		http.Error(w, "no data in DB", http.StatusNoContent)
		return
	}
	// формируем батч, отправляем
	var respOrders []model.ResponceOrders
	for _, orderFromDB := range ordersFromDB {
		respOrder := model.ResponceOrders{
			Number:     strconv.FormatInt(orderFromDB.Number, 10),
			Status:     orderFromDB.Status,
			Accrual:    float64(orderFromDB.Accrual) / config.ExchangeRateCoefficient,
			UploadedAt: time.Unix(0, orderFromDB.UploadedAt).Format(time.RFC3339),
		}
		respOrders = append(respOrders, respOrder)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// записываем ответ
	err := json.NewEncoder(w).Encode(respOrders)
	if err != nil {
		logger.Log.Warnln("error while encoding json in OrdersGet handler", err)
	}
}
