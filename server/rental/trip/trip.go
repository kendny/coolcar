package trip

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	trippb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/share/auth"
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
func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.CreateTripResponse, error) {
	// get accountID from context
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("create trip", zap.String("start", req.Start), zap.String("account_id", aid))
	return nil, status.Error(codes.Unimplemented, "")
}
