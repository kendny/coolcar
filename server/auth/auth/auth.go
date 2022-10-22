package auth

import (
	"context"
	authpb "coolcar/server/auth/api/gen/v1"
	"coolcar/server/auth/dao"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements auth service
type Service struct {
	OpenIDResolver OpenIDResolver
	Mongo          *dao.Mongo
	Logger         *zap.Logger // 服务类型一般都用指针
}

// OpenIDResolver resolves an authorization code to an open id.
type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

// Login logs a user in
func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (resp *authpb.LoginResponse, err error) {
	// 获取用户的openid
	openID, err := s.OpenIDResolver.Resolve(req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openid: %v", err)
	}
	s.Logger.Info("received code", zap.String("code", req.Code), zap.String("code", req.GetCode()))

	//拿到openid 到数据库 查询或者保存用户
	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil {
		return nil, status.Error(codes.Internal, "")
	}

	return &authpb.LoginResponse{
		AccessToken: req.Code + "-" + openID + "-, accountID:" + accountID,
		ExpiresIn:   7200,
	}, nil
}
