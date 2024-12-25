package cloudru

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	iamAuthV1 "github.com/cloudru-tech/iam-sdk/api/auth/v1"
	kmsV1 "github.com/cloudru-tech/key-manager-sdk/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
)

// EndpointsResponse is a response from the Cloud.ru API.
type EndpointsResponse struct {
	// Endpoints contains the list of actual API addresses of Cloud.ru products.
	Endpoints []Endpoint `json:"endpoints"`
}

// Endpoint is a product API address.
type Endpoint struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

type Client struct {
	KMS     kmsV1.KeyManagerServiceClient
	kmsConn *grpc.ClientConn
}

func provideClient() (*Client, error) {
	discoveryURL := DiscoveryURL

	if du, ok := os.LookupEnv(EnvDiscoveryURL); ok {
		u, err := url.Parse(discoveryURL)
		if err != nil {
			return nil, fmt.Errorf("invalid %s param: %w", EnvDiscoveryURL, err)
		}

		switch {
		case u.Host == "":
			return nil, fmt.Errorf("invalid %s param: missing host", EnvDiscoveryURL)
		case u.Scheme != "http", u.Scheme != "https":
			return nil, fmt.Errorf("invalid %s param: scheme must be http or https", EnvDiscoveryURL)
		}

		discoveryURL = du
	}

	var ok bool
	var akID, akSecret string
	if akID, ok = os.LookupEnv(EnvAccessKeyID); !ok {
		return nil, fmt.Errorf("missing %s env param", EnvAccessKeyID)
	}
	if akSecret, ok = os.LookupEnv(EnvAccessKeySecret); !ok {
		return nil, fmt.Errorf("missing %s env param", EnvAccessKeySecret)
	}

	endpoints, err := getEndpoints(discoveryURL)
	if err != nil {
		return nil, err
	}

	kmsEndpoint := endpoints.Get("key-manager")
	if kmsEndpoint == nil {
		return nil, errors.New("key-manager API is not available")
	}

	iamEndpoint := endpoints.Get("iam")
	if iamEndpoint == nil {
		return nil, errors.New("iam API is not available")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	iamConn, err := grpc.NewClient(iamEndpoint.Address,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 30,
			Timeout:             time.Second * 5,
			PermitWithoutStream: false,
		}),
		grpc.WithUserAgent("sops"),
	)
	if err != nil {
		return nil, fmt.Errorf("initiate IAM gRPC connection: %w", err)
	}
	defer iamConn.Close()

	iam := iamAuthV1.NewAuthServiceClient(iamConn)
	token, err := iam.GetToken(ctx, &iamAuthV1.GetTokenRequest{
		KeyId:  akID,
		Secret: akSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	kmsConn, err := grpc.NewClient(kmsEndpoint.Address,
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS13})),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Second * 30,
			Timeout:             time.Second * 5,
			PermitWithoutStream: false,
		}),
		grpc.WithUserAgent("sops"),
		grpc.WithUnaryInterceptor(func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			md, ok := metadata.FromOutgoingContext(ctx)
			if !ok {
				md = metadata.New(map[string]string{})
			}
			md.Set("authorization", "Bearer "+token.AccessToken)

			return invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("initiate KMS gRPC connection: %w", err)
	}

	return &Client{
		KMS:     kmsV1.NewKeyManagerServiceClient(kmsConn),
		kmsConn: kmsConn,
	}, nil
}

// getEndpoints returns the actual Cloud.ru API endpoints.
func getEndpoints(url string) (*EndpointsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("construct HTTP request for cloud.ru endpoints: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get cloud.ru endpoints: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get cloud.ru endpoints: unexpected status code %d", resp.StatusCode)
	}

	var endpoints EndpointsResponse
	if err = json.NewDecoder(resp.Body).Decode(&endpoints); err != nil {
		return nil, fmt.Errorf("decode cloud.ru endpoints: %w", err)
	}

	return &endpoints, nil
}

// Get returns the API address of the product by its ID.
// If the product is not found, the function returns nil.
func (er *EndpointsResponse) Get(id string) *Endpoint {
	for i := range er.Endpoints {
		if er.Endpoints[i].ID == id {
			return &er.Endpoints[i]
		}
	}

	return nil
}

// Close closes the KMS gRPC client connection.
func (c *Client) Close() error { return c.kmsConn.Close() }
