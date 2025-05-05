package natsadapter

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"order-service/internal/model"
	"time"
)

type NatsPublisher struct {
	nc *nats.Conn
}

func NewNatsPublisher(nc *nats.Conn) *NatsPublisher {
	return &NatsPublisher{nc: nc}
}

func (p *NatsPublisher) PublishOrderCreated(order *model.Order) error {
	event := map[string]interface{}{
		"action": "order.created",
		"time":   time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"order_id":   order.ID,
			"user_id":    order.UserID,
			"total":      order.Total,
			"order_time": order.Timestamp.Format(time.RFC3339),
		},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Error marshaling order.created event:", err)
		return err
	}
	return p.nc.Publish("order.created", payload)
}
func (p *NatsPublisher) PublishOrderUpdated(order *model.Order) error {
	event := map[string]interface{}{
		"action": "order.updated",
		"time":   time.Now().Format(time.RFC3339),
		"data": map[string]interface{}{
			"order_id":   order.ID,
			"user_id":    order.UserID,
			"total":      order.Total,
			"status":     order.Status,
			"order_time": order.Timestamp.Format(time.RFC3339),
		},
	}
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.nc.Publish("order.updated", payload)
}
