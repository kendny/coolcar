package profile

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/rental/profile/dao"
	"coolcar/server/share/auth"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// Service defines a profile service
type Service struct {
	Mongo  *dao.Mongo
	Logger *zap.Logger
}

// GetProfie gets profile for the current account
func (s *Service) GetProfile(c context.Context, req *rentalpb.GetProfileRequest) (*rentalpb.Profile, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("GetProfie aid:", zap.String("aid:", aid.String()))
	p, err := s.Mongo.GetProfile(c, aid)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &rentalpb.Profile{}, nil
		}
		s.Logger.Error("cannot get profile:", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	return p, nil
}

// SubmitProfile submits a profile.
func (s *Service) SubmitProfile(c context.Context, i *rentalpb.Identity) (*rentalpb.Profile, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}

	p := &rentalpb.Profile{
		Identity:       i,
		IdentityStatus: rentalpb.IdentityStatus_PENDING,
	}
	// 对前置条件的处理
	err = s.Mongo.UpdateProfile(c, aid, rentalpb.IdentityStatus_UNSUBMITTED, p)
	if err != nil {
		s.Logger.Error("cannot update profile:", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}

	// 模拟审核通过
	go func() {
		time.Sleep(3 * time.Second)
		err := s.Mongo.UpdateProfile(context.Background(), aid, rentalpb.IdentityStatus_PENDING, &rentalpb.Profile{
			Identity:       i,
			IdentityStatus: rentalpb.IdentityStatus_VERIFIED,
		})
		if err != nil {
			s.Logger.Error("cannot verify identity: ", zap.Error(err))
		}
	}()
	return p, nil
}

func (s *Service) ClearProfile(c context.Context, req *rentalpb.ClearProfileRequest) (*rentalpb.Profile, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}

	p := &rentalpb.Profile{}
	err = s.Mongo.UpdateProfile(c, aid, rentalpb.IdentityStatus_VERIFIED, p)
	if err != nil {
		s.Logger.Error("cannot clear profile:", zap.Error(err))
		return nil, status.Error(codes.Internal, "")
	}
	return p, nil
}
