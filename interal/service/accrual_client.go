package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
)

func GetAccrual(number int64) model.AccrualOrder {
	url := config.ConfigAddresses.AccuralSystemAddress + "/" + strconv.FormatInt(number, 10)
	logger.Log.Infoln("string uri to accrual:", url)

	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Warnln("Error making GET request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Log.Infoln("API request failed with status:", resp.Status)
	}
	var result model.AccrualOrder
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		logger.Log.Warnln("Error decoding JSON:", err)
	}
	return result
}
