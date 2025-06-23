package cache

import (
	"inventory-service/internal/model"
	"inventory-service/internal/usecase"
)

type MultiCache struct {
	mem   usecase.ProductCache
	redis usecase.ProductCache
}

func NewMultiCache(mem, redis usecase.ProductCache) usecase.ProductCache {
	return &MultiCache{mem: mem, redis: redis}
}

func (c *MultiCache) GetAll() []*model.Product {
	return c.mem.GetAll()
}

func (c *MultiCache) GetByID(id string) (*model.Product, bool) {
	if p, ok := c.mem.GetByID(id); ok {
		return p, true
	}
	if p, ok := c.redis.GetByID(id); ok {
		c.mem.Save(p)
		return p, true
	}
	return nil, false
}

func (c *MultiCache) Save(p *model.Product) {
	c.mem.Save(p)
	c.redis.Save(p)
}

func (c *MultiCache) Update(p *model.Product) {
	c.mem.Update(p)
	c.redis.Update(p)
}

func (c *MultiCache) Delete(id string) {
	c.mem.Delete(id)
	c.redis.Delete(id)
}

func (c *MultiCache) LoadFromDB(products []*model.Product) {
	c.mem.LoadFromDB(products)
	c.redis.LoadFromDB(products)
}
