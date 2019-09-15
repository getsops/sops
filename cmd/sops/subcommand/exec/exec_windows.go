package exec

import (
	"log"
)

func WritePipe(pipe string, contents []byte) {
	log.Fatal("fifos are not available on windows")
}

func GetPipe(dir string) string {
	log.Fatal("fifos not available on windows")
	return ""
}

func SwitchUser(username string) {
	log.Fatal("user switching not available on windows")
}


