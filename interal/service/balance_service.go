package service

import (
	"context"

	"github.com/mrechkunov/gofermart/interal/config"
	"github.com/mrechkunov/gofermart/interal/model"
	"github.com/mrechkunov/gofermart/interal/repository"
)

func GetBalanceByToken(ctx context.Context, token *string) model.Balance {
	storageUsers := repository.NewUsersStorage(config.DBconn)
	login := storageUsers.GetUserByToken(ctx, *token).Login
	storageBalance := repository.NewBalanceStorage(config.DBconn)
	return storageBalance.GetBalanceByLogin(ctx, login)
}
func GetTransactionsByToken(ctx context.Context, token *string) []model.TransactionWithdraw {
	storageUsers := repository.NewUsersStorage(config.DBconn)
	login := storageUsers.GetUserByToken(ctx, *token).Login
	storageBalance := repository.NewBalanceStorage(config.DBconn)
	return storageBalance.GetTransactionsByLogin(ctx, login)
}

func TransactionAdd(ctx context.Context, login *string, amount *int64, order *int64) error {
	storageBalance := repository.NewBalanceStorage(config.DBconn)
	err := storageBalance.TransactionAdd(ctx, *login, *amount, *order)
	if err != nil {
		return err
	}
	return nil
}
