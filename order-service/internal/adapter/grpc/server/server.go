package server

import (
	"context"
	"github.com/google/uuid"
	"order-service/internal/model"
	"order-service/internal/usecase"
	pb "order-service/proto"
)

type OrderGRPCServer struct {
	pb.UnimplementedOrderServiceServer
	usecase usecase.OrderUsecase
}

func NewOrderGRPCServer(uc usecase.OrderUsecase) *OrderGRPCServer {
	return &OrderGRPCServer{usecase: uc}
}

func (s *OrderGRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	order := &model.Order{
		ID:     uuid.New().String(),
		UserID: req.Order.UserId,
		Status: "pending",
	}

	for _, item := range req.Order.Items {
		order.Items = append(order.Items, model.OrderItem{
			ID:        uuid.New().String(),
			OrderID:   order.ID,
			ProductID: item.ProductId,
			Quantity:  int(item.Quantity),
			Price:     item.Price,
		})
	}

	total := 0.0
	for _, item := range order.Items {
		total += item.Price * float64(item.Quantity)
	}
	order.Total = total

	if err := s.usecase.Create(order); err != nil {
		return nil, err
	}

	return &pb.OrderResponse{Order: &pb.Order{
		Id:     order.ID,
		UserId: order.UserID,
		Status: order.Status,
		Total:  order.Total,
		Items:  convertOrderItemsToPB(order.Items),
	}}, nil
}

func (s *OrderGRPCServer) GetOrderByID(ctx context.Context, req *pb.GetOrderRequest) (*pb.OrderResponse, error) {
	order, err := s.usecase.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{Order: &pb.Order{
		Id:     order.ID,
		UserId: order.UserID,
		Status: order.Status,
		Total:  order.Total,
		Items:  convertOrderItemsToPB(order.Items),
	}}, nil
}

func (s *OrderGRPCServer) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.OrderResponse, error) {
	err := s.usecase.UpdateStatus(req.Id, req.Status)
	if err != nil {
		return nil, err
	}
	order, _ := s.usecase.GetByID(req.Id)
	return &pb.OrderResponse{Order: &pb.Order{
		Id:     order.ID,
		UserId: order.UserID,
		Status: order.Status,
		Total:  order.Total,
		Items:  convertOrderItemsToPB(order.Items),
	}}, nil
}

func (s *OrderGRPCServer) ListUserOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.usecase.ListByUser(req.UserId)
	if err != nil {
		return nil, err
	}
	var pbOrders []*pb.Order
	for _, order := range orders {
		pbOrders = append(pbOrders, &pb.Order{
			Id:     order.ID,
			UserId: order.UserID,
			Status: order.Status,
			Total:  order.Total,
			Items:  convertOrderItemsToPB(order.Items),
		})
	}
	return &pb.ListOrdersResponse{Orders: pbOrders}, nil
}

func convertOrderItemsToPB(items []model.OrderItem) []*pb.OrderItem {
	var pbItems []*pb.OrderItem
	for _, item := range items {
		pbItems = append(pbItems, &pb.OrderItem{
			Id:        item.ID,
			OrderId:   item.OrderID,
			ProductId: item.ProductID,
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		})
	}
	return pbItems
}
