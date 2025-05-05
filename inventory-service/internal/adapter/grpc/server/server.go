package server

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"inventory-service/internal/model"
	"inventory-service/internal/usecase"
	pb "inventory-service/proto"
)

type InventoryGRPCServer struct {
	pb.UnimplementedInventoryServiceServer
	usecase usecase.ProductUsecase
}

func NewInventoryGRPCServer(uc usecase.ProductUsecase) *InventoryGRPCServer {
	return &InventoryGRPCServer{usecase: uc}
}

func (s *InventoryGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	if req.Product == nil {
		return nil, status.Errorf(codes.InvalidArgument, "product field is required")
	}
	productID := uuid.New().String()

	p := &model.Product{
		ID:       productID,
		Name:     req.Product.Name,
		Category: req.Product.Category,
		Price:    req.Product.Price,
		Stock:    int(req.Product.Stock),
	}

	if err := s.usecase.Create(p); err != nil {
		return nil, err
	}

	return &pb.ProductResponse{Product: &pb.Product{
		Id:       p.ID,
		Name:     p.Name,
		Category: p.Category,
		Price:    p.Price,
		Stock:    int32(p.Stock),
	}}, nil
}

func (s *InventoryGRPCServer) GetProductByID(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	p, err := s.usecase.GetByID(req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.ProductResponse{Product: &pb.Product{
		Id:       p.ID,
		Name:     p.Name,
		Category: p.Category,
		Price:    p.Price,
		Stock:    int32(p.Stock),
	}}, nil
}

func (s *InventoryGRPCServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	p := &model.Product{
		ID:       req.Product.Id,
		Name:     req.Product.Name,
		Category: req.Product.Category,
		Price:    req.Product.Price,
		Stock:    int(req.Product.Stock),
	}
	if err := s.usecase.Update(p.ID, p); err != nil {
		return nil, err
	}
	return &pb.ProductResponse{Product: req.Product}, nil
}

func (s *InventoryGRPCServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
	if err := s.usecase.Delete(req.Id); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *InventoryGRPCServer) ListProducts(ctx context.Context, _ *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	list, err := s.usecase.List()
	if err != nil {
		return nil, err
	}
	var protoList []*pb.Product
	for _, p := range list {
		protoList = append(protoList, &pb.Product{
			Id:       p.ID,
			Name:     p.Name,
			Category: p.Category,
			Price:    p.Price,
			Stock:    int32(p.Stock),
		})
	}
	return &pb.ListProductsResponse{Products: protoList}, nil
}
