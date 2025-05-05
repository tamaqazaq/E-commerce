package server

import (
	"context"
	"statistics-service/internal/app/service"
	pb "statistics-service/proto"
)

type StatisticsGRPCServer struct {
	pb.UnimplementedStatisticsServiceServer
	service *service.StatisticsService
}

func NewStatisticsGRPCServer(svc *service.StatisticsService) *StatisticsGRPCServer {
	return &StatisticsGRPCServer{service: svc}
}

func (s *StatisticsGRPCServer) GetUserOrdersStatistics(ctx context.Context, req *pb.UserOrderStatisticsRequest) (*pb.UserOrderStatisticsResponse, error) {
	total, hourly, err := s.service.GetUserOrderStats(req.UserId)
	if err != nil {
		return nil, err
	}

	var result []*pb.HourCount
	for hour, count := range hourly {
		result = append(result, &pb.HourCount{
			Hour:  int32(hour),
			Count: int32(count),
		})
	}

	return &pb.UserOrderStatisticsResponse{
		TotalOrders:  int32(total),
		OrdersByHour: result,
	}, nil
}

func (s *StatisticsGRPCServer) GetUserStatistics(ctx context.Context, _ *pb.UserStatisticsRequest) (*pb.UserStatisticsResponse, error) {
	users, products, err := s.service.GetGeneralStats()
	if err != nil {
		return nil, err
	}
	return &pb.UserStatisticsResponse{
		TotalUsers:    int32(users),
		TotalProducts: int32(products),
	}, nil
}
