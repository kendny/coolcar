package wechat

import (
	"fmt"
	weapp "github.com/medivhzhan/weapp/v2"
	"go.uber.org/zap"
)

type Service struct {
	AppID     string
	AppSecret string
	Logger    *zap.Logger // 服务类型一般都用指针
}

func (s *Service) Resolve(code string) (string, error) {
	resp, err := weapp.Login(s.AppID, s.AppSecret, code)
	if err != nil {
		return "", fmt.Errorf("weapp.Login failed: %v", err)
	}
	if resp.GetResponseError() != nil {
		return "", fmt.Errorf("weapp response error: %v", err)
	}
	fmt.Printf("LoginResponse: %v\n", resp)
	// 微信没邦开放平台， UnionID为空
	s.Logger.Info("received code", zap.String("OpenID: ", resp.OpenID), zap.String("UnionID: ", resp.UnionID))
	return resp.OpenID, nil
}
