package dao

import (
	"context"
	rentalpb "coolcar/server/rental/api/gen/v1"
	"coolcar/server/share/id"
	mgutil "coolcar/server/share/mongo"
	"coolcar/server/share/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	tripField      = "trip"
	accountIDField = tripField + ".accountid"
	statusField    = tripField + ".status"
)

// Mongo defines a mongo dao.
type Mongo struct {
	col      *mongo.Collection
	newObjID func() primitive.ObjectID // 生成ID的函数
}

//TODO...: 同一个account最多只能有一个进行中的Trip
//TODO...：强类型化tripID
// TODO...表格驱动测试

// NewMongo creates a new mongo dao.
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:      db.Collection("trip"), // 注意这里切库
		newObjID: primitive.NewObjectID,
	}
}

// TripRecord 定义表的结构
// TripRecord defines a trip record in mongo db
type TripRecord struct {
	mgutil.ObjID         `bson:"inline"` // todo... bson:"inline"????
	mgutil.UpdateAtField `bson:"inline"`
	Trip                 *rentalpb.Trip `bson:"trip"`
}

// CreateTrip creates a trip
func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = mgutil.NewObjID()
	r.UpdatedAt = mgutil.UpdatedAt()

	_, err := m.col.InsertOne(c, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTrip gets a trip.
func (m *Mongo) GetTrip(c context.Context, id id.TripID, accountID id.AccountID) (*TripRecord, error) {
	objID, err := objid.FromID(id) //primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}
	res := m.col.FindOne(c, bson.M{
		mgutil.IDFieldName: objID,
		accountIDField:     accountID,
	})

	if err := res.Err(); err != nil {
		return nil, err
	}

	var tr TripRecord
	err = res.Decode(&tr)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %v\n", err)
	}
	return &tr, nil
}

// GetTrips gets trips for the account by status.
// If status is not specified, gets all trips for the account.
func (m *Mongo) GetTrips(c context.Context, accountID id.AccountID, status rentalpb.TripStatus) ([]*TripRecord, error) {
	filter := bson.M{
		accountIDField: accountID.String(),
	}
	if status != rentalpb.TripStatus_TS_NOT_SPECIFIED {
		filter[statusField] = status
	}

	res, err := m.col.Find(c, filter)
	if err != nil {
		return nil, err
	}

	var trips []*TripRecord
	for res.Next(c) {
		var trip TripRecord
		err := res.Decode(&trip)
		if err != nil {
			//return nil, err
			continue
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

/**
使用 updatedAt 乐观锁解决同时更新的问题
**/
// UpdateTrip updates a trip.
func (m *Mongo) UpdateTrip(c context.Context, tid id.TripID, aid id.AccountID, updatedAt int64, trip *rentalpb.Trip) error {
	objID, err := objid.FromID(tid)
	if err != nil {
		return fmt.Errorf("invalid id: %v\n", err)
	}

	newUpdatedAt := mgutil.UpdatedAt()
	res, err := m.col.UpdateOne(c, bson.M{
		mgutil.IDFieldName:        objID,
		accountIDField:            aid.String(),
		mgutil.UpdatedAtFieldName: updatedAt,
	}, mgutil.Set(bson.M{
		tripField:                 trip,
		mgutil.UpdatedAtFieldName: newUpdatedAt,
	}))
	if err != nil {
		return err
	}

	// 可能找不到任何文档
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
