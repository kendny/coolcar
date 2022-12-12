package main

import (
	"context"
	blobpb "coolcar/server/blob/api/gen/v1"
	"fmt"
	"google.golang.org/grpc"
)

func main() {
	//grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial("localhost:8083", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	c := blobpb.NewBlobServiceClient(conn)

	ctx := context.Background()
	//// 获取上传图片的地址
	//res, err := c.CreateBlob(ctx, &blobpb.CreateBlobRequest{
	//	AccountId:           "account_2",
	//	UploadUrlTimeoutSec: 1000,
	//})

	//// 获取上传图片
	//res, err := c.GetBlob(ctx, &blobpb.GetBlobRequest{
	//	Id: "639446f1ce6ee5f76eab5220",
	//})

	// 获取预签名的URL并进行展示
	res, err := c.GetBlobURL(ctx, &blobpb.GetBlobURLRequest{
		Id:         "639446f1ce6ee5f76eab5220",
		TimeoutSec: 2000,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", res)
}
