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

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m)) // 确保测试每次在新的docker环境中运行
}
