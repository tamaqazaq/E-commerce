package handler

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"statistics-service/internal/usecase"
)

type OrderCreatedEvent struct {
	UserID    string `json:"user_id"`
	OrderTime string `json:"order_time"`
}

func NewNATSHandler(nc *nats.Conn, service usecase.StatisticsService) {
	_, err := nc.Subscribe("order.created", func(msg *nats.Msg) {
		var event OrderCreatedEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			fmt.Println("Failed to parse NATS event:", err)
			return
		}
		err := service.SaveUserOrder(event.UserID, event.OrderTime)
		if err != nil {
			fmt.Println("Failed to save order stat:", err)
			return
		}
		fmt.Println("Order stat saved for user", event.UserID)
	})

	if err != nil {
		fmt.Println("NATS subscription failed:", err)
	}
}
