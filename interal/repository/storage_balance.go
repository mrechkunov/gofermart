package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"strconv"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/logger"
	"github.com/mrechkunov/gofermart/interal/model"
)

type StorageBalance struct {
	DBconnection *sql.DB
}

// создаем новый сторадж для работы с таблицей пользователей
func NewBalanceStorage(DBconn *sql.DB) StorageBalance {
	return StorageBalance{DBconnection: DBconn}
}

// запрос данных по логину
func (sb *StorageBalance) GetBalanceByLogin(ctx context.Context, login string) model.Balance {
	var result model.Balance
	var curBal, withdrawnBal int64
	err := sb.DBconnection.QueryRowContext(ctx, "SELECT user_id, current_balance, withdrawn_balance, updated_at FROM balances WHERE user_id=$1", login).
		Scan(&result.UserID, &curBal, &withdrawnBal, &result.UpdatedAt)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("users balance is not exist in DB")
	}
	result.CurrentBalance = float64(curBal) / config.ExchangeRateCoefficient         // так как в бд храним и работаем с int преобразуем при запросе
	result.WithdrawnBalance = float64(withdrawnBal) / config.ExchangeRateCoefficient // так как в бд храним и работаем с int преобразуем при запросе
	return result
}

// запрос транзакций по логину
func (sb *StorageBalance) GetTransactionsByLogin(ctx context.Context, login string) []model.TransactionWithdraw {
	var result []model.TransactionWithdraw
	sqlStatement := `SELECT order_id, amount, created_at FROM transactions
		WHERE user_id = $1 AND withdraw is true ORDER BY created_at DESC`

	rows, err := sb.DBconnection.QueryContext(ctx, sqlStatement, login)
	if err == sql.ErrNoRows {
		logger.Log.Infoln("withdraw orders for user", login, "is not exist in DB")
	}
	if err != nil {
		logger.Log.Warnln("error while select transactions by login from DB", err)
	}
	defer rows.Close()
	for rows.Next() {
		var orderNumber int64
		var amount int64
		var created_at int64
		if err := rows.Scan(&orderNumber, &amount, &created_at); err != nil {
			logger.Log.Errorln(err)
		}
		order := model.TransactionWithdraw{
			OrderNumber: strconv.FormatInt(orderNumber, 10),
			Sum:         float64(amount) / -10000,
			ProcessedAt: time.Unix(0, created_at).Format(time.RFC3339),
		}
		result = append(result, order)
	}
	// Проверка на ошибки после цикла
	if err = rows.Err(); err != nil {
		logger.Log.Fatal(err)
	}
	return result
}

// добавление данных в БД
func (sb *StorageBalance) TransactionAdd(ctx context.Context, userID string, amount int64, orderID int64) error {
	//Начало
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := sb.DBconnection.BeginTx(ctxWithTimeout, nil)
	if err != nil {
		logger.Log.Errorln(err)
	}
	// generate NEW TransactionID
	id := make([]byte, 4)
	_, err = rand.Read(id)
	if err != nil {
		logger.Log.Infoln("error while generate TransactionID", err)
	}
	t_id := hex.EncodeToString(id)
	// Откат при ошибке (defer)
	defer func() {
		err := tx.Rollback()
		if err != nil {
			logger.Log.Infoln("transaction", t_id, err)
		}
	}()
	// запись в таблицу транзакций
	// подготовка данных
	created_at := time.Now().UnixNano()
	var withdraw = false
	if amount <= 0 {
		withdraw = true
	}
	sqlStatementTransactions := `INSERT INTO transactions 
			(user_id, t_id, amount, order_id, created_at, withdraw) 
			VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctxWithTimeout, sqlStatementTransactions, userID, t_id, amount, orderID, created_at, withdraw)
	if err != nil {
		logger.Log.Infoln("error while insert transaction", err)
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
	_, err = tx.ExecContext(ctxWithTimeout, sqlStatementBalance, amount, created_at, userID)
	if err != nil {
		logger.Log.Infoln("error while update balance", err)
		return err // Rollback
	}

	if amount < 0 {
		amountIntABS := amount * -1
		sqlStatementBalance := `UPDATE balances
				SET withdrawn_balance = withdrawn_balance + $1,
				updated_at = $2
				WHERE user_id = $3;`
		_, err = tx.ExecContext(ctxWithTimeout, sqlStatementBalance, amountIntABS, created_at, userID)
		if err != nil {
			logger.Log.Infoln("error while update withdrawn_balance", err)
			return err // Rollback
		}
	}

	//Фиксация транзакции
	if err := tx.Commit(); err != nil {
		logger.Log.Errorln(err)
		return err
	}
	return nil
}

func (sb *StorageBalance) AddUserBalance(ctx context.Context, login string) error {
	currentBalance := 0
	withdrawnBalance := 0
	updated_at := time.Now().UnixNano()
	sqlStatement := `INSERT INTO balances
			(user_id, current_balance, withdrawn_balance, updated_at) 
			VALUES ($1, $2, $3, $4)`
	_, err := sb.DBconnection.ExecContext(ctx, sqlStatement, login, currentBalance, withdrawnBalance, updated_at)
	if err != nil {
		logger.Log.Errorln("error while insert new user`s balance to db", err)
		return err
	}
	return nil
}
