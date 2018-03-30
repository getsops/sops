// Copyright 2017, Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// AUTO-GENERATED CODE. DO NOT EDIT.

package logging

import (
	"math"
	"time"

	"cloud.google.com/go/internal/version"
	gax "github.com/googleapis/gax-go"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	loggingpb "google.golang.org/genproto/googleapis/logging/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ConfigCallOptions contains the retry settings for each method of ConfigClient.
type ConfigCallOptions struct {
	ListSinks  []gax.CallOption
	GetSink    []gax.CallOption
	CreateSink []gax.CallOption
	UpdateSink []gax.CallOption
	DeleteSink []gax.CallOption
}

func defaultConfigClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("logging.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultConfigCallOptions() *ConfigCallOptions {
	retry := map[[2]string][]gax.CallOption{
		{"default", "idempotent"}: {
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.DeadlineExceeded,
					codes.Internal,
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        1000 * time.Millisecond,
					Multiplier: 1.2,
				})
			}),
		},
	}
	return &ConfigCallOptions{
		ListSinks:  retry[[2]string{"default", "idempotent"}],
		GetSink:    retry[[2]string{"default", "idempotent"}],
		CreateSink: retry[[2]string{"default", "non_idempotent"}],
		UpdateSink: retry[[2]string{"default", "non_idempotent"}],
		DeleteSink: retry[[2]string{"default", "idempotent"}],
	}
}

// ConfigClient is a client for interacting with Stackdriver Logging API.
type ConfigClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	configClient loggingpb.ConfigServiceV2Client

	// The call options for this service.
	CallOptions *ConfigCallOptions

	// The metadata to be sent with each request.
	xGoogHeader []string
}

// NewConfigClient creates a new config service v2 client.
//
// Service for configuring sinks used to export log entries outside of
// Stackdriver Logging.
func NewConfigClient(ctx context.Context, opts ...option.ClientOption) (*ConfigClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultConfigClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &ConfigClient{
		conn:        conn,
		CallOptions: defaultConfigCallOptions(),

		configClient: loggingpb.NewConfigServiceV2Client(conn),
	}
	c.SetGoogleClientInfo()
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *ConfigClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *ConfigClient) Close() error {
	return c.conn.Close()
}

// SetGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *ConfigClient) SetGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", version.Go()}, keyval...)
	kv = append(kv, "gapic", version.Repo, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogHeader = []string{gax.XGoogHeader(kv...)}
}

// ConfigProjectPath returns the path for the project resource.
func ConfigProjectPath(project string) string {
	return "" +
		"projects/" +
		project +
		""
}

// ConfigSinkPath returns the path for the sink resource.
func ConfigSinkPath(project, sink string) string {
	return "" +
		"projects/" +
		project +
		"/sinks/" +
		sink +
		""
}

// ListSinks lists sinks.
func (c *ConfigClient) ListSinks(ctx context.Context, req *loggingpb.ListSinksRequest, opts ...gax.CallOption) *LogSinkIterator {
	ctx = insertXGoog(ctx, c.xGoogHeader)
	opts = append(c.CallOptions.ListSinks[0:len(c.CallOptions.ListSinks):len(c.CallOptions.ListSinks)], opts...)
	it := &LogSinkIterator{}
	it.InternalFetch = func(pageSize int, pageToken string) ([]*loggingpb.LogSink, string, error) {
		var resp *loggingpb.ListSinksResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.configClient.ListSinks(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Sinks, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	return it
}

// GetSink gets a sink.
func (c *ConfigClient) GetSink(ctx context.Context, req *loggingpb.GetSinkRequest, opts ...gax.CallOption) (*loggingpb.LogSink, error) {
	ctx = insertXGoog(ctx, c.xGoogHeader)
	opts = append(c.CallOptions.GetSink[0:len(c.CallOptions.GetSink):len(c.CallOptions.GetSink)], opts...)
	var resp *loggingpb.LogSink
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.configClient.GetSink(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateSink creates a sink that exports specified log entries to a destination.  The
// export of newly-ingested log entries begins immediately, unless the current
// time is outside the sink's start and end times or the sink's
// writer_identity is not permitted to write to the destination.  A sink can
// export log entries only from the resource owning the sink.
func (c *ConfigClient) CreateSink(ctx context.Context, req *loggingpb.CreateSinkRequest, opts ...gax.CallOption) (*loggingpb.LogSink, error) {
	ctx = insertXGoog(ctx, c.xGoogHeader)
	opts = append(c.CallOptions.CreateSink[0:len(c.CallOptions.CreateSink):len(c.CallOptions.CreateSink)], opts...)
	var resp *loggingpb.LogSink
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.configClient.CreateSink(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateSink updates a sink. If the named sink doesn't exist, then this method is
// identical to
// sinks.create (at /logging/docs/api/reference/rest/v2/projects.sinks/create).
// If the named sink does exist, then this method replaces the following
// fields in the existing sink with values from the new sink: destination,
// filter, output_version_format, start_time, and end_time.
// The updated filter might also have a new writer_identity; see the
// unique_writer_identity field.
func (c *ConfigClient) UpdateSink(ctx context.Context, req *loggingpb.UpdateSinkRequest, opts ...gax.CallOption) (*loggingpb.LogSink, error) {
	ctx = insertXGoog(ctx, c.xGoogHeader)
	opts = append(c.CallOptions.UpdateSink[0:len(c.CallOptions.UpdateSink):len(c.CallOptions.UpdateSink)], opts...)
	var resp *loggingpb.LogSink
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.configClient.UpdateSink(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteSink deletes a sink. If the sink has a unique writer_identity, then that
// service account is also deleted.
func (c *ConfigClient) DeleteSink(ctx context.Context, req *loggingpb.DeleteSinkRequest, opts ...gax.CallOption) error {
	ctx = insertXGoog(ctx, c.xGoogHeader)
	opts = append(c.CallOptions.DeleteSink[0:len(c.CallOptions.DeleteSink):len(c.CallOptions.DeleteSink)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.configClient.DeleteSink(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// LogSinkIterator manages a stream of *loggingpb.LogSink.
type LogSinkIterator struct {
	items    []*loggingpb.LogSink
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*loggingpb.LogSink, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *LogSinkIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *LogSinkIterator) Next() (*loggingpb.LogSink, error) {
	var item *loggingpb.LogSink
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *LogSinkIterator) bufLen() int {
	return len(it.items)
}

func (it *LogSinkIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}
