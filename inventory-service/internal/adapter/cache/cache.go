package cache

import (
	"inventory-service/internal/model"
	"sync"
	"time"
)

type InMemoryProductCache struct {
	mu       sync.RWMutex
	products map[string]*model.Product
	lastLoad time.Time
}

func NewInMemoryProductCache() *InMemoryProductCache {
	return &InMemoryProductCache{
		products: make(map[string]*model.Product),
		lastLoad: time.Now(),
	}
}

func (c *InMemoryProductCache) GetAll() []*model.Product {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]*model.Product, 0, len(c.products))
	for _, p := range c.products {
		result = append(result, p)
	}
	return result
}

func (c *InMemoryProductCache) GetByID(id string) (*model.Product, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, exists := c.products[id]
	return p, exists
}

func (c *InMemoryProductCache) Save(p *model.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.products[p.ID] = p
}

func (c *InMemoryProductCache) Update(p *model.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.products[p.ID] = p
}

func (c *InMemoryProductCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.products, id)
}

func (c *InMemoryProductCache) LoadFromDB(products []*model.Product) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.products = make(map[string]*model.Product)
	for _, p := range products {
		c.products[p.ID] = p
	}
	c.lastLoad = time.Now()
}
