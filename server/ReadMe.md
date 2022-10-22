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

###  as some methods are missing: mustEmbedUnimplementedAuthServiceServer()
解决方案：
protoc版本上的变动和机制的改变导致了这个问题
1.protoc加一个flag（require_unimplemented_servers=false）就可以了：
protoc -I$GOPATH/src -I.
--go-grpc_out=require_unimplemented_servers=false:$GOPATH/src *.proto
2.具体信息参照https://github.com/grpc/grpc-go/blob/master/cmd/protoc-gen-go-grpc/README.md

为啥加？
```shell
--go-grpc_out=require_unimplemented_servers=false:$GO_OUT_PATH
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

