package service

import "github.com/AugustSerenity/order-service/internal/model"

type Storage interface {
	GetAllOrders() ([]model.Order, error)
	SaveOrder(order model.Order) error
	GetByID(id string) (model.Order, error)
}
