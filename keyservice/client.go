package keyservice

import (
	"golang.org/x/net/context"

	"google.golang.org/grpc"
)

// LocalClient is a key service client that performs all operations locally
type LocalClient struct {
	Server Server
}

// NewLocalClient creates a new local client
func NewLocalClient() LocalClient {
	return LocalClient{Server{}}
}

// Decrypt processes a decrypt request locally
// See keyservice/server.go for more details
func (c LocalClient) Decrypt(ctx context.Context,
	req *DecryptRequest, opts ...grpc.CallOption) (*DecryptResponse, error) {
	return c.Server.Decrypt(ctx, req)
}

// Encrypt processes an encrypt request locally
// See keyservice/server.go for more details
func (c LocalClient) Encrypt(ctx context.Context,
	req *EncryptRequest, opts ...grpc.CallOption) (*EncryptResponse, error) {
	return c.Server.Encrypt(ctx, req)
}
