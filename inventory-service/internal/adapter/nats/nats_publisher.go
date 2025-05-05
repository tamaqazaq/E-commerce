package natsadapter

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"inventory-service/internal/model"
	"time"
)

type NatsPublisher struct {
	nc *nats.Conn
}

func NewNatsPublisher(nc *nats.Conn) *NatsPublisher {
	return &NatsPublisher{nc: nc}
}

func (p *NatsPublisher) PublishProductCreated(product *model.Product) error {
	return p.publish("product.created", product)
}

func (p *NatsPublisher) PublishProductUpdated(product *model.Product) error {
	return p.publish("product.updated", product)
}

func (p *NatsPublisher) PublishProductDeleted(productID string) error {
	data := map[string]interface{}{
		"id": productID,
	}
	return p.publish("product.deleted", data)
}

func (p *NatsPublisher) publish(topic string, data interface{}) error {
	event := map[string]interface{}{
		"action": topic,
		"time":   time.Now().Format(time.RFC3339),
		"data":   data,
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.nc.Publish(topic, payload)
}
