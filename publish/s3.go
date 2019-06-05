package publish

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Destination struct {
	s3Bucket string
	s3Prefix string
}

func NewS3Destination(s3Bucket string, s3Prefix string) *S3Destination {
	return &S3Destination{s3Bucket, s3Prefix}
}

func (s3d *S3Destination) Path(fileName string) string {
	return fmt.Sprintf("s3://%s/%s%s", s3d.s3Bucket, s3d.s3Prefix, fileName)
}

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
