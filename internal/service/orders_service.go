package service

import (
	"context"

	"github.com/mrechkunov/gofermart/internal/config"
	"github.com/mrechkunov/gofermart/internal/model"
	"github.com/mrechkunov/gofermart/internal/repository"
)

func GetOrderByNumber(ctx context.Context, number int64) model.Orders {
	storageOrders := repository.NewOrdersStorage(config.DBconn)
	return storageOrders.GetByNumber(ctx, number)
}

func InsertOrder(ctx context.Context, order *model.Orders) error {
	storageOrders := repository.NewOrdersStorage(config.DBconn)
	return storageOrders.InsertOrder(ctx, *order)
}

func GetOrdersSliceByLogin(ctx context.Context, login string) []model.Orders {
	storageOrders := repository.NewOrdersStorage(config.DBconn)
	return storageOrders.GetByLogin(ctx, login)
}
