package publish

type Destination interface {
	Upload(fileContents []byte, fileName string) error
	Path(fileName string) string
}
