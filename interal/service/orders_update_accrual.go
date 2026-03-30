package service

import (
	"fmt"

	"github.com/mrechkunov/gofermart/interal/logger"
)

func UpdateOrderWorker(chanToUpdate chan int64) {
	{
		for orderNumber := range chanToUpdate {
			logger.Log.Infoln("Order number to Update from accrual:", orderNumber)
			// записываем номер заказа в слайс //мапу key - orderNumber, value - status(done/run)

			// идем к accrual получаем данные по номеру заказа

			orderFromAccrual := GetAccrual(orderNumber)

			fmt.Println("data from Accrual", orderFromAccrual.Order, orderFromAccrual.Status, orderFromAccrual.Accrual)
			// если accrual просит подождать, ждем и идем еще раз
			// идем в бд, получаем статус по номеру заказа
			// сравниваем оба статуса, если accrual отличный от БД, обновляем БД статус и начисления
			// если статус из accrual пришел конечный (imvalid/processed), то удаляем номер заказа из слайса,
			// иначе переходим к следующему заказу
		}
	}
}
