package rds

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/awstesting/unit"

	"github.com/stretchr/testify/assert"
)

func TestPresignWithPresignNotSet(t *testing.T) {
	reqs := map[string]*request.Request{}
	svc := New(unit.Session, &aws.Config{Region: aws.String("us-west-2")})

	assert.NotPanics(t, func() {
		// Doesn't panic on nil input
		req, _ := svc.CopyDBSnapshotRequest(nil)
		req.Sign()
	})

	reqs[opCopyDBSnapshot], _ = svc.CopyDBSnapshotRequest(&CopyDBSnapshotInput{
		SourceRegion:               aws.String("us-west-1"),
		SourceDBSnapshotIdentifier: aws.String("foo"),
		TargetDBSnapshotIdentifier: aws.String("bar"),
	})

	reqs[opCreateDBInstanceReadReplica], _ = svc.CreateDBInstanceReadReplicaRequest(&CreateDBInstanceReadReplicaInput{
		SourceRegion:               aws.String("us-west-1"),
		SourceDBInstanceIdentifier: aws.String("foo"),
		DBInstanceIdentifier:       aws.String("bar"),
	})

	for op, req := range reqs {
		req.Sign()
		b, _ := ioutil.ReadAll(req.HTTPRequest.Body)
		q, _ := url.ParseQuery(string(b))

		u, _ := url.QueryUnescape(q.Get("PreSignedUrl"))
		assert.Regexp(t, fmt.Sprintf(`^https://rds.us-west-1\.amazonaws\.com/\?Action=%s.+?DestinationRegion=us-west-2.+`, op), u)
	}
}

func TestPresignWithPresignSet(t *testing.T) {
	reqs := map[string]*request.Request{}
	svc := New(unit.Session, &aws.Config{Region: aws.String("us-west-2")})

	assert.NotPanics(t, func() {
		// Doesn't panic on nil input
		req, _ := svc.CopyDBSnapshotRequest(nil)
		req.Sign()
	})

	reqs[opCopyDBSnapshot], _ = svc.CopyDBSnapshotRequest(&CopyDBSnapshotInput{
		SourceRegion:               aws.String("us-west-1"),
		SourceDBSnapshotIdentifier: aws.String("foo"),
		TargetDBSnapshotIdentifier: aws.String("bar"),
		PreSignedUrl:               aws.String("presignedURL"),
	})

	reqs[opCreateDBInstanceReadReplica], _ = svc.CreateDBInstanceReadReplicaRequest(&CreateDBInstanceReadReplicaInput{
		SourceRegion:               aws.String("us-west-1"),
		SourceDBInstanceIdentifier: aws.String("foo"),
		DBInstanceIdentifier:       aws.String("bar"),
		PreSignedUrl:               aws.String("presignedURL"),
	})

	for _, req := range reqs {
		req.Sign()

		b, _ := ioutil.ReadAll(req.HTTPRequest.Body)
		q, _ := url.ParseQuery(string(b))

		u, _ := url.QueryUnescape(q.Get("PreSignedUrl"))
		assert.Regexp(t, `presignedURL`, u)
	}
}

func TestPresignWithSourceNotSet(t *testing.T) {
	reqs := map[string]*request.Request{}
	svc := New(unit.Session, &aws.Config{Region: aws.String("us-west-2")})

	assert.NotPanics(t, func() {
		// Doesn't panic on nil input
		req, _ := svc.CopyDBSnapshotRequest(nil)
		req.Sign()
	})

	reqs[opCopyDBSnapshot], _ = svc.CopyDBSnapshotRequest(&CopyDBSnapshotInput{
		SourceDBSnapshotIdentifier: aws.String("foo"),
		TargetDBSnapshotIdentifier: aws.String("bar"),
	})

	for _, req := range reqs {
		_, err := req.Presign(5 * time.Minute)
		if err != nil {
			t.Fatal(err)
		}
	}
}
