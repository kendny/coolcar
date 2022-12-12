package mgutil

//common field names
import (
	"coolcar/server/share/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// 将mongo中的操作进行提取共享

const (
	IDFieldName        = "_id"       // IDFieldName defineds the field name for mongo  document id.
	UpdatedAtFieldName = "updatedat" // 千万不能写错了，要不出错难排查
)

// ObjID defines the object id field
type ObjID struct {
	ID primitive.ObjectID `bson:"_id"`
}

type UpdateAtField struct {
	UpdatedAt int64 `bson:"updatedat"`
}

// NewObjID 封装一个返回确定值的方法
// NewObjID generates a new object id.
var NewObjID = primitive.NewObjectID

// NewObjIDWithValue sets id for next objectID generation.
func NewObjIDWithValue(id fmt.Stringer) {
	NewObjID = func() primitive.ObjectID {
		return objid.MustFromID(id)
	}
}

// UpdatedAt returns a value suitable for UpdatedAt field.
var UpdatedAt = func() int64 {
	return time.Now().UnixNano()
}

// Set returns a $set update document
func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

// SetOnInsert 不存在就插入并返回，存在直接查询返回
func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}

// ZeroOrDoesNotExist generates a filter expression with
// filed equal to zero or field does not exits
func ZeroOrDoesNotExist(field string, zero interface{}) bson.M {
	return bson.M{
		"$or": []bson.M{
			{
				field: zero,
			},
			{
				field: bson.M{
					"$exists": false,
				},
			},
		},
	}
}
