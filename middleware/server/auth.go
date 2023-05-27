package server

import (
	"context"
	"fmt"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	ClientHeaderAccessKey = "client-id"
	ClientHeaderSecretKey = "client-secret"
)

func NewClientCredential(ak, sk string) metadata.MD {
	return metadata.MD{
		ClientHeaderAccessKey: []string{ak},
		ClientHeaderSecretKey: []string{sk},
	}
}

func NewAuthUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	//grpcAuther := &GrpcAuther{}
	//return grpcAuther.Auth

	return (&GrpcAuther{}).UnaryServerInterceptor
}

func NewAuthStreamServerInterceptor() grpc.StreamServerInterceptor {
	return (&GrpcAuther{}).StreamServerInterceptorfunc
}

type grpcAuther struct {
	log logger.Logger
}

func newGrpcAuthLogger() *grpcAuther {
	return &grpcAuther{
		log: zap.L().Named("Grpc Auther"),
	}
}

// stream rpc interceptor
func (a *GrpcAuther) StreamServerInterceptorfunc(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	md, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return fmt.Errorf("ctx is not an grpc incoming context")
	}

	// 从metadata中获取客户端传递过来的凭证
	clientId, clientSecret := a.getClientCredentialsFromMeta(md)

	if err := a.validateServiceCredentian(clientId, clientSecret); err != nil {
		return err
	}
	return handler(srv, ss)
}

type GrpcAuther struct {
}

// request-response interceptor
func (a *GrpcAuther) UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	authLog := newGrpcAuthLogger()
	authLog.log.Infof("get client id and client secret from context")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("ctx is not an grpc incoming context")
	}

	// 从metadata中获取客户端传递过来的凭证
	clientId, clientSecret := a.getClientCredentialsFromMeta(md)

	if err := a.validateServiceCredentian(clientId, clientSecret); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func (a *GrpcAuther) getClientCredentialsFromMeta(md metadata.MD) (clientId, clientSecret string) {
	cakList := md[ClientHeaderAccessKey]
	if len(cakList) > 0 {
		clientId = cakList[0]
	}
	cskList := md[ClientHeaderSecretKey]
	if len(cskList) > 0 {
		clientSecret = cskList[0]
	}

	return clientId, clientSecret
}

func (a *GrpcAuther) validateServiceCredentian(clientId, clientSecret string) error {
	if !(clientId == "admin" && clientSecret == "123456") {
		// 返回认证错误，结束rpc调用
		return status.Errorf(codes.Unauthenticated, "client id or client secret not correct")
	}
	return nil
}
