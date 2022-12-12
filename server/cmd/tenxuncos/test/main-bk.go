package main

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"time"
)

func main() {
	// 将 examplebucket-1250000000 和 COS_REGION 修改为真实的信息
	// 存储桶名称，由bucketname-appid 组成，appid必须填入，可以在COS控制台查看存储桶名称。https://console.cloud.tencent.com/cos5/bucket
	// COS_REGION 可以在控制台查看，https://console.cloud.tencent.com/cos5/bucket, 关于地域的详情见 https://cloud.tencent.com/document/product/436/6224
	u, err := url.Parse("https://wuhan-1259722894.cos.ap-shanghai.myqcloud.com")
	if err != nil {
		panic(err)
	}
	b := &cos.BaseURL{BucketURL: u}
	sID := "AKIDbkfNr78vUq32pOhoiQxHMDpDPPESeicR"
	sKey := "jHiz62oV8lK1Zv78yeEGE90hHl48zc1B"
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  sID,  // 替换为用户的 SecretId，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
			SecretKey: sKey, // 替换为用户的 SecretKey，请登录访问管理控制台进行查看和管理，https://console.cloud.tencent.com/cam/capi
		},
	})
	// 对象键（Key）是对象在存储桶中的唯一标识。
	// 例如，在对象的访问域名 `examplebucket-1250000000.cos.COS_REGION.myqcloud.com/test/objectPut.go` 中，对象键为 test/objectPut.go
	name := "1.png" // 已经上传
	ctx := context.Background()
	// 获取预签名URL
	presignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodPut, name, sID, sKey, 5*time.Minute, nil)
	if err != nil {
		panic(err)
	}
	//生成上传下载链接， MethodGet： 生成的URL可用来下载， MethodPut: 生成的URL可用来下载
	fmt.Println(presignedURL)

	//// 1.通过字符串上传对象
	//f := strings.NewReader("./")
	//_, err = c.Object.Put(ctx, name, f, nil)
	//if err != nil {
	//	panic(err)
	//}

	//// 2.通过本地文件上传对象
	//_, err = c.Object.PutFromFile(context.Background(), name, "../test", nil)
	//if err != nil {
	//	panic(err)
	//}
	//// 3.通过文件流上传对象
	//fd, err := os.Open("./test")
	//if err != nil {
	//	panic(err)
	//}
	//defer fd.Close()
	//_, err = c.Object.Put(context.Background(), name, fd, nil)
	//if err != nil {
	//	panic(err)
	//}
}
