package service

import (
	"api-gateway/proto"
	"google.golang.org/grpc"
	"log"
)

type GRPCClient struct {
	InventoryClient proto.InventoryServiceClient
	OrderClient     proto.OrderServiceClient
}

func NewGRPCClient(inventoryAddr, orderAddr string) (*GRPCClient, error) {
	inventoryConn, err := grpc.Dial(inventoryAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to inventory service: %v", err)
		return nil, err
	}

	orderConn, err := grpc.Dial(orderAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to order service: %v", err)
		return nil, err
	}

	client := &GRPCClient{
		InventoryClient: proto.NewInventoryServiceClient(inventoryConn),
		OrderClient:     proto.NewOrderServiceClient(orderConn),
	}

	return client, nil
}
