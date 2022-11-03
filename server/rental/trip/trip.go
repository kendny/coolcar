package trip

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	trippb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/rental/trip/dao"
	"coolcar/server/share/auth"
	"coolcar/server/share/id"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

//所有实现必须嵌入UnimplementedTripServiceServer
//向前兼容
type Service struct {
	POIManager                            POIManager
	Mongo                                 *dao.Mongo
	Logger                                *zap.Logger
	trippb.UnimplementedTripServiceServer // 必须引用，不然报错
}

// POIManager resolves POI(Point Of Interest).
type POIManager interface {
	Resolve(context.Context, *rentalpb.Location) (string, error)
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
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, status.Error(codes.Unimplemented, "")
	}
	tid := id.TripID(req.Id)

	// 如果有同时updateTrip操作，这个地方就会产生脏数据， 如果 begin trans .... commit ... 就不会产生
	// begin trans
	tr, err := s.Mongo.GetTrip(c, tid, aid)

	cur := tr.Trip.Current.Location
	if req.Current != nil {
		cur = req.Current
	}

	if req.Current != nil {
		tr.Trip.Current = s.calcCurrentStatus(c, tr.Trip.Current, cur)
	}

	if req.EndTrip {
		tr.Trip.End = tr.Trip.Current
		tr.Trip.Status = rentalpb.TripStatus_FINISHED
	}
	// 更新
	err = s.Mongo.UpdateTrip(c, tid, aid, tr.UpdatedAt, tr.Trip)
	return nil, status.Error(codes.Unimplemented, "")
	//if err != nil {
	//	return nil, status.Error(codes.Aborted, "")
	//}
	//// commit
	//return tr.Trip, nil
}

var nowFunc = func() int64 {
	return time.Now().Unix()
}

const (
	centsPerSec = 0.7
	kmPerSec    = 0.02
)

// 计算当前状态(含费用)
func (s *Service) calcCurrentStatus(c context.Context, last *rentalpb.LocationStatus, cur *rentalpb.Location) *rentalpb.LocationStatus {
	//now := nowFunc()
	//elapsedSec := float64(now - last.TimestampSec)
	//poi, err := s.POIManager.Resolve(c, cur)
	//if err != nil {
	//	s.Logger.Info("cannot resolve poi", zap.Stringer("location", cur), zap.Error(err))
	//}
	//return &rentalpb.LocationStatus{
	//	Location:     cur,
	//	FeeCent:      last.FeeCent + int32(centsPerSec*elapsedSec*2*rand.Float64()),
	//	KmDriven:     last.KmDriven + kmPerSec*elapsedSec*2*rand.Float64(),
	//	TimestampSec: now,
	//	PoiName:      poi,
	//}
	return nil
}
