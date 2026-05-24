package service

import (
	"context"

	"github.com/mrechkunov/gofermart/internal/config"
	"github.com/mrechkunov/gofermart/internal/model"
	"github.com/mrechkunov/gofermart/internal/repository"
)

func GetBalanceByLogin(ctx context.Context, login string) model.Balance {
	storageBalance := repository.NewBalanceStorage(config.DBconn)
	return storageBalance.GetBalanceByLogin(ctx, login)
}
func GetTransactionsByLogin(ctx context.Context, login string) []model.TransactionWithdraw {
	storageBalance := repository.NewBalanceStorage(config.DBconn)
	return storageBalance.GetTransactionsByLogin(ctx, login)
}

func TransactionAdd(ctx context.Context, login string, amount int64, order int64) error {
	storageBalance := repository.NewBalanceStorage(config.DBconn)
	err := storageBalance.TransactionAdd(ctx, login, amount, order)
	if err != nil {
		return err
	}
	return nil
}
