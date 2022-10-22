package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"))

	if err != nil {
		t.Fatalf("cannt connect mongodb: %v\n", err)
	}

	m := NewMongo(mc.Database("coolcar"))
	id, err := m.ResolveAccountID(c, "aaaa")
	if err != nil {
		t.Errorf("faild resolve account id for aaaa: %v\n", err)
	} else {
		want := "6352921b22327fd5955ae001"
		if id != want {
			t.Errorf("resolve account id: want:%q, got:%q\n", want, id)
		}
	}
}
