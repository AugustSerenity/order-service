package consumer

import "github.com/AugustSerenity/order-service/internal/model"

type OrderService interface {
	ProcessOrder(order model.Order) error
}
