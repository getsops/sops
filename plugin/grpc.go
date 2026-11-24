package plugin

import (
	"context"
	"fmt"

	"github.com/getsops/sops/v3/plugin/protocol"
	hplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type grpcPlugin struct {
	hplugin.NetRPCUnsupportedPlugin
	Impl KeyProvider
}

func (p *grpcPlugin) GRPCServer(_ *hplugin.GRPCBroker, server *grpc.Server) error {
	protocol.RegisterKeyProviderServer(server, &grpcServer{provider: p.Impl})
	return nil
}

func (p *grpcPlugin) GRPCClient(ctx context.Context, _ *hplugin.GRPCBroker, conn *grpc.ClientConn) (interface{}, error) {
	return &grpcClient{client: protocol.NewKeyProviderClient(conn)}, nil
}

type grpcServer struct {
	protocol.UnimplementedKeyProviderServer
	provider KeyProvider
}

func (s *grpcServer) Wrap(ctx context.Context, request *protocol.WrapKeyRequest) (*protocol.WrapKeyResponse, error) {
	config, err := decodeObject(request.ConfigJson)
	if err != nil {
		return nil, fmt.Errorf("decode plugin configuration: %w", err)
	}
	response, err := s.provider.Wrap(ctx, WrapKeyRequest{Config: config, DataKey: request.DataKey})
	if err != nil {
		return nil, err
	}
	metadata, err := encodeObject(response.Metadata)
	if err != nil {
		return nil, fmt.Errorf("encode plugin metadata: %w", err)
	}
	return &protocol.WrapKeyResponse{
		EncryptedDataKey: response.EncryptedDataKey,
		MetadataJson:     metadata,
	}, nil
}

func (s *grpcServer) Unwrap(ctx context.Context, request *protocol.UnwrapKeyRequest) (*protocol.UnwrapKeyResponse, error) {
	config, err := decodeObject(request.ConfigJson)
	if err != nil {
		return nil, fmt.Errorf("decode plugin configuration: %w", err)
	}
	metadata, err := decodeObject(request.MetadataJson)
	if err != nil {
		return nil, fmt.Errorf("decode plugin metadata: %w", err)
	}
	response, err := s.provider.Unwrap(ctx, UnwrapKeyRequest{
		Config:           config,
		Metadata:         metadata,
		EncryptedDataKey: request.EncryptedDataKey,
	})
	if err != nil {
		return nil, err
	}
	return &protocol.UnwrapKeyResponse{DataKey: response.DataKey}, nil
}

type grpcClient struct {
	client protocol.KeyProviderClient
}

func (c *grpcClient) Wrap(ctx context.Context, request WrapKeyRequest) (WrapKeyResponse, error) {
	config, err := encodeObject(request.Config)
	if err != nil {
		return WrapKeyResponse{}, err
	}
	response, err := c.client.Wrap(ctx, &protocol.WrapKeyRequest{
		ConfigJson: config,
		DataKey:    request.DataKey,
	})
	if err != nil {
		return WrapKeyResponse{}, err
	}
	metadata, err := decodeObject(response.MetadataJson)
	if err != nil {
		return WrapKeyResponse{}, fmt.Errorf("decode plugin metadata: %w", err)
	}
	return WrapKeyResponse{
		EncryptedDataKey: response.EncryptedDataKey,
		Metadata:         metadata,
	}, nil
}

func (c *grpcClient) Unwrap(ctx context.Context, request UnwrapKeyRequest) (UnwrapKeyResponse, error) {
	config, err := encodeObject(request.Config)
	if err != nil {
		return UnwrapKeyResponse{}, err
	}
	metadata, err := encodeObject(request.Metadata)
	if err != nil {
		return UnwrapKeyResponse{}, err
	}
	response, err := c.client.Unwrap(ctx, &protocol.UnwrapKeyRequest{
		ConfigJson:       config,
		MetadataJson:     metadata,
		EncryptedDataKey: request.EncryptedDataKey,
	})
	if err != nil {
		return UnwrapKeyResponse{}, err
	}
	return UnwrapKeyResponse{DataKey: response.DataKey}, nil
}
