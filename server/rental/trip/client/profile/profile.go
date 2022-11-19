package profile

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/share/id"
	base64 "encoding/base64"
	"fmt"
	"google.golang.org/protobuf/proto"
)

// Fetcher defines the interface to fetch profile.
type Fetcher interface {
	GetProfile(c context.Context, req *rentalpb.GetProfileRequest) (*rentalpb.Profile, error)
}

type Manager struct {
	Fetcher Fetcher
}

// Verify verifies account identity.
func (m *Manager) Verify(c context.Context, aid id.AccountID) (id.IdentityID, error) {
	nilID := id.IdentityID("")
	p, err := m.Fetcher.GetProfile(c, &rentalpb.GetProfileRequest{})
	if err != nil {
		return nilID, fmt.Errorf("cannot get profile: %v", err)
	}

	if p.IdentityStatus != rentalpb.IdentityStatus_VERIFIED {
		return nilID, fmt.Errorf("invalid identity status")
	}

	// 将 Identity 其序列化成二进制数据
	b, err := proto.Marshal(p.Identity)
	if err != nil {
		return nilID, fmt.Errorf("cannot marshal identity: %v", err)
	}

	// 将 []byte 转化为 base64 进行传输
	return id.IdentityID(base64.StdEncoding.EncodeToString(b)), nil
}
