package dao

import (
	"context"
	"coolcar/server/share/id"
	mgutil "coolcar/server/share/mongo"
	"coolcar/server/share/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const openIDField = "open_id"

// Mongo defines a mongo dao.
type Mongo struct {
	col      *mongo.Collection
	newObjID func() primitive.ObjectID // 生成ID的函数
}

// NewMongo creates a new mongo dao.
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:      db.Collection("account"),
		newObjID: primitive.NewObjectID,
	}
}

// ResolveAccountID resolves an account id from open id
func (m *Mongo) ResolveAccountID(c context.Context, openID string) (id.AccountID, error) {

	insertedID := mgutil.NewObjID() // m.newObjID()
	res := m.col.FindOneAndUpdate(c, bson.M{
		openIDField: openID,
	}, mgutil.SetOnInsert(bson.M{
		mgutil.IDFieldName: insertedID,
		openIDField:        openID,
	}), options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After))

	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate: %v\n", err)
	}

	var row mgutil.ObjID

	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot decode result: %v\n", err)
	}

	return objid.ToAccountID(row.ID), nil
}
