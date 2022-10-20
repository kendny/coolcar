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

