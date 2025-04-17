package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"order-service/internal/model"
	"order-service/internal/usecase"
	pb "order-service/proto"
	"time"
)

type OrderGRPCServer struct {
	pb.UnimplementedOrderServiceServer
	usecase       usecase.OrderUsecase
	reviewUsecase usecase.ReviewUsecase
}

func NewOrderGRPCServer(orderUC usecase.OrderUsecase, reviewUC usecase.ReviewUsecase) *OrderGRPCServer {
	return &OrderGRPCServer{
		usecase:       orderUC,
		reviewUsecase: reviewUC,
	}
}

func (s *OrderGRPCServer) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	if req.Order == nil {
		return nil, status.Errorf(codes.InvalidArgument, "order field is required")
	}
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
func (s *OrderGRPCServer) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.ReviewResponse, error) {
	r := req.Review
	if r.Rating < 1.0 || r.Rating > 5.0 {
		return nil, status.Errorf(codes.InvalidArgument, "rating is invalid")
	}
	review := &model.Review{
		ProductID: r.ProductId,
		UserID:    r.UserId,
		Rating:    r.Rating,
		Comment:   r.Comment,
	}
	err := s.reviewUsecase.CreateReview(review)
	if err != nil {
		return nil, err
	}
	return &pb.ReviewResponse{Review: &pb.Review{
		Id:        review.ID,
		ProductId: review.ProductID,
		UserId:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt.Format(time.RFC3339),
		UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
	}}, nil
}

func (s *OrderGRPCServer) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.ReviewResponse, error) {
	r := req.Review
	if r.Rating < 1.0 || r.Rating > 5.0 {
		return nil, status.Errorf(codes.InvalidArgument, "rating is invalid")
	}
	review := &model.Review{
		ID:      r.Id,
		Rating:  r.Rating,
		Comment: r.Comment,
	}
	err := s.reviewUsecase.UpdateReview(review)
	if err != nil {
		return nil, err
	}
	review.UpdatedAt = time.Now()
	return &pb.ReviewResponse{Review: &pb.Review{
		Id:        review.ID,
		ProductId: review.ProductID,
		UserId:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
	}}, nil
}
