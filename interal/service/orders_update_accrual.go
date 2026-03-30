package service

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/repository"
)

func UpdateOrderListener(ctx context.Context, chanToUpdate chan int64) {
	{
		ctxWithCancel := context.WithoutCancel(ctx)
		// создаем группу
		var wg sync.WaitGroup
		for orderNumber := range chanToUpdate {
			wg.Add(1)                                                    // добавляем в группу запуск горутины
			go UpdateOrderAccrualWorker(ctxWithCancel, orderNumber, &wg) // запускаем горутину на апдейт до конечных статусов
		}
		// ждем всех из группы
		wg.Wait()
	}
}

func UpdateOrderAccrualWorker(ctx context.Context, number int64, wg *sync.WaitGroup) {
	defer wg.Done() // сообщим в группу что мы закончили, когда закончим
	url := config.ConfigAddresses.AccuralSystemAddress + "/api/orders/" + strconv.FormatInt(number, 10)
	isDone := false // это что бы
	for !isDone {
		resp, err := http.Get(url)
		if err != nil {
			logger.Log.Warnln("Error making GET request:", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusNoContent {
			logger.Log.Infoln("accrual responce no content, sleep 3 sec")
			time.Sleep(3 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusTooManyRequests {
			tts := resp.Header.Get("Retry-After")
			timeToSleep, err := strconv.Atoi(tts)
			if err != nil {
				logger.Log.Warnln("error while convert Retry-After string to int")
			}
			time.Sleep(time.Duration(timeToSleep) * time.Second)
			logger.Log.Infoln("accrual responce to many requests, sleep", timeToSleep, " sec")
			continue
		}

		if resp.StatusCode != http.StatusOK {
			logger.Log.Infoln("API request failed with status:", resp.Status)
			continue
		}

		var respOrder model.AccrualOrder
		err = json.NewDecoder(resp.Body).Decode(&respOrder)
		if err != nil {
			logger.Log.Warnln("Error decoding JSON:", err)
		}
		if respOrder.Status == "PROCESSED" || respOrder.Status == "INVALID" {
			isDone = true
		}

		var order model.Orders
		order.Number, err = strconv.ParseInt(respOrder.Order, 10, 64)
		if err != nil {
			logger.Log.Infoln("error while conver order number from string to int64(updateOrderAccrualWorker)", err)
		}
		if respOrder.Status != "REGISTERED" {
			order.Status = respOrder.Status
		}
		order.Accrual = int(respOrder.Accrual * 100) // храним в БД сумму в копейках
		if order.Accrual > 0 {
			storageBalance := repository.NewBalanceStorage(config.DBconn)
			err := storageBalance.TransactionAdd(ctx, order.CreatedBy, order.Accrual, order.Number)
			if err != nil {
				logger.Log.Warnln("error while transaction add at order update", err)
			}
		}
		storageOrders := repository.NewOrdersStorage(config.DBconn)
		storageOrders.UpdateOrder(order)
	}
}
