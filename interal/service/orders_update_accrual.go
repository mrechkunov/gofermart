package service

import (
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

func UpdateOrderListener(chanToUpdate chan int64) {
	{
		for orderNumber := range chanToUpdate {
			logger.Log.Infoln("Order number to Update from accrual:", orderNumber)
			// создаем группу
			var wg sync.WaitGroup

			wg.Add(1)                               // добавляем в группу запуск горутины
			go UpdateOrderAccrual(orderNumber, &wg) // запускаем горутину на апдейт до конечных статусов

			// // записываем номер заказа в слайс //мапу key - orderNumber, value - status(done/run)

			// // идем к accrual получаем данные по номеру заказа

			// orderFromAccrual, err := GetAccrual(orderNumber)
			// ErrNoContent := errors.New("204 no content")
			// ErrToManyRequests := errors.New("429 too many requests")

			// if err == ErrToManyRequests {
			// 	time.Sleep(10 * time.Second)
			// }

			// fmt.Println("data from Accrual", orderFromAccrual.Order, orderFromAccrual.Status, orderFromAccrual.Accrual)
			// если accrual просит подождать, ждем и идем еще раз
			// идем в бд, получаем статус по номеру заказа
			// сравниваем оба статуса, если accrual отличный от БД, обновляем БД статус и начисления
			// если статус из accrual пришел конечный (imvalid/processed), то удаляем номер заказа из слайса,
			// иначе переходим к следующему заказу
			wg.Wait()
		}
	}
}

func UpdateOrderAccrual(number int64, wg *sync.WaitGroup) {
	defer wg.Done()
	url := config.ConfigAddresses.AccuralSystemAddress + "/api/orders/" + strconv.FormatInt(number, 10)
	logger.Log.Infoln("string uri to accrual:", url)
	isDone := false
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
		}

		if resp.StatusCode != http.StatusOK {
			logger.Log.Infoln("API request failed with status:", resp.Status)
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
			logger.Log.Infoln("error while conver order number srom string to int64", err)
		}
		if respOrder.Status != "REGISTERED" {
			order.Status = respOrder.Status
		}

		order.Accrual = int(respOrder.Accrual * 100) // храним в БД сумму в копейках
		storageOrders := repository.NewOrdersStorage(config.DBconn)
		storageOrders.UpdateOrder(order)
	}
}
