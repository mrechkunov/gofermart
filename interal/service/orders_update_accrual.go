package service

import "github.com/mrechkunov/gofermart/interal/logger"

func UpdateOrderWorker(chanToUpdate chan int64) {
	{
		for orderNumber := range chanToUpdate {
			logger.Log.Infoln("Order number to Update from accrual:", orderNumber)
		}
	}
}
