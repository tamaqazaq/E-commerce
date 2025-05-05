package service

import (
	"api-gateway/proto"
	"google.golang.org/grpc"
	"log"
)

type GRPCClient struct {
	InventoryClient  proto.InventoryServiceClient
	OrderClient      proto.OrderServiceClient
	StatisticsClient proto.StatisticsServiceClient
}

func NewGRPCClient(inventoryAddr, orderAddr, statsAddr string) (*GRPCClient, error) {
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
	statsConn, err := grpc.Dial(statsAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to statistics service: %v", err)
		return nil, err
	}

	client := &GRPCClient{
		InventoryClient:  proto.NewInventoryServiceClient(inventoryConn),
		OrderClient:      proto.NewOrderServiceClient(orderConn),
		StatisticsClient: proto.NewStatisticsServiceClient(statsConn),
	}

	return client, nil
}
