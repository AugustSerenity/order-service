package service

import (
	"encoding/json"
	"fmt"

	"github.com/AugustSerenity/order-service/internal/cache"
	"github.com/AugustSerenity/order-service/internal/model"
	"github.com/sirupsen/logrus"
)

type Service struct {
	cache   *cache.Cache
	storage Storage
}

func NewOrderService(cache *cache.Cache, st Storage) *Service {
	return &Service{
		cache:   cache,
		storage: st,
	}
}

func (s *Service) RestoreCache() error {
	orders, err := s.storage.GetAllOrders()
	if err != nil {
		return err
	}

	s.cache.Fill(orders)
	logrus.Infof("Cache restored with %d orders", len(orders))
	return nil
}

func (s *Service) ProcessOrder(order model.Order) error {
	if err := s.storage.SaveOrder(order); err != nil {
		logrus.WithError(err).Error("Failed to save order to DB")
		return fmt.Errorf("failed to save order: %w", err)
	}

	s.cache.Add(order.OrderUID, order)

	logrus.Infof("Order %s successfully saved and cached", order.OrderUID)
	return nil
}

func (s *Service) GetOrderByID(id string) ([]byte, bool) {
	order, foundInCache := s.cache.Get(id)
	if foundInCache {
		data, err := json.Marshal(order)
		if err != nil {
			logrus.WithError(err).Error("failed to marshal order from cache")
			return nil, false
		}
		return data, true
	}

	order, err := s.storage.GetByID(id)
	if err != nil {
		logrus.WithError(err).Error("failed to get order from DB")
		return nil, false
	}

	s.cache.Add(id, order)

	data, err := json.Marshal(order)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal order from DB")
		return nil, false
	}

	return data, true
}
