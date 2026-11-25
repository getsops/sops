package keyservice

import (
	"context"

	"google.golang.org/grpc"
)

// LocalClient is a key service client that performs all operations locally
type LocalClient struct {
	Server KeyServiceServer
}

// NewLocalClient creates a local client with its own plugin runtime. Call
// Close when finished to stop plugin processes owned by the client.
func NewLocalClient() LocalClient {
	return NewLocalClientWithContext(context.Background())
}

// NewLocalClientWithContext creates a local client whose plugin processes are
// stopped when ctx is canceled or the client is closed.
func NewLocalClientWithContext(ctx context.Context) LocalClient {
	return LocalClient{Server: NewServer(ctx)}
}

// NewCustomLocalClient creates a new local client with a non-default backing
// KeyServiceServer implementation
func NewCustomLocalClient(server KeyServiceServer) LocalClient {
	return LocalClient{Server: server}
}

// Close releases resources owned by the local key-service server, if any.
func (c LocalClient) Close() error {
	if closer, ok := c.Server.(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
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

type contextualClient struct {
	client KeyServiceClient
	ctx    context.Context
}

// WithContext returns a client that also cancels requests when ctx is done.
// The request's own deadline and cancellation remain in effect.
func WithContext(client KeyServiceClient, ctx context.Context) KeyServiceClient {
	if ctx == nil {
		return client
	}
	return contextualClient{client: client, ctx: ctx}
}

func (c contextualClient) Encrypt(ctx context.Context, req *EncryptRequest, opts ...grpc.CallOption) (*EncryptResponse, error) {
	requestContext, cancel := combineContexts(ctx, c.ctx)
	defer cancel()
	return c.client.Encrypt(requestContext, req, opts...)
}

func (c contextualClient) Decrypt(ctx context.Context, req *DecryptRequest, opts ...grpc.CallOption) (*DecryptResponse, error) {
	requestContext, cancel := combineContexts(ctx, c.ctx)
	defer cancel()
	return c.client.Decrypt(requestContext, req, opts...)
}

func combineContexts(request, operation context.Context) (context.Context, context.CancelFunc) {
	if request == nil {
		request = context.Background()
	}
	combined, cancel := context.WithCancel(request)
	stop := context.AfterFunc(operation, cancel)
	if operation.Err() != nil {
		cancel()
	}
	return combined, func() {
		stop()
		cancel()
	}
}
