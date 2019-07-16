package publish

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Destination is the AWS S3 implementation of the Destination interface
type S3Destination struct {
	s3Bucket string
	s3Prefix string
}

// NewS3Destination is the constructor for an S3 Destination
func NewS3Destination(s3Bucket, s3Prefix string) *S3Destination {
	return &S3Destination{s3Bucket, s3Prefix}
}

// Path returns the S3 path of a file in an S3 Destination (bucket)
func (s3d *S3Destination) Path(fileName string) string {
	return fmt.Sprintf("s3://%s/%s%s", s3d.s3Bucket, s3d.s3Prefix, fileName)
}

// Upload uploads contents to a file in an S3 Destination (bucket)
func (s3d *S3Destination) Upload(fileContents []byte, fileName string) error {
	sess := session.Must(session.NewSession())
	svc := s3.New(sess)
	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(bytes.NewReader(fileContents)),
		Bucket: aws.String(s3d.s3Bucket),
		Key:    aws.String(s3d.s3Prefix + fileName),
	}
	_, err := svc.PutObject(input)
	if err != nil {
		return err
	}
	return nil
}

// Returns NotImplementedError
func (s3d *S3Destination) UploadUnencrypted(data map[string]interface{}, fileName string) error {
	return &NotImplementedError{"S3 does not support uploading the unencrypted file contents."}
}
