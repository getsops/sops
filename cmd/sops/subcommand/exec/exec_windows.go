package exec

import (
	"os/exec"
)

func BuildCommand(command string, shell string) *exec.Cmd {
	if shell == "" {
		return exec.Command("cmd.exe", "/C", command)
	}
	return exec.Command(shell, "/C", command)
}

func WritePipe(pipe string, contents []byte) {
	log.Fatal("fifos are not available on windows")
}

func GetPipe(dir, filename string) string {
	log.Fatal("fifos are not available on windows")
	return ""
}

func SwitchUser(username string) {
	log.Fatal("user switching not available on windows")
}
