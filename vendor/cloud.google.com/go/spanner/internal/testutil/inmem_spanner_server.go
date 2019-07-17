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
	emptypb "github.com/golang/protobuf/ptypes/empty"
	structpb "github.com/golang/protobuf/ptypes/struct"
	spannerpb "google.golang.org/genproto/googleapis/spanner/v1"
)

import (
	"context"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

// StatementResultType indicates the type of result returned by a SQL
// statement.
type StatementResultType int

const (
	// StatementResultError indicates that the sql statement returns an error.
	StatementResultError StatementResultType = 0
	// StatementResultResultSet indicates that the sql statement returns a
	// result set.
	StatementResultResultSet StatementResultType = 1
	// StatementResultUpdateCount indicates that the sql statement returns an
	// update count.
	StatementResultUpdateCount StatementResultType = 2
)

// The method names that can be used to register execution times and errors.
const (
	MethodBeginTransaction    string = "BEGIN_TRANSACTION"
	MethodCommitTransaction   string = "COMMIT_TRANSACTION"
	MethodCreateSession       string = "CREATE_SESSION"
	MethodDeleteSession       string = "DELETE_SESSION"
	MethodGetSession          string = "GET_SESSION"
	MethodExecuteStreamingSql string = "EXECUTE_STREAMING_SQL"
)

// StatementResult represents a mocked result on the test server. Th result can
// be either a ResultSet, an update count or an error.
type StatementResult struct {
	Type        StatementResultType
	Err         error
	ResultSet   *spannerpb.ResultSet
	UpdateCount int64
}

// Converts a ResultSet to a PartialResultSet. This method is used to convert
// a mocked result to a PartialResultSet when one of the streaming methods are
// called.
func (s *StatementResult) toPartialResultSet() *spannerpb.PartialResultSet {
	values := make([]*structpb.Value,
		len(s.ResultSet.Rows)*len(s.ResultSet.Metadata.RowType.Fields))
	var idx int
	for _, row := range s.ResultSet.Rows {
		for colIdx := range s.ResultSet.Metadata.RowType.Fields {
			values[idx] = row.Values[colIdx]
			idx++
		}
	}
	return &spannerpb.PartialResultSet{
		Metadata: s.ResultSet.Metadata,
		Values:   values,
	}
}

func (s *StatementResult) updateCountToPartialResultSet(exact bool) *spannerpb.PartialResultSet {
	return &spannerpb.PartialResultSet{
		Stats: s.convertUpdateCountToResultSet(exact).Stats,
	}
}

// Converts an update count to a ResultSet, as DML statements also return the
// update count as the statistics of a ResultSet.
func (s *StatementResult) convertUpdateCountToResultSet(exact bool) *spannerpb.ResultSet {
	if exact {
		return &spannerpb.ResultSet{
			Stats: &spannerpb.ResultSetStats{
				RowCount: &spannerpb.ResultSetStats_RowCountExact{
					RowCountExact: s.UpdateCount,
				},
			},
		}
	}
	return &spannerpb.ResultSet{
		Stats: &spannerpb.ResultSetStats{
			RowCount: &spannerpb.ResultSetStats_RowCountLowerBound{
				RowCountLowerBound: s.UpdateCount,
			},
		},
	}
}

// SimulatedExecutionTime represents the time the execution of a method
// should take, and any errors that should be returned by the method.
type SimulatedExecutionTime struct {
	MinimumExecutionTime time.Duration
	RandomExecutionTime  time.Duration
	Errors               []error
	// Keep error after execution. The error will continue to be returned until
	// it is cleared.
	KeepError bool
}

// InMemSpannerServer contains the SpannerServer interface plus a couple
// of specific methods for adding mocked results and resetting the server.
type InMemSpannerServer interface {
	spannerpb.SpannerServer

	// Stops this server.
	Stop()

	// Resets the in-mem server to its default state, deleting all sessions and
	// transactions that have been created on the server. Mocked results are
	// not deleted.
	Reset()

	// Sets an error that will be returned by the next server call. The server
	// call will also automatically clear the error.
	SetError(err error)

	// Puts a mocked result on the server for a specific sql statement. The
	// server does not parse the SQL string in any way, it is merely used as
	// a key to the mocked result. The result will be used for all methods that
	// expect a SQL statement, including (batch) DML methods.
	PutStatementResult(sql string, result *StatementResult) error

	// Removes a mocked result on the server for a specific sql statement.
	RemoveStatementResult(sql string)

	// Aborts the specified transaction . This method can be used to test
	// transaction retry logic.
	AbortTransaction(id []byte)

	// Puts a simulated execution time for one of the Spanner methods.
	PutExecutionTime(method string, executionTime SimulatedExecutionTime)
	// Freeze stalls all requests.
	Freeze()
	// Unfreeze restores processing requests.
	Unfreeze()

	TotalSessionsCreated() uint
	TotalSessionsDeleted() uint

	ReceivedRequests() chan interface{}
	DumpSessions() map[string]bool
	ClearPings()
	DumpPings() []string
}

type inMemSpannerServer struct {
	// Embed for forward compatibility.
	// Tests will keep working if more methods are added
	// in the future.
	spannerpb.SpannerServer

	mu sync.Mutex

	// If set, all calls return this error.
	err error
	// The mock server creates session IDs using this counter.
	sessionCounter uint64
	// The sessions that have been created on this mock server.
	sessions map[string]*spannerpb.Session
	// Last use times per session.
	sessionLastUseTime map[string]time.Time

	// The mock server creates transaction IDs per session using these
	// counters.
	transactionCounters map[string]*uint64
	// The transactions that have been created on this mock server.
	transactions map[string]*spannerpb.Transaction
	// The transactions that have been (manually) aborted on the server.
	abortedTransactions map[string]bool
	// The transactions that are marked as PartitionedDMLTransaction
	partitionedDmlTransactions map[string]bool

	// The mocked results for this server.
	statementResults map[string]*StatementResult
	// The simulated execution times per method.
	executionTimes map[string]*SimulatedExecutionTime
	// Server will stall on any requests.
	freezed chan struct{}

	totalSessionsCreated uint
	totalSessionsDeleted uint
	receivedRequests     chan interface{}
	// Session ping history.
	pings []string
}

// NewInMemSpannerServer creates a new in-mem test server.
func NewInMemSpannerServer() InMemSpannerServer {
	res := &inMemSpannerServer{}
	res.initDefaults()
	res.statementResults = make(map[string]*StatementResult)
	res.executionTimes = make(map[string]*SimulatedExecutionTime)
	res.receivedRequests = make(chan interface{}, 1000000)
	// Produce a closed channel, so the default action of ready is to not block.
	res.Freeze()
	res.Unfreeze()
	return res
}

func (s *inMemSpannerServer) Stop() {
	close(s.receivedRequests)
}

// Resets the test server to its initial state, deleting all sessions and
// transactions that have been created on the server. This method will not
// remove mocked results.
func (s *inMemSpannerServer) Reset() {
	close(s.receivedRequests)
	s.receivedRequests = make(chan interface{}, 1000000)
	s.initDefaults()
}

func (s *inMemSpannerServer) SetError(err error) {
	s.err = err
}

// Registers a mocked result for a SQL statement on the server.
func (s *inMemSpannerServer) PutStatementResult(sql string, result *StatementResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statementResults[sql] = result
	return nil
}

func (s *inMemSpannerServer) RemoveStatementResult(sql string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.statementResults, sql)
}

func (s *inMemSpannerServer) AbortTransaction(id []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.abortedTransactions[string(id)] = true
}

func (s *inMemSpannerServer) PutExecutionTime(method string, executionTime SimulatedExecutionTime) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executionTimes[method] = &executionTime
}

// Freeze stalls all requests.
func (s *inMemSpannerServer) Freeze() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.freezed = make(chan struct{})
}

// Unfreeze restores processing requests.
func (s *inMemSpannerServer) Unfreeze() {
	s.mu.Lock()
	defer s.mu.Unlock()
	close(s.freezed)
}

// ready checks conditions before executing requests
func (s *inMemSpannerServer) ready() {
	s.mu.Lock()
	freezed := s.freezed
	s.mu.Unlock()
	// check if server should be freezed
	<-freezed
}

func (s *inMemSpannerServer) TotalSessionsCreated() uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.totalSessionsCreated
}

func (s *inMemSpannerServer) TotalSessionsDeleted() uint {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.totalSessionsDeleted
}

func (s *inMemSpannerServer) ReceivedRequests() chan interface{} {
	return s.receivedRequests
}

// ClearPings clears the ping history from the server.
func (s *inMemSpannerServer) ClearPings() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pings = nil
}

// DumpPings dumps the ping history.
func (s *inMemSpannerServer) DumpPings() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.pings...)
}

// DumpSessions dumps the internal session table.
func (s *inMemSpannerServer) DumpSessions() map[string]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	st := map[string]bool{}
	for s := range s.sessions {
		st[s] = true
	}
	return st
}

func (s *inMemSpannerServer) initDefaults() {
	s.sessionCounter = 0
	s.sessions = make(map[string]*spannerpb.Session)
	s.sessionLastUseTime = make(map[string]time.Time)
	s.transactions = make(map[string]*spannerpb.Transaction)
	s.abortedTransactions = make(map[string]bool)
	s.partitionedDmlTransactions = make(map[string]bool)
	s.transactionCounters = make(map[string]*uint64)
}

func (s *inMemSpannerServer) generateSessionName(database string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessionCounter++
	return fmt.Sprintf("%s/sessions/%d", database, s.sessionCounter)
}

func (s *inMemSpannerServer) findSession(name string) (*spannerpb.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session := s.sessions[name]
	if session == nil {
		return nil, gstatus.Error(codes.NotFound, fmt.Sprintf("Session %s not found", name))
	}
	return session, nil
}

func (s *inMemSpannerServer) updateSessionLastUseTime(session string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessionLastUseTime[session] = time.Now()
}

func getCurrentTimestamp() *timestamp.Timestamp {
	t := time.Now()
	return &timestamp.Timestamp{Seconds: t.Unix(), Nanos: int32(t.Nanosecond())}
}

// Gets the transaction id from the transaction selector. If the selector
// specifies that a new transaction should be started, this method will start
// a new transaction and return the id of that transaction.
func (s *inMemSpannerServer) getTransactionID(session *spannerpb.Session, txSelector *spannerpb.TransactionSelector) []byte {
	var res []byte
	if txSelector.GetBegin() != nil {
		// Start a new transaction.
		res = s.beginTransaction(session, txSelector.GetBegin()).Id
	} else if txSelector.GetId() != nil {
		res = txSelector.GetId()
	}
	return res
}

func (s *inMemSpannerServer) generateTransactionName(session string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	counter, ok := s.transactionCounters[session]
	if !ok {
		counter = new(uint64)
		s.transactionCounters[session] = counter
	}
	*counter++
	return fmt.Sprintf("%s/transactions/%d", session, *counter)
}

func (s *inMemSpannerServer) beginTransaction(session *spannerpb.Session, options *spannerpb.TransactionOptions) *spannerpb.Transaction {
	id := s.generateTransactionName(session.Name)
	res := &spannerpb.Transaction{
		Id:            []byte(id),
		ReadTimestamp: getCurrentTimestamp(),
	}
	s.mu.Lock()
	s.transactions[id] = res
	s.partitionedDmlTransactions[id] = options.GetPartitionedDml() != nil
	s.mu.Unlock()
	return res
}

func (s *inMemSpannerServer) getTransactionByID(id []byte) (*spannerpb.Transaction, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	tx, ok := s.transactions[string(id)]
	if !ok {
		return nil, gstatus.Error(codes.NotFound, "Transaction not found")
	}
	aborted, ok := s.abortedTransactions[string(id)]
	if ok && aborted {
		return nil, gstatus.Error(codes.Aborted, "Transaction has been aborted")
	}
	return tx, nil
}

func (s *inMemSpannerServer) removeTransaction(tx *spannerpb.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.transactions, string(tx.Id))
	delete(s.partitionedDmlTransactions, string(tx.Id))
}

func (s *inMemSpannerServer) getStatementResult(sql string) (*StatementResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	result, ok := s.statementResults[sql]
	if !ok {
		return nil, gstatus.Error(codes.Internal, fmt.Sprintf("No result found for statement %v", sql))
	}
	return result, nil
}

func (s *inMemSpannerServer) simulateExecutionTime(method string, req interface{}) error {
	s.receivedRequests <- req
	s.ready()
	s.mu.Lock()
	if s.err != nil {
		err := s.err
		s.err = nil
		s.mu.Unlock()
		return err
	}
	executionTime, ok := s.executionTimes[method]
	s.mu.Unlock()
	if ok {
		var randTime int64
		if executionTime.RandomExecutionTime > 0 {
			randTime = rand.Int63n(int64(executionTime.RandomExecutionTime))
		}
		totalExecutionTime := time.Duration(int64(executionTime.MinimumExecutionTime) + randTime)
		<-time.After(totalExecutionTime)
		if executionTime.Errors != nil && len(executionTime.Errors) > 0 {
			err := executionTime.Errors[0]
			if !executionTime.KeepError {
				executionTime.Errors = executionTime.Errors[1:]
			}
			return err
		}
	}
	return nil
}

func (s *inMemSpannerServer) CreateSession(ctx context.Context, req *spannerpb.CreateSessionRequest) (*spannerpb.Session, error) {
	if err := s.simulateExecutionTime(MethodCreateSession, req); err != nil {
		return nil, err
	}
	if req.Database == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing database")
	}
	sessionName := s.generateSessionName(req.Database)
	ts := getCurrentTimestamp()
	session := &spannerpb.Session{Name: sessionName, CreateTime: ts, ApproximateLastUseTime: ts}
	s.mu.Lock()
	s.totalSessionsCreated++
	s.sessions[sessionName] = session
	s.mu.Unlock()
	return session, nil
}

func (s *inMemSpannerServer) GetSession(ctx context.Context, req *spannerpb.GetSessionRequest) (*spannerpb.Session, error) {
	if err := s.simulateExecutionTime(MethodGetSession, req); err != nil {
		return nil, err
	}
	s.mu.Lock()
	s.pings = append(s.pings, req.Name)
	s.mu.Unlock()
	if req.Name == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Name)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *inMemSpannerServer) ListSessions(ctx context.Context, req *spannerpb.ListSessionsRequest) (*spannerpb.ListSessionsResponse, error) {
	s.receivedRequests <- req
	if req.Database == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing database")
	}
	expectedSessionName := req.Database + "/sessions/"
	var sessions []*spannerpb.Session
	s.mu.Lock()
	for _, session := range s.sessions {
		if strings.Index(session.Name, expectedSessionName) == 0 {
			sessions = append(sessions, session)
		}
	}
	s.mu.Unlock()
	sort.Slice(sessions[:], func(i, j int) bool {
		return sessions[i].Name < sessions[j].Name
	})
	res := &spannerpb.ListSessionsResponse{Sessions: sessions}
	return res, nil
}

func (s *inMemSpannerServer) DeleteSession(ctx context.Context, req *spannerpb.DeleteSessionRequest) (*emptypb.Empty, error) {
	if err := s.simulateExecutionTime(MethodDeleteSession, req); err != nil {
		return nil, err
	}
	if req.Name == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	if _, err := s.findSession(req.Name); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.totalSessionsDeleted++
	delete(s.sessions, req.Name)
	return &emptypb.Empty{}, nil
}

func (s *inMemSpannerServer) ExecuteSql(ctx context.Context, req *spannerpb.ExecuteSqlRequest) (*spannerpb.ResultSet, error) {
	s.receivedRequests <- req
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	var id []byte
	s.updateSessionLastUseTime(session.Name)
	if id = s.getTransactionID(session, req.Transaction); id != nil {
		_, err = s.getTransactionByID(id)
		if err != nil {
			return nil, err
		}
	}
	statementResult, err := s.getStatementResult(req.Sql)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	isPartitionedDml := s.partitionedDmlTransactions[string(id)]
	s.mu.Unlock()
	switch statementResult.Type {
	case StatementResultError:
		return nil, statementResult.Err
	case StatementResultResultSet:
		return statementResult.ResultSet, nil
	case StatementResultUpdateCount:
		return statementResult.convertUpdateCountToResultSet(!isPartitionedDml), nil
	}
	return nil, gstatus.Error(codes.Internal, "Unknown result type")
}

func (s *inMemSpannerServer) ExecuteStreamingSql(req *spannerpb.ExecuteSqlRequest, stream spannerpb.Spanner_ExecuteStreamingSqlServer) error {
	if err := s.simulateExecutionTime(MethodExecuteStreamingSql, req); err != nil {
		return err
	}
	if req.Session == "" {
		return gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return err
	}
	s.updateSessionLastUseTime(session.Name)
	var id []byte
	if id = s.getTransactionID(session, req.Transaction); id != nil {
		_, err = s.getTransactionByID(id)
		if err != nil {
			return err
		}
	}
	statementResult, err := s.getStatementResult(req.Sql)
	if err != nil {
		return err
	}
	s.mu.Lock()
	isPartitionedDml := s.partitionedDmlTransactions[string(id)]
	s.mu.Unlock()
	switch statementResult.Type {
	case StatementResultError:
		return statementResult.Err
	case StatementResultResultSet:
		part := statementResult.toPartialResultSet()
		if err := stream.Send(part); err != nil {
			return err
		}
		return nil
	case StatementResultUpdateCount:
		part := statementResult.updateCountToPartialResultSet(!isPartitionedDml)
		if err := stream.Send(part); err != nil {
			return err
		}
		return nil
	}
	return gstatus.Error(codes.Internal, "Unknown result type")
}

func (s *inMemSpannerServer) ExecuteBatchDml(ctx context.Context, req *spannerpb.ExecuteBatchDmlRequest) (*spannerpb.ExecuteBatchDmlResponse, error) {
	s.receivedRequests <- req
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	s.updateSessionLastUseTime(session.Name)
	var id []byte
	if id = s.getTransactionID(session, req.Transaction); id != nil {
		_, err = s.getTransactionByID(id)
		if err != nil {
			return nil, err
		}
	}
	s.mu.Lock()
	isPartitionedDml := s.partitionedDmlTransactions[string(id)]
	s.mu.Unlock()
	resp := &spannerpb.ExecuteBatchDmlResponse{}
	resp.ResultSets = make([]*spannerpb.ResultSet, len(req.Statements))
	for idx, batchStatement := range req.Statements {
		statementResult, err := s.getStatementResult(batchStatement.Sql)
		if err != nil {
			return nil, err
		}
		switch statementResult.Type {
		case StatementResultError:
			resp.Status = &status.Status{Code: int32(codes.Unknown)}
		case StatementResultResultSet:
			return nil, gstatus.Error(codes.InvalidArgument, fmt.Sprintf("Not an update statement: %v", batchStatement.Sql))
		case StatementResultUpdateCount:
			resp.ResultSets[idx] = statementResult.convertUpdateCountToResultSet(!isPartitionedDml)
			resp.Status = &status.Status{Code: int32(codes.OK)}
		}
	}
	return resp, nil
}

func (s *inMemSpannerServer) Read(ctx context.Context, req *spannerpb.ReadRequest) (*spannerpb.ResultSet, error) {
	s.receivedRequests <- req
	return nil, gstatus.Error(codes.Unimplemented, "Method not yet implemented")
}

func (s *inMemSpannerServer) StreamingRead(req *spannerpb.ReadRequest, stream spannerpb.Spanner_StreamingReadServer) error {
	s.receivedRequests <- req
	return gstatus.Error(codes.Unimplemented, "Method not yet implemented")
}

func (s *inMemSpannerServer) BeginTransaction(ctx context.Context, req *spannerpb.BeginTransactionRequest) (*spannerpb.Transaction, error) {
	if err := s.simulateExecutionTime(MethodBeginTransaction, req); err != nil {
		return nil, err
	}
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	s.updateSessionLastUseTime(session.Name)
	tx := s.beginTransaction(session, req.Options)
	return tx, nil
}

func (s *inMemSpannerServer) Commit(ctx context.Context, req *spannerpb.CommitRequest) (*spannerpb.CommitResponse, error) {
	if err := s.simulateExecutionTime(MethodCommitTransaction, req); err != nil {
		return nil, err
	}
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	s.updateSessionLastUseTime(session.Name)
	var tx *spannerpb.Transaction
	if req.GetSingleUseTransaction() != nil {
		tx = s.beginTransaction(session, req.GetSingleUseTransaction())
	} else if req.GetTransactionId() != nil {
		tx, err = s.getTransactionByID(req.GetTransactionId())
		if err != nil {
			return nil, err
		}
	} else {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing transaction in commit request")
	}
	s.removeTransaction(tx)
	return &spannerpb.CommitResponse{CommitTimestamp: getCurrentTimestamp()}, nil
}

func (s *inMemSpannerServer) Rollback(ctx context.Context, req *spannerpb.RollbackRequest) (*emptypb.Empty, error) {
	s.receivedRequests <- req
	if req.Session == "" {
		return nil, gstatus.Error(codes.InvalidArgument, "Missing session name")
	}
	session, err := s.findSession(req.Session)
	if err != nil {
		return nil, err
	}
	s.updateSessionLastUseTime(session.Name)
	tx, err := s.getTransactionByID(req.TransactionId)
	if err != nil {
		return nil, err
	}
	s.removeTransaction(tx)
	return &emptypb.Empty{}, nil
}

func (s *inMemSpannerServer) PartitionQuery(ctx context.Context, req *spannerpb.PartitionQueryRequest) (*spannerpb.PartitionResponse, error) {
	s.receivedRequests <- req
	return nil, gstatus.Error(codes.Unimplemented, "Method not yet implemented")
}

func (s *inMemSpannerServer) PartitionRead(ctx context.Context, req *spannerpb.PartitionReadRequest) (*spannerpb.PartitionResponse, error) {
	s.receivedRequests <- req
	return nil, gstatus.Error(codes.Unimplemented, "Method not yet implemented")
}
