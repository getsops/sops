package publish

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

// GCSDestination represents the Google Cloud Storage destination
type GCSDestination struct {
	gcsBucket string
	gcsPrefix string
}

// NewGCSDestination is the constructor for a new Google Cloud Storage destination
func NewGCSDestination(gcsBucket string, gcsPrefix string) *GCSDestination {
	return &GCSDestination{gcsBucket, gcsPrefix}
}

// Path returns a the GCS path for a file within this GCS Destination
func (gcsd *GCSDestination) Path(fileName string) string {
	return fmt.Sprintf("gcs://%s/%s%s", gcsd.gcsBucket, gcsd.gcsPrefix, fileName)
}

// Upload uploads contents to a new file in GCS
func (gcsd *GCSDestination) Upload(fileContents []byte, fileName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	wc := client.Bucket(gcsd.gcsBucket).Object(gcsd.gcsPrefix + fileName).NewWriter(ctx)
	defer wc.Close()
	_, err = wc.Write(fileContents)
	if err != nil {
		return err
	}
	return nil
}
