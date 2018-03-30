// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package firestore

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"testing"
	"time"

	"cloud.google.com/go/internal/pretty"
	"cloud.google.com/go/internal/testutil"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/type/latlng"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

func TestMain(m *testing.M) {
	initIntegrationTest()
	status := m.Run()
	cleanupIntegrationTest()
	os.Exit(status)
}

const (
	envProjID     = "GCLOUD_TESTS_GOLANG_FIRESTORE_PROJECT_ID"
	envPrivateKey = "GCLOUD_TESTS_GOLANG_FIRESTORE_KEY"
)

var (
	iClient       *Client
	iColl         *CollectionRef
	collectionIDs = testutil.NewUIDSpace("go-integration-test")
)

func initIntegrationTest() {
	flag.Parse() // needed for testing.Short()
	if testing.Short() {
		return
	}
	ctx := context.Background()
	testProjectID := os.Getenv(envProjID)
	if testProjectID == "" {
		log.Println("Integration tests skipped. See CONTRIBUTING.md for details")
		return
	}
	ts := testutil.TokenSourceEnv(ctx, envPrivateKey,
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/datastore")
	if ts == nil {
		log.Fatal("The project key must be set. See CONTRIBUTING.md for details")
	}
	ti := &testInterceptor{dbPath: "projects/" + testProjectID + "/databases/(default)"}
	c, err := NewClient(ctx, testProjectID,
		option.WithTokenSource(ts),
		option.WithGRPCDialOption(grpc.WithUnaryInterceptor(ti.interceptUnary)),
		option.WithGRPCDialOption(grpc.WithStreamInterceptor(ti.interceptStream)),
	)
	if err != nil {
		log.Fatalf("NewClient: %v", err)
	}
	iClient = c
	iColl = c.Collection(collectionIDs.New())
	refDoc := iColl.NewDoc()
	integrationTestMap["ref"] = refDoc
	wantIntegrationTestMap["ref"] = refDoc
	integrationTestStruct.Ref = refDoc
}

type testInterceptor struct {
	dbPath string
}

func (ti *testInterceptor) interceptUnary(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ti.checkMetadata(ctx, method)
	return invoker(ctx, method, req, res, cc, opts...)
}

func (ti *testInterceptor) interceptStream(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	ti.checkMetadata(ctx, method)
	return streamer(ctx, desc, cc, method, opts...)
}

func (ti *testInterceptor) checkMetadata(ctx context.Context, method string) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		log.Fatalf("method %s: bad metadata", method)
	}
	for _, h := range []string{"google-cloud-resource-prefix", "x-goog-api-client"} {
		v, ok := md[h]
		if !ok {
			log.Fatalf("method %s, header %s missing", method, h)
		}
		if len(v) != 1 {
			log.Fatalf("method %s, header %s: bad value %v", method, h, v)
		}
	}
	v := md["google-cloud-resource-prefix"][0]
	if v != ti.dbPath {
		log.Fatalf("method %s: bad resource prefix header:  %q", method, v)
	}
}

func cleanupIntegrationTest() {
	if iClient == nil {
		return
	}
	// TODO(jba): delete everything in integrationColl.
	iClient.Close()
}

// integrationClient should be called by integration tests to get a valid client. It will never
// return nil. If integrationClient returns, an integration test can proceed without
// further checks.
func integrationClient(t *testing.T) *Client {
	if testing.Short() {
		t.Skip("Integration tests skipped in short mode")
	}
	if iClient == nil {
		t.SkipNow() // log message printed in initIntegrationTest
	}
	return iClient
}

func integrationColl(t *testing.T) *CollectionRef {
	_ = integrationClient(t)
	return iColl
}

type integrationTestStructType struct {
	Int         int
	Str         string
	Bool        bool
	Float       float32
	Null        interface{}
	Bytes       []byte
	Time        time.Time
	Geo, NilGeo *latlng.LatLng
	Ref         *DocumentRef
}

var (
	integrationTime = time.Date(2017, 3, 20, 1, 2, 3, 456789, time.UTC)
	// Firestore times are accurate only to microseconds.
	wantIntegrationTime = time.Date(2017, 3, 20, 1, 2, 3, 456000, time.UTC)

	integrationGeo = &latlng.LatLng{Latitude: 30, Longitude: 70}

	// Use this when writing a doc.
	integrationTestMap = map[string]interface{}{
		"int":   1,
		"str":   "two",
		"bool":  true,
		"float": 3.14,
		"null":  nil,
		"bytes": []byte("bytes"),
		"*":     map[string]interface{}{"`": 4},
		"time":  integrationTime,
		"geo":   integrationGeo,
		"ref":   nil, // populated by initIntegrationTest
	}

	// The returned data is slightly different.
	wantIntegrationTestMap = map[string]interface{}{
		"int":   int64(1),
		"str":   "two",
		"bool":  true,
		"float": 3.14,
		"null":  nil,
		"bytes": []byte("bytes"),
		"*":     map[string]interface{}{"`": int64(4)},
		"time":  wantIntegrationTime,
		"geo":   integrationGeo,
		"ref":   nil, // populated by initIntegrationTest
	}

	integrationTestStruct = integrationTestStructType{
		Int:    1,
		Str:    "two",
		Bool:   true,
		Float:  3.14,
		Null:   nil,
		Bytes:  []byte("bytes"),
		Time:   integrationTime,
		Geo:    integrationGeo,
		NilGeo: nil,
		Ref:    nil, // populated by initIntegrationTest
	}
)

func TestIntegration_Create(t *testing.T) {
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()
	start := time.Now()
	wr := mustCreate("Create #1", t, doc, integrationTestMap)
	end := time.Now()
	checkTimeBetween(t, wr.UpdateTime, start, end)
	_, err := doc.Create(ctx, integrationTestMap)
	codeEq(t, "Create on a present doc", codes.AlreadyExists, err)
}

func TestIntegration_Get(t *testing.T) {
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()
	mustCreate("Get #1", t, doc, integrationTestMap)
	ds, err := doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if ds.CreateTime != ds.UpdateTime {
		t.Errorf("create time %s != update time %s", ds.CreateTime, ds.UpdateTime)
	}
	got := ds.Data()
	if want := wantIntegrationTestMap; !testEqual(got, want) {
		t.Errorf("got\n%v\nwant\n%v", pretty.Value(got), pretty.Value(want))
	}

	//
	_, err = integrationColl(t).NewDoc().Get(ctx)
	codeEq(t, "Get on a missing doc", codes.NotFound, err)
}

func TestIntegration_GetAll(t *testing.T) {
	type getAll struct{ N int }

	coll := integrationColl(t)
	ctx := context.Background()
	var docRefs []*DocumentRef
	for i := 0; i < 5; i++ {
		doc := coll.NewDoc()
		docRefs = append(docRefs, doc)
		mustCreate("GetAll #1", t, doc, getAll{N: i})
	}
	docSnapshots, err := iClient.GetAll(ctx, docRefs)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(docSnapshots), len(docRefs); got != want {
		t.Fatalf("got %d snapshots, want %d", got, want)
	}
	for i, ds := range docSnapshots {
		var got getAll
		if err := ds.DataTo(&got); err != nil {
			t.Fatal(err)
		}
		want := getAll{N: i}
		if got != want {
			t.Errorf("%d: got %+v, want %+v", i, got, want)
		}
	}
}

func TestIntegration_Add(t *testing.T) {
	start := time.Now()
	_, wr, err := integrationColl(t).Add(context.Background(), integrationTestMap)
	if err != nil {
		t.Fatal(err)
	}
	end := time.Now()
	checkTimeBetween(t, wr.UpdateTime, start, end)
}

func TestIntegration_Set(t *testing.T) {
	coll := integrationColl(t)
	ctx := context.Background()

	// Set Should be able to create a new doc.
	doc := coll.NewDoc()
	wr1, err := doc.Set(ctx, integrationTestMap)
	if err != nil {
		t.Fatal(err)
	}
	// Calling Set on the doc completely replaces the contents.
	// The update time should increase.
	newData := map[string]interface{}{
		"str": "change",
		"x":   "1",
	}
	wr2, err := doc.Set(ctx, newData)
	if err != nil {
		t.Fatal(err)
	}
	if !wr1.UpdateTime.Before(wr2.UpdateTime) {
		t.Errorf("update time did not increase: old=%s, new=%s", wr1.UpdateTime, wr2.UpdateTime)
	}
	ds, err := doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := ds.Data()
	if !testEqual(got, newData) {
		t.Errorf("got %v, want %v", got, newData)
	}

	newData = map[string]interface{}{
		"str": "1",
		"x":   "2",
		"y":   "3",
	}
	// SetOptions:
	// Only fields mentioned in the Merge option will be changed.
	// In this case, "str" will not be changed to "1".
	wr3, err := doc.Set(ctx, newData, Merge("x", "y"))
	if err != nil {
		t.Fatal(err)
	}
	ds, err = doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got = ds.Data()
	want := map[string]interface{}{
		"str": "change",
		"x":   "2",
		"y":   "3",
	}
	if !testEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if !wr2.UpdateTime.Before(wr3.UpdateTime) {
		t.Errorf("update time did not increase: old=%s, new=%s", wr2.UpdateTime, wr3.UpdateTime)
	}

	// Another way to change only x and y is to pass a map with only
	// those keys, and use MergeAll.
	wr4, err := doc.Set(ctx, map[string]interface{}{"x": "4", "y": "5"}, MergeAll)
	if err != nil {
		t.Fatal(err)
	}
	ds, err = doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got = ds.Data()
	want = map[string]interface{}{
		"str": "change",
		"x":   "4",
		"y":   "5",
	}
	if !testEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if !wr3.UpdateTime.Before(wr4.UpdateTime) {
		t.Errorf("update time did not increase: old=%s, new=%s", wr3.UpdateTime, wr4.UpdateTime)
	}
}

func TestIntegration_Delete(t *testing.T) {
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()
	mustCreate("Delete #1", t, doc, integrationTestMap)
	wr, err := doc.Delete(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// Confirm that doc doesn't exist.
	if _, err := doc.Get(ctx); grpc.Code(err) != codes.NotFound {
		t.Fatalf("got error <%v>, want NotFound", err)
	}

	er := func(_ *WriteResult, err error) error { return err }

	codeEq(t, "Delete on a missing doc", codes.OK,
		er(doc.Delete(ctx)))
	// TODO(jba): confirm that the server should return InvalidArgument instead of
	// FailedPrecondition.
	wr = mustCreate("Delete #2", t, doc, integrationTestMap)
	codeEq(t, "Delete with wrong LastUpdateTime", codes.FailedPrecondition,
		er(doc.Delete(ctx, LastUpdateTime(wr.UpdateTime.Add(-time.Millisecond)))))
	codeEq(t, "Delete with right LastUpdateTime", codes.OK,
		er(doc.Delete(ctx, LastUpdateTime(wr.UpdateTime))))
}

func TestIntegration_UpdateMap(t *testing.T) {
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()
	mustCreate("UpdateMap", t, doc, integrationTestMap)
	um := map[string]interface{}{
		"bool":        false,
		"time":        17,
		"null":        Delete,
		"noSuchField": Delete, // deleting a non-existent field is a no-op
	}
	wr, err := doc.UpdateMap(ctx, um)
	if err != nil {
		t.Fatal(err)
	}
	ds, err := doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := ds.Data()
	want := copyMap(wantIntegrationTestMap)
	want["bool"] = false
	want["time"] = int64(17)
	delete(want, "null")
	if !testEqual(got, want) {
		t.Errorf("got\n%#v\nwant\n%#v", got, want)
	}

	er := func(_ *WriteResult, err error) error { return err }
	codeEq(t, "UpdateMap on missing doc", codes.NotFound,
		er(integrationColl(t).NewDoc().UpdateMap(ctx, um)))
	codeEq(t, "UpdateMap with wrong LastUpdateTime", codes.FailedPrecondition,
		er(doc.UpdateMap(ctx, um, LastUpdateTime(wr.UpdateTime.Add(-time.Millisecond)))))
	codeEq(t, "UpdateMap with right LastUpdateTime", codes.OK,
		er(doc.UpdateMap(ctx, um, LastUpdateTime(wr.UpdateTime))))
}

func TestIntegration_UpdateStruct(t *testing.T) {
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()
	mustCreate("UpdateStruct", t, doc, integrationTestStruct)
	fields := []string{"Bool", "Time", "Null", "noSuchField"}
	wr, err := doc.UpdateStruct(ctx, fields,
		integrationTestStructType{
			Bool: false,
			Time: aTime2,
		})
	if err != nil {
		t.Fatal(err)
	}
	ds, err := doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var got integrationTestStructType
	if err := ds.DataTo(&got); err != nil {
		t.Fatal(err)
	}
	want := integrationTestStruct
	want.Bool = false
	want.Time = aTime2
	if !testEqual(got, want) {
		t.Errorf("got\n%#v\nwant\n%#v", got, want)
	}

	er := func(_ *WriteResult, err error) error { return err }
	codeEq(t, "UpdateStruct on missing doc", codes.NotFound,
		er(integrationColl(t).NewDoc().UpdateStruct(ctx, fields, integrationTestStruct)))
	codeEq(t, "UpdateStruct with wrong LastUpdateTime", codes.FailedPrecondition,
		er(doc.UpdateStruct(ctx, fields, integrationTestStruct, LastUpdateTime(wr.UpdateTime.Add(-time.Millisecond)))))
	codeEq(t, "UpdateStruct with right LastUpdateTime", codes.OK,
		er(doc.UpdateStruct(ctx, fields, integrationTestStruct, LastUpdateTime(wr.UpdateTime))))
}

func TestIntegration_UpdatePaths(t *testing.T) {
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()
	mustCreate("UpdatePaths", t, doc, integrationTestMap)
	fpus := []FieldPathUpdate{
		{Path: []string{"bool"}, Value: false},
		{Path: []string{"time"}, Value: 17},
		{Path: []string{"*", "`"}, Value: 18},
		{Path: []string{"null"}, Value: Delete},
		{Path: []string{"noSuchField"}, Value: Delete}, // deleting a non-existent field is a no-op
	}
	wr, err := doc.UpdatePaths(ctx, fpus)
	if err != nil {
		t.Fatal(err)
	}
	ds, err := doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := ds.Data()
	want := copyMap(wantIntegrationTestMap)
	want["bool"] = false
	want["time"] = int64(17)
	want["*"] = map[string]interface{}{"`": int64(18)}
	delete(want, "null")
	if !testEqual(got, want) {
		t.Errorf("got\n%#v\nwant\n%#v", got, want)
	}

	er := func(_ *WriteResult, err error) error { return err }

	codeEq(t, "UpdatePaths on missing doc", codes.NotFound,
		er(integrationColl(t).NewDoc().UpdatePaths(ctx, fpus)))
	codeEq(t, "UpdatePaths with wrong LastUpdateTime", codes.FailedPrecondition,
		er(doc.UpdatePaths(ctx, fpus, LastUpdateTime(wr.UpdateTime.Add(-time.Millisecond)))))
	codeEq(t, "UpdatePaths with right LastUpdateTime", codes.OK,
		er(doc.UpdatePaths(ctx, fpus, LastUpdateTime(wr.UpdateTime))))
}

func TestIntegration_Collections(t *testing.T) {
	ctx := context.Background()
	c := integrationClient(t)
	got, err := c.Collections(ctx).GetAll()
	if err != nil {
		t.Fatal(err)
	}
	// There should be at least one collection.
	if len(got) == 0 {
		t.Error("got 0 top-level collections, want at least one")
	}

	doc := integrationColl(t).NewDoc()
	got, err = doc.Collections(ctx).GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("got %d collections, want 0", len(got))
	}
	var want []*CollectionRef
	for i := 0; i < 3; i++ {
		id := collectionIDs.New()
		cr := doc.Collection(id)
		want = append(want, cr)
		mustCreate("Collections", t, cr.NewDoc(), integrationTestMap)
	}
	got, err = doc.Collections(ctx).GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if !testEqual(got, want) {
		t.Errorf("got\n%#v\nwant\n%#v", got, want)
	}
}

func TestIntegration_ServerTimestamp(t *testing.T) {
	type S struct {
		A int
		B time.Time
		C time.Time `firestore:"C.C,serverTimestamp"`
		D map[string]interface{}
		E time.Time `firestore:",omitempty,serverTimestamp"`
	}
	data := S{
		A: 1,
		B: aTime,
		// C is unset, so will get the server timestamp.
		D: map[string]interface{}{"x": ServerTimestamp},
		// E is unset, so will get the server timestamp.
	}
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()
	// Bound times of the RPC, with some slack for clock skew.
	start := time.Now()
	mustCreate("ServerTimestamp", t, doc, data)
	end := time.Now()
	ds, err := doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var got S
	if err := ds.DataTo(&got); err != nil {
		t.Fatal(err)
	}
	if !testEqual(got.B, aTime) {
		t.Errorf("B: got %s, want %s", got.B, aTime)
	}
	checkTimeBetween(t, got.C, start, end)
	if g, w := got.D["x"], got.C; !testEqual(g, w) {
		t.Errorf(`D["x"] = %s, want equal to C (%s)`, g, w)
	}
	if g, w := got.E, got.C; !testEqual(g, w) {
		t.Errorf(`E = %s, want equal to C (%s)`, g, w)
	}
}

func TestIntegration_MergeServerTimestamp(t *testing.T) {
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()

	// Create a doc with an ordinary field "a" and a ServerTimestamp field "b".
	_, err := doc.Set(ctx, map[string]interface{}{
		"a": 1,
		"b": ServerTimestamp})
	if err != nil {
		t.Fatal(err)
	}
	docSnap, err := doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	data1 := docSnap.Data()
	// Merge with a document with a different value of "a". However,
	// specify only "b" in the list of merge fields.
	_, err = doc.Set(ctx,
		map[string]interface{}{"a": 2, "b": ServerTimestamp},
		Merge("b"))
	if err != nil {
		t.Fatal(err)
	}
	// The result should leave "a" unchanged, while "b" is updated.
	docSnap, err = doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	data2 := docSnap.Data()
	if got, want := data2["a"], data1["a"]; got != want {
		t.Errorf("got %v, want %v", got, want)
	}
	t1 := data1["b"].(time.Time)
	t2 := data2["b"].(time.Time)
	if !t1.Before(t2) {
		t.Errorf("got t1=%s, t2=%s; want t1 before t2", t1, t2)
	}
}

func TestIntegration_MergeNestedServerTimestamp(t *testing.T) {
	ctx := context.Background()
	doc := integrationColl(t).NewDoc()

	// Create a doc with an ordinary field "a" a ServerTimestamp field "b",
	// and a second ServerTimestamp field "c.d".
	_, err := doc.Set(ctx, map[string]interface{}{
		"a": 1,
		"b": ServerTimestamp,
		"c": map[string]interface{}{"d": ServerTimestamp},
	})
	if err != nil {
		t.Fatal(err)
	}
	docSnap, err := doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	data1 := docSnap.Data()
	// Merge with a document with a different value of "a". However,
	// specify only "c.d" in the list of merge fields.
	_, err = doc.Set(ctx,
		map[string]interface{}{
			"a": 2,
			"b": ServerTimestamp,
			"c": map[string]interface{}{"d": ServerTimestamp},
		},
		Merge("c.d"))
	if err != nil {
		t.Fatal(err)
	}
	// The result should leave "a" and "b" unchanged, while "c.d" is updated.
	docSnap, err = doc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	data2 := docSnap.Data()
	if got, want := data2["a"], data1["a"]; got != want {
		t.Errorf("a: got %v, want %v", got, want)
	}
	want := data1["b"].(time.Time)
	got := data2["b"].(time.Time)
	if !got.Equal(want) {
		t.Errorf("b: got %s, want %s", got, want)
	}
	t1 := data1["c"].(map[string]interface{})["d"].(time.Time)
	t2 := data2["c"].(map[string]interface{})["d"].(time.Time)
	if !t1.Before(t2) {
		t.Errorf("got t1=%s, t2=%s; want t1 before t2", t1, t2)
	}
}

func TestIntegration_WriteBatch(t *testing.T) {
	ctx := context.Background()
	b := integrationClient(t).Batch()
	doc1 := iColl.NewDoc()
	doc2 := iColl.NewDoc()
	b.Create(doc1, integrationTestMap)
	b.Set(doc2, integrationTestMap)
	b.UpdateMap(doc1, map[string]interface{}{"bool": false})
	b.UpdateMap(doc1, map[string]interface{}{"str": Delete})

	wrs, err := b.Commit(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(wrs), 4; got != want {
		t.Fatalf("got %d WriteResults, want %d", got, want)
	}
	ds, err := doc1.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got1 := ds.Data()
	want := copyMap(wantIntegrationTestMap)
	want["bool"] = false
	delete(want, "str")
	if !testEqual(got1, want) {
		t.Errorf("got\n%#v\nwant\n%#v", got1, want)
	}
	ds, err = doc2.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got2 := ds.Data()
	if !testEqual(got2, wantIntegrationTestMap) {
		t.Errorf("got\n%#v\nwant\n%#v", got2, wantIntegrationTestMap)
	}
	// TODO(jba): test two updates to the same document when it is supported.
	// TODO(jba): test verify when it is supported.
}

func TestIntegration_Query(t *testing.T) {
	ctx := context.Background()
	coll := integrationColl(t)
	var docs []*DocumentRef
	var wants []map[string]interface{}
	for i := 0; i < 3; i++ {
		doc := coll.NewDoc()
		docs = append(docs, doc)
		// To support running this test in parallel with the others, use a field name
		// that we don't use anywhere else.
		mustCreate(fmt.Sprintf("Query #%d", i), t, doc,
			map[string]interface{}{
				"q": i,
				"x": 1,
			})
		wants = append(wants, map[string]interface{}{"q": int64(i)})
	}
	q := coll.Select("q").OrderBy("q", Asc)
	for i, test := range []struct {
		q    Query
		want []map[string]interface{}
	}{
		{q, wants},
		{q.Where("q", ">", 1), wants[2:]},
		{q.WherePath([]string{"q"}, ">", 1), wants[2:]},
		{q.Offset(1).Limit(1), wants[1:2]},
		{q.StartAt(1), wants[1:]},
		{q.StartAfter(1), wants[2:]},
		{q.EndAt(1), wants[:2]},
		{q.EndBefore(1), wants[:1]},
	} {
		gotDocs, err := test.q.Documents(ctx).GetAll()
		if err != nil {
			t.Errorf("#%d: %+v: %v", i, test.q, err)
			continue
		}
		if len(gotDocs) != len(test.want) {
			t.Errorf("#%d: %+v: got %d docs, want %d", i, test.q, len(gotDocs), len(test.want))
			continue
		}
		for j, g := range gotDocs {
			if got, want := g.Data(), test.want[j]; !testEqual(got, want) {
				t.Errorf("#%d: %+v, #%d: got\n%+v\nwant\n%+v", i, test.q, j, got, want)
			}
		}
	}
	_, err := coll.Select("q").Where("x", "==", 1).OrderBy("q", Asc).Documents(ctx).GetAll()
	codeEq(t, "Where and OrderBy on different fields without an index", codes.FailedPrecondition, err)

	// Using the collection itself as the query should return the full documents.
	allDocs, err := coll.Documents(ctx).GetAll()
	if err != nil {
		t.Fatal(err)
	}
	seen := map[int64]bool{} // "q" values we see
	for _, d := range allDocs {
		data := d.Data()
		q, ok := data["q"]
		if !ok {
			// A document from another test.
			continue
		}
		if seen[q.(int64)] {
			t.Errorf("%v: duplicate doc", data)
		}
		seen[q.(int64)] = true
		if data["x"] != int64(1) {
			t.Errorf("%v: wrong or missing 'x'", data)
		}
		if len(data) != 2 {
			t.Errorf("%v: want two keys", data)
		}
	}
	if got, want := len(seen), len(wants); got != want {
		t.Errorf("got %d docs with 'q', want %d", len(seen), len(wants))
	}
}

// Test the special DocumentID field in queries.
func TestIntegration_QueryName(t *testing.T) {
	ctx := context.Background()

	checkIDs := func(q Query, wantIDs []string) {
		gots, err := q.Documents(ctx).GetAll()
		if err != nil {
			t.Fatal(err)
		}
		if len(gots) != len(wantIDs) {
			t.Fatalf("got %d, want %d", len(gots), len(wantIDs))
		}
		for i, g := range gots {
			if got, want := g.Ref.ID, wantIDs[i]; got != want {
				t.Errorf("#%d: got %s, want %s", i, got, want)
			}
		}
	}

	coll := integrationColl(t)
	var wantIDs []string
	for i := 0; i < 3; i++ {
		doc := coll.NewDoc()
		mustCreate(fmt.Sprintf("Query #%d", i), t, doc, map[string]interface{}{"nm": 1})
		wantIDs = append(wantIDs, doc.ID)
	}
	sort.Strings(wantIDs)
	q := coll.Where("nm", "==", 1).OrderBy(DocumentID, Asc)
	checkIDs(q, wantIDs)

	// Empty Select.
	q = coll.Select().Where("nm", "==", 1).OrderBy(DocumentID, Asc)
	checkIDs(q, wantIDs)

	// Test cursors with __name__.
	checkIDs(q.StartAt(wantIDs[1]), wantIDs[1:])
	checkIDs(q.EndAt(wantIDs[1]), wantIDs[:2])
}

func TestIntegration_QueryNested(t *testing.T) {
	ctx := context.Background()
	coll1 := integrationColl(t)
	doc1 := coll1.NewDoc()
	coll2 := doc1.Collection(collectionIDs.New())
	doc2 := coll2.NewDoc()
	wantData := map[string]interface{}{"x": int64(1)}
	mustCreate("QueryNested", t, doc2, wantData)
	q := coll2.Select("x")
	got, err := q.Documents(ctx).GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d docs, want 1", len(got))
	}
	if gotData := got[0].Data(); !testEqual(gotData, wantData) {
		t.Errorf("got\n%+v\nwant\n%+v", gotData, wantData)
	}
}

func TestIntegration_RunTransaction(t *testing.T) {
	ctx := context.Background()
	type Player struct {
		Name  string
		Score int
		Star  bool `firestore:"*"`
	}
	pat := Player{Name: "Pat", Score: 3, Star: false}
	client := integrationClient(t)
	patDoc := iColl.Doc("pat")
	var anError error
	incPat := func(_ context.Context, tx *Transaction) error {
		doc, err := tx.Get(patDoc)
		if err != nil {
			return err
		}
		score, err := doc.DataAt("Score")
		if err != nil {
			return err
		}
		// Since the Star field is called "*", we must use DataAtPath to get it.
		star, err := doc.DataAtPath([]string{"*"})
		if err != nil {
			return err
		}
		err = tx.UpdateStruct(patDoc, []string{"Score"},
			Player{Score: int(score.(int64) + 7)})
		if err != nil {
			return err
		}
		// Since the Star field is called "*", we must use UpdatePaths to change it.
		err = tx.UpdatePaths(patDoc,
			[]FieldPathUpdate{{Path: []string{"*"}, Value: !star.(bool)}})
		if err != nil {
			return err
		}
		return anError
	}
	mustCreate("RunTransaction", t, patDoc, pat)
	err := client.RunTransaction(ctx, incPat)
	if err != nil {
		t.Fatal(err)
	}
	ds, err := patDoc.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var got Player
	if err := ds.DataTo(&got); err != nil {
		t.Fatal(err)
	}
	want := Player{Name: "Pat", Score: 10, Star: true}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}

	// Function returns error, so transaction is rolled back and no writes happen.
	anError = errors.New("bad")
	err = client.RunTransaction(ctx, incPat)
	if err != anError {
		t.Fatalf("got %v, want %v", err, anError)
	}
	if err := ds.DataTo(&got); err != nil {
		t.Fatal(err)
	}
	// want is same as before.
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func codeEq(t *testing.T, msg string, code codes.Code, err error) {
	if grpc.Code(err) != code {
		t.Fatalf("%s:\ngot <%v>\nwant code %s", msg, err, code)
	}
}

func mustCreate(msg string, t *testing.T, doc *DocumentRef, data interface{}) *WriteResult {
	wr, err := doc.Create(context.Background(), data)
	if err != nil {
		t.Fatalf("%s: creating: %v", msg, err)
	}
	return wr
}

func copyMap(m map[string]interface{}) map[string]interface{} {
	c := map[string]interface{}{}
	for k, v := range m {
		c[k] = v
	}
	return c
}

func checkTimeBetween(t *testing.T, got, low, high time.Time) {
	// Allow slack for clock skew.
	const slack = 2 * time.Second
	low = low.Add(-slack)
	high = high.Add(slack)
	if got.Before(low) || got.After(high) {
		t.Fatalf("got %s, not in [%s, %s]", got, low, high)
	}
}
