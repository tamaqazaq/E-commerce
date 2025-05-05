package natslistener

import (
	"encoding/json"
	"log"
	"statistics-service/internal/model"
	"statistics-service/internal/usecase"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

func SubscribeToEvents(nc *nats.Conn, repo usecase.StatisticsRepository) {
	subscribe(nc, "order.created", func(data []byte) {
		var rawEvent struct {
			UserID    string  `json:"user_id"`
			OrderID   string  `json:"order_id"`
			Total     float64 `json:"total"`
			OrderTime string  `json:"order_time"`
		}

		if err := json.Unmarshal(data, &rawEvent); err != nil {
			log.Println("Failed to parse order.created event:", err)
			return
		}

		if rawEvent.OrderTime == "" {
			return
		}

		timestamp, err := time.Parse(time.RFC3339, rawEvent.OrderTime)
		if err != nil {
			log.Println("Failed to parse order_time format:", rawEvent.OrderTime, err)
			return
		}

		event := model.OrderEvent{
			UserID:    rawEvent.UserID,
			OrderID:   rawEvent.OrderID,
			Total:     rawEvent.Total,
			Timestamp: timestamp,
			Action:    "created",
		}

		log.Printf("Saving order.CREATED event: %+v\n", event)
		if err := repo.SaveOrderEvent(&event); err != nil {
			log.Println("SaveOrderEvent error:", err)
		}
	})

	subscribe(nc, "order.updated", func(data []byte) {
		var raw map[string]interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			log.Println("Failed to parse order.updated event:", err)
			return
		}

		dataMap, ok := raw["data"].(map[string]interface{})
		if !ok {
			log.Println("Missing data field in order.updated")
			return
		}

		timestampStr, _ := dataMap["order_time"].(string)
		t, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			log.Println("Invalid timestamp in order.updated:", err)
			return
		}

		event := model.OrderEvent{
			OrderID:   getString(dataMap["order_id"]),
			UserID:    getString(dataMap["user_id"]),
			Total:     getFloat(dataMap["total"]),
			Status:    getString(dataMap["status"]),
			Timestamp: t,
			Action:    "updated",
		}

		log.Printf("Saving order.UPDATED event: %+v\n", event)
		if err := repo.SaveOrderEvent(&event); err != nil {
			log.Println("SaveOrderEvent error:", err)
		}
	})

	subscribe(nc, "product.created", handleProductEvent(repo, "created"))
	subscribe(nc, "product.updated", handleProductEvent(repo, "updated"))
	subscribe(nc, "product.deleted", handleProductEvent(repo, "deleted"))
}

func subscribe(nc *nats.Conn, topic string, handler func([]byte)) {
	_, err := nc.Subscribe(topic, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	if err != nil {
		log.Printf("Failed to subscribe to %s: %v", topic, err)
	}
}

func handleProductEvent(repo usecase.StatisticsRepository, action string) func([]byte) {
	return func(data []byte) {
		var raw map[string]interface{}
		if err := json.Unmarshal(data, &raw); err != nil {
			log.Printf("Failed to parse product.%s event: %v", strings.ToUpper(action), err)
			return
		}

		dataMap, ok := raw["data"].(map[string]interface{})
		if !ok {
			log.Printf("Invalid 'data' field in product.%s event", strings.ToUpper(action))
			return
		}

		event := model.ProductEvent{
			ProductID: getString(dataMap["id"]),
			Name:      getString(dataMap["name"]),
			Category:  getString(dataMap["category"]),
			Price:     getFloat(dataMap["price"]),
			Stock:     getInt(dataMap["stock"]),
			Action:    action,
			Timestamp: time.Now(),
		}

		log.Printf("Saving product.%s event: %+v", strings.ToUpper(action), event)
		if err := repo.SaveProductEvent(&event); err != nil {
			log.Printf("Failed to save product.%s event: %v", strings.ToUpper(action), err)
		}
	}
}

func getString(v interface{}) string {
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func getFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

func getInt(v interface{}) int {
	if f, ok := v.(float64); ok {
		return int(f)
	}
	return 0
}
