package car

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/share/id"
)

// 进行跟真实的设备交互，先全部返回成功
// todo... 逻辑待处理

//Manager defines a car manager
type Manager struct {
}

// Verify verifies car status
func (c *Manager) Verify(context.Context, id.CardID, *rentalpb.Location) error {
	return nil
}

// Unlock verifies a car
func (c *Manager) Unlock(context.Context, id.CardID) error {
	return nil
}
