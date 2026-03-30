package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
)

type StorageBalance struct {
	DBconnection *sql.DB
}

// создаем новый сторадж для работы с таблицей пользователей
func NewBalanceStorage(DBconn *sql.DB) StorageBalance {
	var sb StorageBalance
	sb.DBconnection = DBconn
	return sb
}

// запрос данных по логину
func (sb *StorageBalance) GetByLogin(uLogin string) model.Balance {
	var result model.Balance
	err := sb.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	var curBal, withdrawnBal int
	err = sb.DBconnection.QueryRow("SELECT user_id, current_balance, withdrawn_balance, updated_at FROM balances WHERE user_id=$1", uLogin).
		Scan(&result.UserID, &curBal, &withdrawnBal, &result.Updated_at)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("users balance is not exist in DB")
	}

	fmt.Println("----------debug GET BY LOGIN--int------")
	fmt.Println("CurrentBalance", curBal)
	fmt.Println("Withdrawn_balance", withdrawnBal)
	fmt.Println("----------------------")

	result.CurrentBalance = float64(curBal)          // так как в бд храним и работаем с int преобразуем при запросе
	result.Withdrawn_balance = float64(withdrawnBal) // так как в бд храним и работаем с int преобразуем при запросе

	fmt.Println("----------debug GET BY LOGIN--float------")
	fmt.Println("CurrentBalance", result.CurrentBalance)
	fmt.Println("Withdrawn_balance", result.Withdrawn_balance)
	fmt.Println("----------------------")

	result.CurrentBalance = result.CurrentBalance / 100
	result.Withdrawn_balance = result.Withdrawn_balance / 100
	fmt.Println("----------debug GET BY LOGIN--float/100------")
	fmt.Println("CurrentBalance", result.CurrentBalance)
	fmt.Println("Withdrawn_balance", result.Withdrawn_balance)
	fmt.Println("----------------------")
	return result
}

// добавление данных в БД
func (sb *StorageBalance) TransactionAdd(ctx context.Context, userID string, amount int, orderID int64) error {
	err := sb.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	//Начало транзакции
	tx, err := sb.DBconnection.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Errorln(err)
	}
	// Откат при ошибке (defer)
	defer tx.Rollback()
	// запись в таблицу транзакций

	// подготовка данных
	// generate NEW TransactionID
	id := make([]byte, 4)
	_, err = rand.Read(id)
	if err != nil {
		logger.Log.Infoln("error while generate TransactionID", err)
	}
	t_id := hex.EncodeToString(id)
	//amount cast to int
	created_at := time.Now().UnixNano()

	sqlStatementTransactions := `INSERT INTO transactions 
			(user_id, t_id, amount, order_id, created_at) 
			VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, sqlStatementTransactions, userID, t_id, amount, orderID, created_at)
	if err != nil {
		return err // Rollback
	}

	// ----------------------------------
	// апдейт таблицы с балансом
	// подготовка данных
	// если amount отрицательный то делаем два апдейта
	sqlStatementBalance := `UPDATE balances 
				SET current_balance = current_balance + $1,
				updated_at = $2 
				WHERE user_id = $3;`
	_, err = tx.ExecContext(ctx, sqlStatementBalance, amount, created_at, userID)
	if err != nil {
		return err // Rollback
	}

	if amount < 0 {
		amountIntABS := int(math.Abs(float64(amount)))
		sqlStatementBalance := `UPDATE balances 
				SET withdrawn_balance = withdrawn_balance + $1,
				updated_at = $2
				WHERE user_id = $3;`
		_, err = tx.ExecContext(ctx, sqlStatementBalance, amountIntABS, created_at, userID)
		if err != nil {
			return err // Rollback
		}
	}

	//Фиксация транзакции
	if err := tx.Commit(); err != nil {
		logger.Log.Errorln(err)
		return err
	}

	logger.Log.Infoln("Transaction added")
	return nil
}
func (sb *StorageBalance) AddUserBalance(login string) error {
	err := sb.DBconnection.Ping()
	if err != nil {
		logger.Log.Warnln(err)
	}
	currentBalance := 0
	withdrawnBalance := 0
	updated_at := time.Now().UnixNano()
	sqlStatement := `INSERT INTO balances
			(user_id, current_balance, withdrawn_balance, updated_at) 
			VALUES ($1, $2, $3, $4)`
	_, err = sb.DBconnection.Exec(sqlStatement, login, currentBalance, withdrawnBalance, updated_at)
	if err != nil {
		logger.Log.Errorln("error while insert new user`s balance to db", err)
		return err
	}
	return nil
}

func (sb *StorageBalance) Close() error {
	return sb.DBconnection.Close()
}
