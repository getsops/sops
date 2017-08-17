package keyservice

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

type LocalClient struct {
	Server Server
}

func NewLocalClient() LocalClient {
	return LocalClient{Server{}}
}

func (c LocalClient) Decrypt(ctx context.Context,
	req *DecryptRequest, opts ...grpc.CallOption) (*DecryptResponse, error) {
	return c.Server.Decrypt(ctx, req)
}

func (c LocalClient) Encrypt(ctx context.Context,
	req *EncryptRequest, opts ...grpc.CallOption) (*EncryptResponse, error) {
	return c.Server.Encrypt(ctx, req)
}
