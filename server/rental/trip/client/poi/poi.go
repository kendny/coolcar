package poi

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"github.com/golang/protobuf/proto"
	"hash/fnv"
)

/*
client 文件夹是完成和其他系统交互的逻辑
*/
// 构造假地标
var poi = []string{
	"中关村",
	"天安门",
	"陆家嘴",
	"迪士尼",
	"天河体育中心",
	"广州塔",
}

// Manager defines a poi manager
type Manager struct {
}

// Resolve resolves the given location
func (*Manager) Resolve(c context.Context, loc *rentalpb.Location) (string, error) {
	b, err := proto.Marshal(loc)
	if err != nil {
		return "", nil
	}
	h := fnv.New32()
	h.Write(b)

	// h.Sum32() 获取哈希
	return poi[int(h.Sum32())%len(poi)], nil
}
