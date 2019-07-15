package ssh

import (
	"context"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathFetchPublicKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: `public_key`,

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathFetchPublicKey,
		},

		HelpSynopsis:    `Retrieve the public key.`,
		HelpDescription: `This allows the public key, that this backend has been configured with, to be fetched.`,
	}
}

func (b *backend) pathFetchPublicKey(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	publicKeyEntry, err := caKey(ctx, req.Storage, caPublicKey)
	if err != nil {
		return nil, err
	}
	if publicKeyEntry == nil || publicKeyEntry.Key == "" {
		return nil, nil
	}

	response := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPContentType: "text/plain",
			logical.HTTPRawBody:     []byte(publicKeyEntry.Key),
			logical.HTTPStatusCode:  200,
		},
	}

	return response, nil
}
