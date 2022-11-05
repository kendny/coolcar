package trip

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/rental/trip/client/poi"
	"coolcar/server/rental/trip/dao"
	"coolcar/server/share/auth"
	"coolcar/server/share/id"
	mgutil "coolcar/server/share/mongo"
	"coolcar/server/share/server"
	mongotesting "coolcar/server/share/testing"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestCreateTrip(t *testing.T) {
	//c := context.Background()
	c := auth.ContextWithAccountID(context.Background(), id.AccountID("account1"))
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot create mongo client: %v", err)
	}

	logger, err := server.NewZapLogger()
	if err != nil {
		t.Fatalf("cannot create logger: %v\n", err)
	}

	pm := &profileManager{}
	cm := &carManager{}
	s := &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		Mongo:          dao.NewMongo(mc.Database("coolcar")),
		Logger:         logger,
	}

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.123,
			Longitude: 114.2523,
		},
	}

	pm.iID = "identity1"
	golden := `{"account_id":"account1","car_id":"car1","start":{"location":{"latitude":32.123,"longitude":114.2523},"poi_name":"中关村"},"current":{"location":{"latitude":32.123,"longitude":114.2523},"poi_name":"中关村"},"status":1,"identity_id":"identity1"}`
	//	 设计表格案例
	cases := []struct {
		name         string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:   "normal_create",
			tripID: "5f8132eb12714bf629489054",
			want:   golden,
		}, {
			name:       "profile_err",
			tripID:     "5f8132eb12714bf629489055",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		}, {
			name:         "car_verify_err",
			tripID:       "5f8132eb12714bf629489056",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		}, {
			name:         "car_unlock_err",
			tripID:       "5f8132eb12714bf629489057",
			carUnlockErr: fmt.Errorf("unlock"),
			want:         golden,
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			mgutil.NewObjIDWithValue(id.TripID(cc.tripID)) // cc.tripID =
			pm.err = cc.profileErr
			cm.unlockErr = cc.carUnlockErr
			cm.verifyErr = cc.carVerifyErr
			// 创建行程
			// todo... 为啥用 context.Background() 就会报 panic: runtime error:
			// todo... invalid memory address or nil pointer dereference？
			res, err := s.CreateTrip(c, req)
			if cc.wantErr {
				if err == nil {
					t.Errorf("want error; go none\n")
				} else {
					return
				}
			}

			if err != nil {
				t.Errorf("error creating trip: %v\n", err)
			}
			if res.Id != cc.tripID {
				t.Errorf("incorrect id; want: %q, got: %q\n", cc.tripID, res.Id)
			}
			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("cannot marshall response:%v\n", err)
			}
			got := string(b)
			if cc.want != got {
				t.Errorf("incorrect response:want %s, got: %s\n", cc.want, got)
			}

		})
	}
}

type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(context.Context, id.AccountID) (id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyErr error
	unlockErr error
}

func (c *carManager) Verify(context.Context, id.CardID, *rentalpb.Location) error {
	return c.verifyErr
}

func (c *carManager) Unlock(context.Context, id.CardID) error {
	return c.unlockErr
}

func TestMain(m *testing.M) {
	// 测试跑在真实的数据库环境
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
