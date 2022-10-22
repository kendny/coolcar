package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 将mongo中的操作进行提取共享

// IDField defineds the field name for mongo  document id.
const IDField = "_id"

// ObjID defines the object id field
type ObjID struct {
	ID primitive.ObjectID `bson:"_id"`
}

// Set returns a $set update document
func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}
