package publish

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
)

type GCSDestination struct {
	gcsBucket string
	gcsPrefix string
}

func NewGCSDestination(gcsBucket, gcsPrefix string) *GCSDestination {
	return &GCSDestination{gcsBucket, gcsPrefix}
}

func (gcsd *GCSDestination) Path(fileName string) string {
	return fmt.Sprintf("gcs://%s/%s%s", gcsd.gcsBucket, gcsd.gcsPrefix, fileName)
}

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

// Returns NotImplementedError
func (gcsd *GCSDestination) UploadUnencrypted(data map[string]interface{}, fileName string) error {
	return &NotImplementedError{"GCS does not support uploading the unencrypted file contents."}
}
