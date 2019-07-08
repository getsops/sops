package publish

// Destination represents actions which all destination types
// must implement in order to be used by SOPS
type Destination interface {
	Upload(fileContents []byte, fileName string) error
	Path(fileName string) string
}
