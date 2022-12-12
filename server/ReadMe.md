## 笔记

```go
mux := runtime.NewServeMux(runtime.WithMarshalerOption(
runtime.MIMEWildcard, &runtime.JSONPb{
MarshalOptions: protojson.MarshalOptions{
UseEnumNumbers: true, // 将枚举类型生成数值
UseProtoNames:  true,
},
UnmarshalOptions: protojson.UnmarshalOptions{
DiscardUnknown: true,
},
},
))
```

- 说明：GRPC工具链在2020年11月有一次较大的更新。针对这个问题

```
属性名称改变：EnumsAsInts -> MarshalOptions.UseEnumNumbers，OrigName -> MarshalOptions.UseProtoNames，

增加属性：UnmarshalOptions. DiscardUnknown = true，该属性用来忽略proto中不存在的字段
```

参考： https://git.imooc.com/class-108/coolcar/pulls/1

https://github.com/grpc/grpc/releases/tag/
https://github.com/storyicon/powerproto

### as some methods are missing: mustEmbedUnimplementedAuthServiceServer()

解决方案： protoc版本上的变动和机制的改变导致了这个问题 1.protoc加一个flag（require_unimplemented_servers=false）就可以了： protoc -I$GOPATH/src -I.
--go-grpc_out=require_unimplemented_servers=false:$GOPATH/src *.proto
2.具体信息参照https://github.com/grpc/grpc-go/blob/master/cmd/protoc-gen-go-grpc/README.md

为啥加？

```shell
--go-grpc_out=require_unimplemented_servers=false:$GO_OUT_PATH
```

3. 定义service时加 UnimplementedAuthServiceServer //所有实现必须嵌入UnimplementedTripServiceServer //向前兼容

```shell
type Service struct {
	OpenIDResolver                        OpenIDResolver
	Mongo                                 *dao.Mongo
	TokenGenerator                        TokenGenerator
	TokenExpire                           time.Duration
	Logger                                *zap.Logger // 服务类型一般都用指针
	authpb.UnimplementedAuthServiceServer             // 必须引用，不然报错
}
```

### 编译 ts

在packpage.json中添加 编译 命令

```shell
 "scripts": {
      "compile": "./node_modules/typescript/bin/tsc",
      "tsc": "node ./node_modules/typescript/lib/tsc.js"
  }
  
// 运行下面命令  
npm run tsc
```

### go 操作MongoDB 报错原因

```shell
panic: error decoding key _id: an ObjectID string must be exactly 12 bytes long (got 5)
```

> ```当然该数据库中还有其他数据，因为ObjectID的原因，程序panic了```
参考：https://learnku.com/articles/66231

#### mongo 测试用例报错

```shell
panic: the provided hex string is not a valid ObjectID [recovered]
	panic: the provided hex string is not a valid ObjectID
```

解决： mongo的主键ID mgo.IDField：长度不对导致的

```shell
bson.M{
			mgutil.IDField: mustObjID("5f7c245ab0361e00ffb9fd6f11"), // 会报错长度应该是 5f7c245ab0361e00ffb9fd6f
			openIDField: "openid_1",
		},
```

### --grpc-gateway_out: must not set request body when http method is GET: CreateTrip

```shell
http:
    rules:
    - selector: rental.v1.TripService.CreateTrip
      get: /v1/trip
      body: "*"
```

### --grpc-gateway_out: failed to parse gRPC API Configuration from YAML in './rental/api/rental.yaml': proto: (line 1:67): error parsing "post", oneof google.api.HttpRule.pattern is already set

```shell
http:
    rules:
    - selector: rental.v1.TripService.CreateTrip
      post: /v1/trip
      get: /v1/trip
      body: "*"
```

### 如何保证 同一个account最多只能有一个进行中的Trip

```js
db.trip.createIndex({
    "trip.accountid": 1,
    "trip.status": 1,
}, {
    unique: true,
    partialFilterExpression: {
        "trip.status": 1, // 指的是值为1
    }
})
```

### {"code":12,"message":"unknown service rental.v1.TripService"}

在gateway上服务没有注册成功

```shell
	for _, s := range serverConfig {
		err := s.registerFunc(c,
			mux, // mux:multiplexer
			"localhost:8081",
			[]grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			lg.Sugar().Fatalf("cannot register service %s:%v\n", s.name, s.addr)
		}
	}
```

### docker 启动mongo 跑测试失败

需要查看docker mongo 测试启动日志，可能是空间不足导致
[docker清理缓存](https://blog.csdn.net/m0_67390963/article/details/126327604)

### panic: runtime error: invalid memory address or nil pointer dereference

[signal SIGSEGV: segmentation violation code=0x1 addr=0x18 pc=0x1771afd]
> “某个指针变量声明之后, 还没有经过赋值时默认指向nil, 直接调用指针就会报错”

```shell
rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager: &car.Manager{},
				ProfileManager: &profClient.Manager{},
				POIManager: &poi.Manager{},
				Mongo:      tripdao.NewMongo(db),
				Logger:     logger,
			})
```

> 上面的 &profClient.Manager{} 未挂服务

```shell
	rentalpb.RegisterTripServiceServer(s, &trip.Service{
				CarManager: &car.Manager{},
				ProfileManager: &profClient.Manager{
					Fetcher: profService,
				},
				POIManager: &poi.Manager{},
				Mongo:      tripdao.NewMongo(db),
				Logger:     logger,
			})
```

### 微信小程序坑点之一

小程序上传接口

```shell
wx.uploaldFile()
```

> 中的method是小写的post，导致其连自身的对象存储服务都不能很好的对接，需要使用wx.request进行文件的上传

### imports coolcar/server/blob/api/gen/v1: import cycle not allowed

```shell
可能导致的原因，go_package包声明路径不对
```

### panic: grpc: no transport security set (use grpc.WithTransportCredentials(insecure.NewCredentials()) explicitly or set credentials)

```shell
	conn, err := grpc.Dial("localhost:8083")
```

改成

```shell
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
#	或
    conn, err := grpc.Dial("localhost:8083", grpc.WithTransportCredentials(insecure.NewCredentials()))
```

#### Mongo 更新语句

> 存在更新，不存在插入

```JS
db.profile.update({
    accountid: "xxxx",
    "$or": [
        {
            "profile.identitystatus": 0,
        }, {
            "profile.identitystatus": {
                "$exists": false,
            },
        },
    ],
}, {
    "$set": {
        accountid: 'xxxx',
        profile: {
            identity: {
                name: "abc",
            },
            identitystatus: 1
        }
    }
}, {
    upsert: true
})
```