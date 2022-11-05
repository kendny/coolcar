package trip

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
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
	POIManager     POIManager
	ProfileManager ProfileManager
	CarManager     CarManager
	Mongo          *dao.Mongo
	Logger         *zap.Logger
	//trippb.UnimplementedTripServiceServer // 必须引用，不然报错
}

// ProfileManager 防止入侵层
// ProfileManager defines the ACL (Anti Corruption Layer)
// for profile verification logic
type ProfileManager interface {
	// Verify 验证有没有租车的资质
	Verify(context.Context, id.AccountID) (id.IdentityID, error)
}

type CarManager interface {
	//Verify 验证车是否可以租用, *rentalpb.Location 确定人车的位置
	Verify(context.Context, id.CardID, *rentalpb.Location) error
	//Unlock 开锁
	Unlock(context.Context, id.CardID) error
}

// POIManager 查询坐标的能力
// POIManager resolves POI(Point Of Interest).
type POIManager interface {
	Resolve(context.Context, *rentalpb.Location) (string, error)
}

// CreateTrip creates a trip
func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	// 获取用户身份
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	// 1. 验证驾驶者身份
	iID, err := s.ProfileManager.Verify(c, aid)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	// 2. 检查车辆状态
	carID := id.CardID(req.CarId)
	err = s.CarManager.Verify(c, carID, req.Start)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	poi, err := s.POIManager.Resolve(c, req.Start)
	if err != nil {
		s.Logger.Info("cannot resolve poi", zap.Stringer("location", req.Start), zap.Error(err))
	}

	// 3. 创建行程：写入数据库，开始计费 (保证用户开锁后有行程)
	ls := &rentalpb.LocationStatus{
		Location: req.Start,
		PoiName:  poi,
	}
	tr, err := s.Mongo.CreateTrip(c, &rentalpb.Trip{
		AccountId:  aid.String(),
		CarId:      carID.String(),
		IdentityId: iID.String(),
		Status:     rentalpb.TripStatus_IN_PROGRESS,
		Start:      ls,
		Current:    ls,
	})
	// 创建行程失败
	if err != nil {
		s.Logger.Warn("cannot create trip", zap.Error(err))
		return nil, status.Error(codes.AlreadyExists, "")
	}

	// 4. 车辆开锁（无法开锁，需要有补救措施, 开锁是个复杂的过程，需要在后台进行开锁）
	go func() {
		err := s.CarManager.Unlock(context.Background(), carID)
		if err != nil {
			s.Logger.Error("cannot unlock car", zap.Error(err))
		}
	}()

	// 返回行程
	return &rentalpb.TripEntity{
		Id:   tr.ID.Hex(),
		Trip: tr.Trip,
	}, nil
}

// GetTrip gets a trip
func (s *Service) GetTrip(c context.Context, req *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// GetTrips get trips
func (s *Service) GetTrips(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
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
