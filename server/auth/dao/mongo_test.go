package dao

import (
	"context"
	"coolcar/server/share/id"
	mgutil "coolcar/server/share/mongo"
	"coolcar/server/share/mongo/objid"
	mongotesting "coolcar/server/share/testing"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"testing"
)

var mongoURI string

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	//"mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"
	//mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	mc, err := mongotesting.NewDefaultClient(c)
	if err != nil {
		t.Fatalf("cannt connect mongodb: %v\n", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	// 构造多个测试实例
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("5f7c245ab0361e00ffb9fd6f")),
			openIDField:        "openid_1",
		},
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id.AccountID("5f7c245ab0361e00ffb9fd70")),
			openIDField:        "openid_2",
		},
	})
	if err != nil {
		t.Fatalf("cannot insert initial values: %v\n", err)
	}
	// 生成固定ID
	mgutil.NewObjID = func() primitive.ObjectID {
		return mustObjID("6352921b22327fd5955ae001")
	}

	//表格驱动测试
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "5f7c245ab0361e00ffb9fd6f",
		}, {
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "5f7c245ab0361e00ffb9fd70",
		}, {
			name:   "new_user",
			openID: "openid_3",
			want:   "6352921b22327fd5955ae001",
		},
	}

	for i, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			id, err := m.ResolveAccountID(context.Background(), cc.openID)
			fmt.Printf("i: %v, c: %v, id: %v, %v\n", i, cc, id, err)
			if err != nil {
				t.Errorf("failed resolve account id for %q:%v\n", cc.openID, err)
			}
			if id.String() != cc.want {
				t.Errorf("resolve account id: want:%q, got:%q\n", cc.want, id)
			}
		})
	}

	//id, err := m.ResolveAccountID(c, "aaaa")
	//if err != nil {
	//	t.Errorf("faild resolve account id for aaaa: %v\n", err)
	//} else {
	//	want := "6352921b22327fd5955ae001"
	//	if id != want {
	//		t.Errorf("resolve account id: want:%q, got:%q\n", want, id)
	//	}
	//}
}

func mustObjID(hex string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return objID
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m)) // 确保测试每次在新的docker环境中运行
}
