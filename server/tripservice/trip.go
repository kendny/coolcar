/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-19 15:58:30
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-06-19 16:51:35
 * @FilePath: /wx/Users/xxxian/go_project/src/coolcar/server/tripservice/trip.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package trip

import (
	"context"
	trippb "coolcar/server/proto/gen/go"
)

// Service implements trip service
type Service struct {
	trippb.UnimplementedTripServiceServer // 必须引用，不然报错
}

func (*Service) GetTrip(c context.Context, req *trippb.GetTripRequest) (*trippb.GetTripResponse, error) {
	return &trippb.GetTripResponse{
		Id: req.Id,
		Trip: &trippb.Trip{
			Start:       "abc",
			End:         "def",
			DurationSec: 3600,
			FeeCent:     10000,
			StartPos: &trippb.Location{
				Latitude:  34,
				Longitude: 45,
			},
			EndPos: &trippb.Location{
				Latitude:  10,
				Longitude: 32,
			},
			PathLocationns: []*trippb.Location{
				{
					Latitude:  1,
					Longitude: 12,
				}, {
					Latitude:  91,
					Longitude: 32,
				},
			},
			Status: trippb.TripStatus_FINISHED,
		},
	}, nil

}
