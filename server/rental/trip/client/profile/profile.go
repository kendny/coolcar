package profile

import (
	"context"
	"coolcar/server/share/id"
)

type Manager struct {
}

func (c *Manager) Verify(context.Context, id.AccountID) (id.IdentityID, error) {
	return id.IdentityID("identity1"), nil
}
