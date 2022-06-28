###
 # @Author: kendny wh_kendny@163.com
 # @Date: 2022-06-20 21:56:05
 # @LastEditors: kendny wh_kendny@163.com
 # @LastEditTime: 2022-06-20 22:03:55
 # @FilePath: /coolcar/server/proto/gen/gen.sh
 # @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
### 

protoc -I=. --go_out=gen/go --go_opt=paths=source_relative --go-grpc_out=gen/go  --go-grpc_opt=paths=source_relative  ./trip.proto 
protoc -I=. --grpc-gateway_out=paths=source_relative,grpc_api_configuration=trip.yaml:gen/go ./trip.proto
