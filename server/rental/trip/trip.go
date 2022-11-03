package trip

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	trippb "coolcar/server/rental/api/gen/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//所有实现必须嵌入UnimplementedTripServiceServer
//向前兼容
type Service struct {
	Logger                                *zap.Logger
	trippb.UnimplementedTripServiceServer // 必须引用，不然报错
}

// CreateTrip creates a trip
func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// GetTrip gets a trip
func (s *Service) GetTrip(c context.Context, req *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// GetTrips get trips
func (s *Service) GetTrips(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsRequest, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// UpdateTrips updates a trip.
func (s *Service) UpdateTrip(c context.Context, req *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	return nil, status.Error(codes.Unimplemented, "")
}
