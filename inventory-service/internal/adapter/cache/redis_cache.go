package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"inventory-service/internal/model"
	"time"
)

type RedisProductCache struct {
	client  *redis.Client
	ctx     context.Context
	prefix  string
	timeout time.Duration
}

func NewRedisProductCache() *RedisProductCache {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	return &RedisProductCache{
		client:  rdb,
		ctx:     context.Background(),
		prefix:  "product:",
		timeout: 24 * time.Hour,
	}
}

func (c *RedisProductCache) GetAll() []*model.Product {
	return []*model.Product{}
}

func (c *RedisProductCache) GetByID(id string) (*model.Product, bool) {
	data, err := c.client.Get(c.ctx, c.prefix+id).Result()
	if err != nil {
		return nil, false
	}
	var product model.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, false
	}
	fmt.Println("📦 REDIS HIT:", id)
	return &product, true
}

func (c *RedisProductCache) Save(product *model.Product) {
	payload, _ := json.Marshal(product)
	c.client.Set(c.ctx, c.prefix+product.ID, payload, c.timeout)
}

func (c *RedisProductCache) Update(product *model.Product) {
	c.Save(product)
}

func (c *RedisProductCache) Delete(id string) {
	c.client.Del(c.ctx, c.prefix+id)
}

func (c *RedisProductCache) LoadFromDB(products []*model.Product) {
	for _, product := range products {
		payload, _ := json.Marshal(product)
		c.client.Set(c.ctx, c.prefix+product.ID, payload, c.timeout)
	}
}
