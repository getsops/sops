// Package codes the exit statuses returned by the sops binary
package codes

// Exit statuses returned by the binary
const (
	ErrorGeneric                           int = 1
	CouldNotReadInputFile                  int = 2
	CouldNotWriteOutputFile                int = 3
	ErrorDumpingTree                       int = 4
	ErrorReadingConfig                     int = 5
	ErrorInvalidKMSEncryptionContextFormat int = 6
	ErrorInvalidSetFormat                  int = 7
	ErrorEncryptingMac                     int = 21
	ErrorEncryptingTree                    int = 23
	ErrorDecryptingMac                     int = 24
	ErrorDecryptingTree                    int = 25
	CannotChangeKeysFromNonExistentFile    int = 49
	MacMismatch                            int = 51
	MacNotFound                            int = 52
	ConfigFileNotFound                     int = 61
	KeyboardInterrupt                      int = 85
	InvalidTreePathFormat                  int = 91
	NoFileSpecified                        int = 100
	CouldNotRetrieveKey                    int = 128
	NoEncryptionKeyFound                   int = 111
	FileHasNotBeenModified                 int = 200
	NoEditorFound                          int = 201
	FailedToCompareVersions                int = 202
)
