package cache

import (
	"log"
	"order-service/internal/model"
	"sync"
)

type InMemoryOrderCache struct {
	byID     map[string]*model.Order
	byUserID map[string][]*model.Order
	mu       sync.RWMutex
}

func NewInMemoryOrderCache() *InMemoryOrderCache {
	log.Println("In-memory cache created")
	return &InMemoryOrderCache{
		byID:     make(map[string]*model.Order),
		byUserID: make(map[string][]*model.Order),
	}
}

func (c *InMemoryOrderCache) GetByID(id string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, ok := c.byID[id]
	if ok {
		log.Printf("Cache hit: order found by ID %s\n", id)
	} else {
		log.Printf("Cache miss: no order found by ID %s\n", id)
	}
	return order, ok
}

func (c *InMemoryOrderCache) ListByUser(userID string) ([]*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	orders, ok := c.byUserID[userID]
	if ok {
		log.Printf("Cache hit: found %d orders for user %s\n", len(orders), userID)
	} else {
		log.Printf("Cache miss: no orders found for user %s\n", userID)
	}
	return orders, ok
}

func (c *InMemoryOrderCache) Save(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byID[order.ID] = order
	c.byUserID[order.UserID] = append(c.byUserID[order.UserID], order)
	log.Printf("Cache save: stored order %s for user %s\n", order.ID, order.UserID)
}

func (c *InMemoryOrderCache) Update(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byID[order.ID] = order
	var filtered []*model.Order
	for _, o := range c.byUserID[order.UserID] {
		if o.ID != order.ID {
			filtered = append(filtered, o)
		}
	}
	filtered = append(filtered, order)
	c.byUserID[order.UserID] = filtered
	log.Printf("Cache update: updated order %s for user %s\n", order.ID, order.UserID)
}

func (c *InMemoryOrderCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if o, ok := c.byID[id]; ok {
		delete(c.byID, id)
		var filtered []*model.Order
		for _, order := range c.byUserID[o.UserID] {
			if order.ID != id {
				filtered = append(filtered, order)
			}
		}
		c.byUserID[o.UserID] = filtered
		log.Printf("Cache delete: removed order %s for user %s\n", id, o.UserID)
	}
}

func (c *InMemoryOrderCache) LoadFromDB(orders []*model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byID = make(map[string]*model.Order)
	c.byUserID = make(map[string][]*model.Order)
	for _, order := range orders {
		c.byID[order.ID] = order
		c.byUserID[order.UserID] = append(c.byUserID[order.UserID], order)
	}
	log.Printf("Cache initialized: loaded %d orders from the database\n", len(orders))
}
