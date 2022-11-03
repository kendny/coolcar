package dao

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/share/id"
	mgutil "coolcar/server/share/mongo"
	"coolcar/server/share/mongo/objid"
	mongotesting "coolcar/server/share/testing"
	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
	"os"
	"testing"
)

var mongoURI string

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	//"mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannt connect mongodb: %v\n", err)
	}

	db := mc.Database("coolcar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannt setup indexes: %v\n", err)
	}
	m := NewMongo(db)

	// 表格
	cases := []struct {
		name       string
		tripID     string
		accountID  string
		tripStatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     "5f8132eb00714bf62948905c",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "another_finished",
			tripID:     "5f8132eb00714bf62948905d",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "in_progress",
			tripID:     "5f8132eb00714bf62948905e",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			name:       "another_in_progress",
			tripID:     "5f8132eb00714bf62948905f",
			accountID:  "account1",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		},
		{
			name:       "in_progress_by_another_account",
			tripID:     "5f8132eb00714bf629489060",
			accountID:  "account2",
			tripStatus: rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, cc := range cases {
		//mgutil.NewObjID = func() primitive.ObjectID {
		//	return objid.MustFromID(id.TripID(cc.tripID))
		//}
		mgutil.NewObjIDWithValue(id.TripID(cc.tripID))
		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: cc.accountID,
			Status:    cc.tripStatus,
		})
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s error expected; got none", cc.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s error creating trip: %v\n", cc.name, err)
			continue
		}
		if tr.ID.Hex() != cc.tripID {
			t.Errorf("%s incorrect trip id; want %q; got %q", cc.name, tr.ID.Hex(), cc.tripID)
		}
	}
}

func TestGetTrip(t *testing.T) {
	c := context.Background()
	//"mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannt connect mongodb: %v\n", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	acct := id.AccountID("account1")
	// TODO... 设置回去才不会挂
	mgutil.NewObjID = primitive.NewObjectID
	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: acct.String(),
		CarId:     "car1",
		Start: &rentalpb.LocationStatus{
			PoiName: "startpoint",
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 120,
			},
		},
		End: &rentalpb.LocationStatus{
			PoiName:  "endpoint",
			FeeCent:  10000,
			KmDriven: 35,
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 115,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})
	if err != nil {
		t.Errorf("cannot create trip: %v\n", err)
	}

	//	 测试 GetTrip
	got, err := m.GetTrip(c, objid.ToTripID(tr.ID), acct)
	if err != nil {
		t.Errorf("cannot get trip: %v\n", err)
	}
	// TODO... 进行破坏测试
	//got.Trip.Start.PoiName = "badsstart"
	// TODO ... protocmp.Transform() ???
	if diff := cmp.Diff(tr, got, protocmp.Transform()); diff != "" {
		t.Errorf("result differs; -want +got: %s", diff)
	}
}

func TestGetTrips(t *testing.T) {
	// 建客户端
	c := context.Background()
	//"mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannt connect mongodb: %v\n", err)
	}

	m := NewMongo(mc.Database("coolcar"))

	rows := []struct {
		id        string
		accountID string
		status    rentalpb.TripStatus
	}{
		{
			id:        "5f8132eb10714bf629489051",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "5f8132eb10714bf629489052",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "5f8132eb10714bf629489053",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "5f8132eb10714bf629489054",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			id:        "5f8132eb10714bf629489055",
			accountID: "account_id_for_get_trips_1",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, r := range rows {
		mgutil.NewObjIDWithValue(id.TripID(r.id))
		_, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: r.accountID,
			Status:    r.status,
		})

		if err != nil {
			t.Fatalf("cannot create rows: %v\n", err)
		}
	}

	cases := []struct {
		name       string
		accountID  string
		status     rentalpb.TripStatus
		wantCount  int
		wantOnlyID string
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},
		{
			name:       "get_in_progress",
			accountID:  "account_id_for_get_trips",
			status:     rentalpb.TripStatus_IN_PROGRESS,
			wantCount:  1,
			wantOnlyID: "5f8132eb10714bf629489054",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			res, err := m.GetTrips(context.Background(), id.AccountID(cc.accountID), cc.status)
			if err != nil {
				t.Errorf("cannot get trips: %v\n", err)
			}

			if cc.wantCount != len(res) {
				t.Errorf("incorrect result count; want %d, got:%d\n", cc.wantCount, len(res))
			}

			if cc.wantOnlyID != "" && len(res) > 0 {
				if cc.wantOnlyID != res[0].ID.Hex() {
					t.Errorf("only_id incorrect; want %q, got:%q\n", cc.wantOnlyID, res[0].ID.Hex())
				}
			}
		})
	}

}

func TestUpdateTrip(t *testing.T) {

	c := context.Background()
	//"mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannt connect mongodb: %v\n", err)
	}

	m := NewMongo(mc.Database("coolcar"))

	//	一条记录两个人同一时刻更改
	tid := id.TripID("5f8132eb12714bf629489054")
	aid := id.AccountID("account_for_update")

	// 设定一个固定的时间
	var now int64 = 10000
	mgutil.NewObjIDWithValue(tid)
	mgutil.UpdatedAt = func() int64 {
		return now
	}
	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi",
		},
	})
	if err != nil {
		t.Fatalf("cannot create trip: %v\n", err)
	}
	if tr.UpdatedAt != 10000 {
		t.Fatalf("wrong updatedat; want:10000, got:%d\n", tr.UpdatedAt)
	}

	//更新的值
	update := &rentalpb.Trip{
		AccountId: aid.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PoiName: "start_poi_updated",
		},
	}

	cases := []struct {
		name          string
		now           int64 // 当前更新的时间
		withUpdatedAt int64 // 更新记录的时间
		wantErr       bool
	}{
		{
			name:          "normal_update",
			now:           20000,
			withUpdatedAt: 10000,
		}, {
			name:          "update_with_stale_timestamp",
			now:           30000,
			withUpdatedAt: 10000,
			wantErr:       true,
		}, {
			name:          "update_with_refetch",
			now:           40000,
			withUpdatedAt: 20000,
		},
	}

	for _, cc := range cases {
		// 控制时间的走动
		now = cc.now
		err := m.UpdateTrip(c, tid, aid, cc.withUpdatedAt, update)
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s: want error; got none", cc.name)
			} else {
				// 期待出错，出错了继续
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s:cannot update:%v\n", cc.name, err)
			}
		}

		// 校验时间戳
		updatedTrip, err := m.GetTrip(c, tid, aid)
		if err != nil {
			t.Errorf("%s:cannot get trip after update: %v\n", cc.name, err)
		}
		if cc.now != updatedTrip.UpdatedAt {
			t.Errorf("%s: incorrect updatedat: want %d, got %d\n", cc.name, cc.now, updatedTrip.UpdatedAt)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m)) // 确保测试每次在新的docker环境中运行
}
