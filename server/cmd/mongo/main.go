package main

import (
	"context"
	mgo "coolcar/server/share/mongo"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://127.0.0.1:27017/?readPreference=primary&ssl=false&directConnection=true"))
	if err != nil {
		panic(err)
	}
	// mongodb 存的是bson的格式
	col := mc.Database("coolcar").Collection("account")
	findRows(c, col)
}

func findRows(c context.Context, col *mongo.Collection) {
	//res := col.FindOne(c, bson.M{
	//	"open_id": "123",
	//})
	//fmt.Printf("%+v\n", res)
	//var row struct {
	//	ID     primitive.ObjectID `bson:"_id"`
	//	OpenID string             `bson:"open_id"`
	//}
	//err := res.Decode(&row)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v\n", row)

	cur, err := col.Find(c, bson.M{})
	if err != nil {
		panic(err)
	}
	for cur.Next(c) {
		var row struct {
			mgo.ObjID
			OpenID string `bson:"open_id"`
		}
		err := cur.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", row)
	}
}

func insertRows(c context.Context, col *mongo.Collection) {
	res, err := col.InsertMany(c, []interface{}{ // []interface{} 任意类型的数组
		bson.M{
			"open_id": "123",
		}, bson.M{
			"open_id": "456",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", res)
}
