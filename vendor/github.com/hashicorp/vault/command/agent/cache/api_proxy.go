package cache

import (
	"context"
	"fmt"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
)

// APIProxy is an implementation of the proxier interface that is used to
// forward the request to Vault and get the response.
type APIProxy struct {
	client *api.Client
	logger hclog.Logger
}

type APIProxyConfig struct {
	Client *api.Client
	Logger hclog.Logger
}

func NewAPIProxy(config *APIProxyConfig) (Proxier, error) {
	if config.Client == nil {
		return nil, fmt.Errorf("nil API client")
	}
	return &APIProxy{
		client: config.Client,
		logger: config.Logger,
	}, nil
}

func (ap *APIProxy) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	client, err := ap.client.Clone()
	if err != nil {
		return nil, err
	}
	client.SetToken(req.Token)
	client.SetHeaders(req.Request.Header)

	fwReq := client.NewRequest(req.Request.Method, req.Request.URL.Path)
	fwReq.BodyBytes = req.RequestBody

	query := req.Request.URL.Query()
	if len(query) != 0 {
		fwReq.Params = query
	}

	// Make the request to Vault and get the response
	ap.logger.Info("forwarding request", "method", req.Request.Method, "path", req.Request.URL.Path)

	resp, err := client.RawRequestWithContext(ctx, fwReq)
	if resp == nil && err != nil {
		// We don't want to cache nil responses, so we simply return the error
		return nil, err
	}

	// Before error checking from the request call, we'd want to initialize a SendResponse to
	// potentially return
	sendResponse, newErr := NewSendResponse(resp, nil)
	if newErr != nil {
		return nil, newErr
	}

	// Bubble back the api.Response as well for error checking/handling at the handler layer.
	return sendResponse, err
}
