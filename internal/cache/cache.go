package cache

import (
	"time"

	"github.com/AugustSerenity/order-service/internal/model"
	"github.com/patrickmn/go-cache"
)

type Cache struct {
	pool *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		pool: cache.New(cache.NoExpiration, time.Minute),
	}
}

func (c *Cache) Add(id string, order model.Order) {
	c.pool.SetDefault(id, order)
}

func (c *Cache) Get(id string) (model.Order, bool) {
	content, ok := c.pool.Get(id)
	if !ok {
		return model.Order{}, false
	}

	order, ok := content.(model.Order)
	return order, ok
}

func (c *Cache) Fill(orders []model.Order) {
	for _, order := range orders {
		c.Add(order.OrderUID, order)
	}
}
