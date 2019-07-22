// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spanner

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	"cloud.google.com/go/spanner/internal/testutil"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/api/option"
	spannerpb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc"
)

// The SQL statements and results that are already mocked for this test server.
const selectFooFromBar = "SELECT FOO FROM BAR"
const selectFooFromBarRowCount int64 = 2
const selectFooFromBarColCount int = 1

var selectFooFromBarResults = [...]int64{1, 2}

const selectSingerIDAlbumIDAlbumTitleFromAlbums = "SELECT SingerId, AlbumId, AlbumTitle FROM Albums"
const selectSingerIDAlbumIDAlbumTitleFromAlbumsRowCount int64 = 3
const selectSingerIDAlbumIDAlbumTitleFromAlbumsColCount int = 3

const updateBarSetFoo = "UPDATE FOO SET BAR=1 WHERE BAZ=2"
const updateBarSetFooRowCount = 5

// An InMemSpannerServer with results for a number of SQL statements readily
// mocked.
type spannerInMemTestServer struct {
	testSpanner testutil.InMemSpannerServer
	server      *grpc.Server
}

// Create a spannerInMemTestServer with default configuration.
func newSpannerInMemTestServer(t *testing.T) (*spannerInMemTestServer, *Client) {
	s := &spannerInMemTestServer{}
	client := s.setup(t)
	return s, client
}

// Create a spannerInMemTestServer with the specified configuration.
func newSpannerInMemTestServerWithConfig(t *testing.T, config ClientConfig) (*spannerInMemTestServer, *Client) {
	s := &spannerInMemTestServer{}
	client := s.setupWithConfig(t, config)
	return s, client
}

func (s *spannerInMemTestServer) setup(t *testing.T) *Client {
	return s.setupWithConfig(t, ClientConfig{})
}

func (s *spannerInMemTestServer) setupWithConfig(t *testing.T, config ClientConfig) *Client {
	s.testSpanner = testutil.NewInMemSpannerServer()
	s.setupFooResults()
	s.setupSingersResults()
	s.server = grpc.NewServer()
	spannerpb.RegisterSpannerServer(s.server, s.testSpanner)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	go s.server.Serve(lis)

	serverAddress := lis.Addr().String()
	ctx := context.Background()
	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	client, err := NewClientWithConfig(ctx, formattedDatabase, config,
		option.WithEndpoint(serverAddress),
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithoutAuthentication(),
	)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func (s *spannerInMemTestServer) setupFooResults() {
	fields := make([]*spannerpb.StructType_Field, selectFooFromBarColCount)
	fields[0] = &spannerpb.StructType_Field{
		Name: "FOO",
		Type: &spannerpb.Type{Code: spannerpb.TypeCode_INT64},
	}
	rowType := &spannerpb.StructType{
		Fields: fields,
	}
	metadata := &spannerpb.ResultSetMetadata{
		RowType: rowType,
	}
	rows := make([]*structpb.ListValue, selectFooFromBarRowCount)
	for idx, value := range selectFooFromBarResults {
		rowValue := make([]*structpb.Value, selectFooFromBarColCount)
		rowValue[0] = &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: strconv.FormatInt(value, 10)},
		}
		rows[idx] = &structpb.ListValue{
			Values: rowValue,
		}
	}
	resultSet := &spannerpb.ResultSet{
		Metadata: metadata,
		Rows:     rows,
	}
	result := &testutil.StatementResult{Type: testutil.StatementResultResultSet, ResultSet: resultSet}
	s.testSpanner.PutStatementResult(selectFooFromBar, result)
	s.testSpanner.PutStatementResult(updateBarSetFoo, &testutil.StatementResult{
		Type:        testutil.StatementResultUpdateCount,
		UpdateCount: updateBarSetFooRowCount,
	})
}

func (s *spannerInMemTestServer) setupSingersResults() {
	fields := make([]*spannerpb.StructType_Field, selectSingerIDAlbumIDAlbumTitleFromAlbumsColCount)
	fields[0] = &spannerpb.StructType_Field{
		Name: "SingerId",
		Type: &spannerpb.Type{Code: spannerpb.TypeCode_INT64},
	}
	fields[1] = &spannerpb.StructType_Field{
		Name: "AlbumId",
		Type: &spannerpb.Type{Code: spannerpb.TypeCode_INT64},
	}
	fields[2] = &spannerpb.StructType_Field{
		Name: "AlbumTitle",
		Type: &spannerpb.Type{Code: spannerpb.TypeCode_STRING},
	}
	rowType := &spannerpb.StructType{
		Fields: fields,
	}
	metadata := &spannerpb.ResultSetMetadata{
		RowType: rowType,
	}
	rows := make([]*structpb.ListValue, selectSingerIDAlbumIDAlbumTitleFromAlbumsRowCount)
	var idx int64
	for idx = 0; idx < selectSingerIDAlbumIDAlbumTitleFromAlbumsRowCount; idx++ {
		rowValue := make([]*structpb.Value, selectSingerIDAlbumIDAlbumTitleFromAlbumsColCount)
		rowValue[0] = &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: strconv.FormatInt(idx+1, 10)},
		}
		rowValue[1] = &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: strconv.FormatInt(idx*10+idx, 10)},
		}
		rowValue[2] = &structpb.Value{
			Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("Album title %d", idx)},
		}
		rows[idx] = &structpb.ListValue{
			Values: rowValue,
		}
	}
	resultSet := &spannerpb.ResultSet{
		Metadata: metadata,
		Rows:     rows,
	}
	result := &testutil.StatementResult{Type: testutil.StatementResultResultSet, ResultSet: resultSet}
	s.testSpanner.PutStatementResult(selectSingerIDAlbumIDAlbumTitleFromAlbums, result)
}

func (s *spannerInMemTestServer) teardown(client *Client) {
	client.Close()
	s.server.Stop()
}
