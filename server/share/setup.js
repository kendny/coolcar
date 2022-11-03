db.account.createIndex({
    open_id: 1,
},{
    unique: true
})

// 设置mongo相关
// 同一个account最多只能有一个进行中的Trip
db.trip.createIndex({
    "trip.accountid": 1,
    "trip.status": 1,
}, {
    unique: true,
    partialFilterExpression: {
        "trip.status": 1, // 指的是值为1
    }
})