package repository

import (
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
)

type StorageOrders struct {
	DBconnection *sql.DB
}

// создаем новый сторадж для работы с таблицей пользователей
func NewOrdersStorage(DBconn *sql.DB) StorageOrders {
	var so StorageOrders
	so.DBconnection = DBconn
	return so
}

// запрос данных по номеру заказа
func (so *StorageOrders) GetByNumber(number int64) model.Orders {
	var result model.Orders
	err := so.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	sqlStatement := `SELECT o_number, o_status, o_accrual, uploaded_at, created_by FROM orders
		WHERE o_number = $1`
	err = so.DBconnection.QueryRow(sqlStatement, number).Scan(&result.Number, &result.Status, &result.Accrual, &result.UploadedAt, &result.CreatedBy)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("order with number", number, "is not exist in DB")
	}
	return result
}

// запрос данных по логину пользователя
func (so *StorageOrders) GetByLogin(login string) []model.Orders {
	var result []model.Orders
	err := so.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	sqlStatement := `SELECT o_number, o_status, o_accrual, uploaded_at, created_by FROM orders
		WHERE created_by = $1 ORDER BY uploaded_at DESC`

	rows, err := so.DBconnection.Query(sqlStatement, login)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("orders created by user", login, "is not exist in DB")
	}
	if err != nil {
		logger.Log.Warnln("error while select from data by login from DB", err)
	}
	defer rows.Close()
	for rows.Next() {
		var o model.Orders
		if err := rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt, &o.CreatedBy); err != nil {
			logger.Log.Errorln(err)
		}
		result = append(result, o)
	}
	// Проверка на ошибки после цикла
	if err = rows.Err(); err != nil {
		logger.Log.Errorln(err)
	}
	return result
}

// обновление данных в БД
func (so *StorageOrders) UpdateOrder(order model.Orders) error {
	err := so.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	sqlStatement := `UPDATE orders 
		SET o_status = $1,
			o_accrual = $2
		WHERE o_number = $3;`
	_, err = so.DBconnection.Exec(sqlStatement, order.Status, order.Accrual, order.Number)
	if err != nil {
		logger.Log.Errorln("error while update order in DB", err)
		return err
	}
	return nil
}

// добавление данных о заказе в БД
func (so *StorageOrders) InsertOrder(order model.Orders) error {
	err := so.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	sqlStatement := `INSERT INTO orders 
			(o_number, o_status, o_accrual, uploaded_at, created_by) 
			VALUES ($1, $2, $3, $4, $5)`
	_, err = so.DBconnection.Exec(sqlStatement, order.Number, order.Status, order.Accrual, order.UploadedAt, order.CreatedBy)
	if err != nil {
		logger.Log.Errorln("error while insert order to db", err)
		return err
	}
	return nil
}
