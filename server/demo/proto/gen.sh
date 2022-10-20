###
# @Author: kendny wh_kendny@163.com
# @Date: 2022-06-20 21:56:05
# @LastEditors: kendny wh_kendny@163.com
# @LastEditTime: 2022-09-30 22:53:17
# @FilePath: /coolcar/server/proto/gen/gen.sh
# @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
###

protoc -I=. --go_out=gen/go --go_opt=paths=source_relative --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative ./trip.proto

# 生成 gateway文件
protoc -I=. --grpc-gateway_out=paths=source_relative,grpc_api_configuration=trip.yaml:gen/go ./trip.proto

PBTS_BIN_DIR=../../wx/miniprogram/node_modules/.bin
PBTS_OUT_DIR=../../wx/miniprogram/service/proto_gen

$PBTS_BIN_DIR/pbjs -t static -w es6 trip.proto --no-create --no-encode --no-decode --no-verify --no-delimited --force-number -o $PBTS_OUT_DIR/trip_pb_tmp.js

echo 'import * as $protobuf from "protobufjs";' >$PBTS_OUT_DIR/trip_pb.js
cat $PBTS_OUT_DIR/trip_pb_tmp.js >>$PBTS_OUT_DIR/trip_pb.js

rm $PBTS_OUT_DIR/trip_pb_tmp.js

$PBTS_BIN_DIR/pbts -o $PBTS_OUT_DIR/trip_pb.d.ts $PBTS_OUT_DIR/trip_pb.js

