package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	client_auth "grpc/middleware/client"
	"grpc/sample/server/pb"
	"time"
)

func main() {
	// 建立网络链接
	// 添加认证信息
	credential := client_auth.NewAuthentication("admin", "123456")
	conn, err := grpc.DialContext(context.Background(), "127.0.0.1:1234", grpc.WithInsecure(), grpc.WithPerRPCCredentials(credential))
	if err != nil {
		panic(err)
	}

	client := pb.NewHelloServiceClient(conn)

	// 添加认证凭证信息
	// 方式一、每个请求都要初始化一个凭证，有些麻烦
	//MD_Credential := server.NewClientCredential("admin", "123456")
	//ctx := metadata.NewOutgoingContext(context.Background(), MD_Credential)

	// 方式二、批量初始化凭证，每个请求自动添加凭证
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs())

	resp, err := client.Hello(ctx, &pb.Request{Value: "bob"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Value)

	stream, err := client.Channel(ctx)
	if err != nil {
		panic(err)
	}

	// 启用一个goroutine 来发送请求
	go func() {
		for {
			err := stream.Send(&pb.Request{Value: "alice"})
			time.Sleep(1 * time.Second)
			if err != nil {
				panic(err)
			}
		}
	}()

	for {
		recv, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Println(recv.Value)
	}

}
