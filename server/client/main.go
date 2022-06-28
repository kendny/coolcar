/*
 * @Author: kendny wh_kendny@163.com
 * @Date: 2022-06-19 16:18:39
 * @LastEditors: kendny wh_kendny@163.com
 * @LastEditTime: 2022-06-19 17:03:02
 * @FilePath: /wx/Users/xxxian/go_project/src/coolcar/server/client/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"context"
	trippb "coolcar/server/proto/gen/go"
	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	log.SetFlags(log.Lshortfile)
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannnot connect server: %v", err)
	}

	tsClient := trippb.NewTripServiceClient(conn)
	r, err := tsClient.GetTrip(context.Background(), &trippb.GetTripRequest{
		Id: "trip457",
	})
	if err != nil {
		log.Fatalf("cantnot call GetTrip:%v", err)
	}
	fmt.Println(r)
}
