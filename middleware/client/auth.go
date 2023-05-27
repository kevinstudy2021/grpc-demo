package client

import (
	"context"
	"grpc/middleware/server"
)

//
//type PerRPCCredentials interface {
//	// GetRequestMetadata gets the current request metadata, refreshing tokens
//	// if required. This should be called by the transport layer on each
//	// request, and the data should be populated in headers or other
//	// context. If a status code is returned, it will be used as the status for
//	// the RPC (restricted to an allowable set of codes as defined by gRFC
//	// A54). uri is the URI of the entry point for the request.  When supported
//	// by the underlying implementation, ctx can be used for timeout and
//	// cancellation. Additionally, RequestInfo data will be available via ctx
//	// to this call.  TODO(zhaoq): Define the set of the qualified keys instead
//	// of leaving it as an arbitrary string.
//	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)
//	// RequireTransportSecurity indicates whether the credentials requires
//	// transport security.
//	RequireTransportSecurity() bool
//}

func NewAuthentication(ak, sk string) *Authentication {
	return &Authentication{
		clientId:     ak,
		clientSecret: sk,
	}
}

type Authentication struct {
	clientId     string
	clientSecret string
}

func (a *Authentication) build() map[string]string {
	return map[string]string{
		server.ClientHeaderAccessKey: a.clientId,
		server.ClientHeaderSecretKey: a.clientSecret,
	}
}

func (a *Authentication) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return a.build(), nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false

}
