package order_repository

import (
	"assignment-2/entity"
	"assignment-2/pkg/errs"
)

type OrderRepository interface {
	CreateOrder(orderPayload entity.Order, itemsPayload []entity.Item) (*entity.Order, errs.MessageErr)
	GetAllOrders() ([]OrderItem, errs.MessageErr)
	UpdateOrder(orderPayload entity.Order, itemsPayload []entity.Item) (*OrderItem, errs.MessageErr)
	DeleteOrder(orderId int) errs.MessageErr
}
