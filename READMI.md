```shell
# 所在目录  D:\go-study\mage-go-project\grpc\sample\server
# protoc -I="." --go_out=.  --go_opt=module="grpc/sample/server" pb/hello.proto
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 补充rpc 接口定义protobuf的代码生成
# protoc -I="." --go_out=.  --go_opt=module="grpc/sample/server" --go-grpc_out=. --go-grpc_opt=module="grpc/sample/server"  pb/hello.proto
```








