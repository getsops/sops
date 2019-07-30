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

package testutil

import (
	"strconv"

	structpb "github.com/golang/protobuf/ptypes/struct"
	spannerpb "google.golang.org/genproto/googleapis/spanner/v1"
	"google.golang.org/grpc/codes"
)

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	apiv1 "cloud.google.com/go/spanner/apiv1"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	gstatus "google.golang.org/grpc/status"
)

// clientOpt is the option tests should use to connect to the test server.
// It is initialized by TestMain.
var serverAddress string
var clientOpt option.ClientOption
var testSpanner InMemSpannerServer

// Mocked selectSQL statement.
const selectSQL = "SELECT FOO FROM BAR"
const selectRowCount int64 = 2
const selectColCount int = 1

var selectValues = [...]int64{1, 2}

// Mocked DML statement.
const updateSQL = "UPDATE FOO SET BAR=1 WHERE ID=ID"
const updateRowCount int64 = 2

func TestMain(m *testing.M) {
	flag.Parse()

	testSpanner = NewInMemSpannerServer()
	serv := grpc.NewServer()
	spannerpb.RegisterSpannerServer(serv, testSpanner)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	go serv.Serve(lis)

	serverAddress = lis.Addr().String()
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	clientOpt = option.WithGRPCConn(conn)

	os.Exit(m.Run())
}

// Resets the mock server to its default values and registers a mocked result
// for the statements "SELECT FOO FROM BAR" and
// "UPDATE FOO SET BAR=1 WHERE ID=ID".
func setup() {
	testSpanner.Reset()
	fields := make([]*spannerpb.StructType_Field, selectColCount)
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
	rows := make([]*structpb.ListValue, selectRowCount)
	for idx, value := range selectValues {
		rowValue := make([]*structpb.Value, selectColCount)
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
	result := &StatementResult{Type: StatementResultResultSet, ResultSet: resultSet}
	testSpanner.PutStatementResult(selectSQL, result)

	updateResult := &StatementResult{Type: StatementResultUpdateCount, UpdateCount: updateRowCount}
	testSpanner.PutStatementResult(updateSQL, updateResult)
}

func TestSpannerCreateSession(t *testing.T) {
	testSpanner.Reset()
	var expectedName = fmt.Sprintf("projects/%s/instances/%s/databases/%s/sessions/", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var request = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}

	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.CreateSession(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(resp.Name, expectedName) != 0 {
		t.Errorf("wrong name %s, should start with %s)", resp.Name, expectedName)
	}
}

func TestSpannerCreateSession_Unavailable(t *testing.T) {
	testSpanner.Reset()
	var expectedName = fmt.Sprintf("projects/%s/instances/%s/databases/%s/sessions/", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var request = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}

	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}
	testSpanner.SetError(gstatus.Error(codes.Unavailable, "Temporary unavailable"))
	resp, err := c.CreateSession(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Index(resp.Name, expectedName) != 0 {
		t.Errorf("wrong name %s, should start with %s)", resp.Name, expectedName)
	}
}

func TestSpannerGetSession(t *testing.T) {
	testSpanner.Reset()
	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}

	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}
	createResp, err := c.CreateSession(context.Background(), createRequest)
	if err != nil {
		t.Fatal(err)
	}
	var getRequest = &spannerpb.GetSessionRequest{
		Name: createResp.Name,
	}
	getResp, err := c.GetSession(context.Background(), getRequest)
	if err != nil {
		t.Fatal(err)
	}
	if getResp.Name != getRequest.Name {
		t.Errorf("wrong name %s, expected %s)", getResp.Name, getRequest.Name)
	}
}

func TestSpannerListSessions(t *testing.T) {
	testSpanner.Reset()
	const expectedNumberOfSessions = 5
	var expectedName = fmt.Sprintf("projects/%s/instances/%s/databases/%s/sessions/", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}

	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < expectedNumberOfSessions; i++ {
		_, err := c.CreateSession(context.Background(), createRequest)
		if err != nil {
			t.Fatal(err)
		}
	}
	var listRequest = &spannerpb.ListSessionsRequest{
		Database: formattedDatabase,
	}
	var sessionCount int
	listResp := c.ListSessions(context.Background(), listRequest)
	for {
		session, err := listResp.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if strings.Index(session.Name, expectedName) != 0 {
			t.Errorf("wrong name %s, should start with %s)", session.Name, expectedName)
		}
		sessionCount++
	}
	if sessionCount != expectedNumberOfSessions {
		t.Errorf("wrong number of sessions: %d, expected %d", sessionCount, expectedNumberOfSessions)
	}
}

func TestSpannerDeleteSession(t *testing.T) {
	testSpanner.Reset()
	const expectedNumberOfSessions = 5
	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}

	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < expectedNumberOfSessions; i++ {
		_, err := c.CreateSession(context.Background(), createRequest)
		if err != nil {
			t.Fatal(err)
		}
	}
	var listRequest = &spannerpb.ListSessionsRequest{
		Database: formattedDatabase,
	}
	var sessionCount int
	listResp := c.ListSessions(context.Background(), listRequest)
	for {
		session, err := listResp.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		var deleteRequest = &spannerpb.DeleteSessionRequest{
			Name: session.Name,
		}
		c.DeleteSession(context.Background(), deleteRequest)
		sessionCount++
	}
	if sessionCount != expectedNumberOfSessions {
		t.Errorf("wrong number of sessions: %d, expected %d", sessionCount, expectedNumberOfSessions)
	}
	// Re-list all sessions. This should now be empty.
	listResp = c.ListSessions(context.Background(), listRequest)
	_, err = listResp.Next()
	if err != iterator.Done {
		t.Errorf("expected empty session iterator")
	}
}

func TestSpannerExecuteSql(t *testing.T) {
	setup()
	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}
	session, err := c.CreateSession(context.Background(), createRequest)
	if err != nil {
		t.Fatal(err)
	}
	request := &spannerpb.ExecuteSqlRequest{
		Session: session.Name,
		Sql:     selectSQL,
		Transaction: &spannerpb.TransactionSelector{
			Selector: &spannerpb.TransactionSelector_SingleUse{
				SingleUse: &spannerpb.TransactionOptions{
					Mode: &spannerpb.TransactionOptions_ReadOnly_{
						ReadOnly: &spannerpb.TransactionOptions_ReadOnly{
							ReturnReadTimestamp: false,
							TimestampBound: &spannerpb.TransactionOptions_ReadOnly_Strong{
								Strong: true,
							},
						},
					},
				},
			},
		},
		Seqno:     1,
		QueryMode: spannerpb.ExecuteSqlRequest_NORMAL,
	}
	response, err := c.ExecuteSql(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	var rowCount int64
	for _, row := range response.Rows {
		if len(row.Values) != selectColCount {
			t.Fatalf("unexpected number of columns: %d, expected %d", len(row.Values), selectColCount)
		}
		rowCount++
	}
	if rowCount != selectRowCount {
		t.Fatalf("unexpected number of rows: %d, expected %d", rowCount, selectRowCount)
	}
}

func TestSpannerExecuteSqlDml(t *testing.T) {
	setup()
	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}
	session, err := c.CreateSession(context.Background(), createRequest)
	if err != nil {
		t.Fatal(err)
	}
	request := &spannerpb.ExecuteSqlRequest{
		Session: session.Name,
		Sql:     updateSQL,
		Transaction: &spannerpb.TransactionSelector{
			Selector: &spannerpb.TransactionSelector_Begin{
				Begin: &spannerpb.TransactionOptions{
					Mode: &spannerpb.TransactionOptions_ReadWrite_{
						ReadWrite: &spannerpb.TransactionOptions_ReadWrite{},
					},
				},
			},
		},
		Seqno:     1,
		QueryMode: spannerpb.ExecuteSqlRequest_NORMAL,
	}
	response, err := c.ExecuteSql(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	var rowCount int64 = response.Stats.GetRowCountExact()
	if rowCount != updateRowCount {
		t.Fatalf("unexpected number of rows updated: %d, expected %d", rowCount, updateRowCount)
	}
}

func TestSpannerExecuteStreamingSql(t *testing.T) {
	setup()
	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}
	session, err := c.CreateSession(context.Background(), createRequest)
	if err != nil {
		t.Fatal(err)
	}
	request := &spannerpb.ExecuteSqlRequest{
		Session: session.Name,
		Sql:     selectSQL,
		Transaction: &spannerpb.TransactionSelector{
			Selector: &spannerpb.TransactionSelector_SingleUse{
				SingleUse: &spannerpb.TransactionOptions{
					Mode: &spannerpb.TransactionOptions_ReadOnly_{
						ReadOnly: &spannerpb.TransactionOptions_ReadOnly{
							ReturnReadTimestamp: false,
							TimestampBound: &spannerpb.TransactionOptions_ReadOnly_Strong{
								Strong: true,
							},
						},
					},
				},
			},
		},
		Seqno:     1,
		QueryMode: spannerpb.ExecuteSqlRequest_NORMAL,
	}
	response, err := c.ExecuteStreamingSql(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	partial, err := response.Recv()
	if err != nil {
		t.Fatal(err)
	}
	var rowIndex int64
	colCount := len(partial.Metadata.RowType.Fields)
	if colCount != selectColCount {
		t.Fatalf("unexpected number of columns: %d, expected %d", colCount, selectColCount)
	}
	for {
		for col := 0; col < colCount; col++ {
			val, err := strconv.ParseInt(partial.Values[rowIndex*int64(colCount)+int64(col)].GetStringValue(), 10, 64)
			if err != nil {
				t.Fatal(err)
			}
			if val != selectValues[rowIndex] {
				t.Fatalf("Unexpected value at index %d. Expected %d, got %d", rowIndex, selectValues[rowIndex], val)
			}
		}
		rowIndex++
		if rowIndex == selectRowCount {
			break
		}
	}
	if rowIndex != selectRowCount {
		t.Fatalf("unexpected number of rows: %d, expected %d", rowIndex, selectRowCount)
	}
}

func TestSpannerExecuteBatchDml(t *testing.T) {
	setup()
	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}
	session, err := c.CreateSession(context.Background(), createRequest)
	if err != nil {
		t.Fatal(err)
	}
	statements := make([]*spannerpb.ExecuteBatchDmlRequest_Statement, 3)
	for idx := 0; idx < len(statements); idx++ {
		statements[idx] = &spannerpb.ExecuteBatchDmlRequest_Statement{Sql: updateSQL}
	}
	executeBatchDmlRequest := &spannerpb.ExecuteBatchDmlRequest{
		Session:    session.Name,
		Statements: statements,
		Transaction: &spannerpb.TransactionSelector{
			Selector: &spannerpb.TransactionSelector_Begin{
				Begin: &spannerpb.TransactionOptions{
					Mode: &spannerpb.TransactionOptions_ReadWrite_{
						ReadWrite: &spannerpb.TransactionOptions_ReadWrite{},
					},
				},
			},
		},
		Seqno: 1,
	}
	response, err := c.ExecuteBatchDml(context.Background(), executeBatchDmlRequest)
	if err != nil {
		t.Fatal(err)
	}
	var totalRowCount int64
	for _, res := range response.ResultSets {
		var rowCount int64 = res.Stats.GetRowCountExact()
		if rowCount != updateRowCount {
			t.Fatalf("unexpected number of rows updated: %d, expected %d", rowCount, updateRowCount)
		}
		totalRowCount += rowCount
	}
	if totalRowCount != updateRowCount*int64(len(statements)) {
		t.Fatalf("unexpected number of total rows updated: %d, expected %d", totalRowCount, updateRowCount*int64(len(statements)))
	}
}

func TestBeginTransaction(t *testing.T) {
	setup()
	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}
	session, err := c.CreateSession(context.Background(), createRequest)
	if err != nil {
		t.Fatal(err)
	}
	beginRequest := &spannerpb.BeginTransactionRequest{
		Session: session.Name,
		Options: &spannerpb.TransactionOptions{
			Mode: &spannerpb.TransactionOptions_ReadWrite_{
				ReadWrite: &spannerpb.TransactionOptions_ReadWrite{},
			},
		},
	}
	tx, err := c.BeginTransaction(context.Background(), beginRequest)
	if err != nil {
		t.Fatal(err)
	}
	expectedName := fmt.Sprintf("%s/transactions/", session.Name)
	if strings.Index(string(tx.Id), expectedName) != 0 {
		t.Errorf("wrong name %s, should start with %s)", string(tx.Id), expectedName)
	}
}

func TestCommitTransaction(t *testing.T) {
	setup()
	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}
	session, err := c.CreateSession(context.Background(), createRequest)
	if err != nil {
		t.Fatal(err)
	}
	beginRequest := &spannerpb.BeginTransactionRequest{
		Session: session.Name,
		Options: &spannerpb.TransactionOptions{
			Mode: &spannerpb.TransactionOptions_ReadWrite_{
				ReadWrite: &spannerpb.TransactionOptions_ReadWrite{},
			},
		},
	}
	tx, err := c.BeginTransaction(context.Background(), beginRequest)
	if err != nil {
		t.Fatal(err)
	}
	commitRequest := &spannerpb.CommitRequest{
		Session: session.Name,
		Transaction: &spannerpb.CommitRequest_TransactionId{
			TransactionId: tx.Id,
		},
	}
	resp, err := c.Commit(context.Background(), commitRequest)
	if err != nil {
		t.Fatal(err)
	}
	if resp.CommitTimestamp == nil {
		t.Fatalf("No commit timestamp returned")
	}
}

func TestRollbackTransaction(t *testing.T) {
	setup()
	c, err := apiv1.NewClient(context.Background(), clientOpt)
	if err != nil {
		t.Fatal(err)
	}

	var formattedDatabase = fmt.Sprintf("projects/%s/instances/%s/databases/%s", "[PROJECT]", "[INSTANCE]", "[DATABASE]")
	var createRequest = &spannerpb.CreateSessionRequest{
		Database: formattedDatabase,
	}
	session, err := c.CreateSession(context.Background(), createRequest)
	if err != nil {
		t.Fatal(err)
	}
	beginRequest := &spannerpb.BeginTransactionRequest{
		Session: session.Name,
		Options: &spannerpb.TransactionOptions{
			Mode: &spannerpb.TransactionOptions_ReadWrite_{
				ReadWrite: &spannerpb.TransactionOptions_ReadWrite{},
			},
		},
	}
	tx, err := c.BeginTransaction(context.Background(), beginRequest)
	if err != nil {
		t.Fatal(err)
	}
	rollbackRequest := &spannerpb.RollbackRequest{
		Session:       session.Name,
		TransactionId: tx.Id,
	}
	err = c.Rollback(context.Background(), rollbackRequest)
	if err != nil {
		t.Fatal(err)
	}
}
