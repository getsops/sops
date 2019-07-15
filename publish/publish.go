package publish

import "fmt"

type Destination interface {
	Upload(fileContents []byte, fileName string) error
	UploadUnencrypted(data map[string]interface{}, fileName string) error
	Path(fileName string) string
}

type NotImplementedError struct {
	message string
}

func (e *NotImplementedError) Error() string {
	return fmt.Sprintf("NotImplementedError: %s", e.message)
}
