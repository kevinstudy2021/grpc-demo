package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc/sample/server/pb"
	"io"
	"log"
	"net"
)

// HelloServiceServer is the server API for HelloService service.
// All implementations must embed UnimplementedHelloServiceServer
// for forward compatibility
//type HelloServiceServer interface {
//	Hello(context.Context, *Request) (*Response, error)
//	mustEmbedUnimplementedHelloServiceServer()
//}

// HelloServiceServer must be embedded to have forward compatible implementations.
type HelloServiceServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *HelloServiceServer) Hello(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return &pb.Response{Value: fmt.Sprintf("hello, " + req.Value)}, nil
}

func (s *HelloServiceServer) Channel(stream pb.HelloService_ChannelServer) error {
	for {
		// 接受client请求
		req, err := stream.Recv()
		if err != nil {
			log.Printf("recv error , %s", err)
			if err == io.EOF {
				log.Printf("client closed")
				return nil
			}
			return err
		}

		resp := &pb.Response{Value: fmt.Sprintf("hello, %s", req.Value)}

		// 响应client请求
		err = stream.Send(resp)
		if err != nil {
			log.Printf("send error  %s", err)
			if err == io.EOF {
				log.Printf("client closed")
				return nil
			}
			return err
		}

	}
}

func main() {
	server := grpc.NewServer()

	// 把实现类注册到grpc server
	pb.RegisterHelloServiceServer(server, new(HelloServiceServer))

	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	log.Printf("grpc listen addr: 127.0.0.1:1234")

	if err := server.Serve(listen); err != nil {
		panic(err)
	}

}
