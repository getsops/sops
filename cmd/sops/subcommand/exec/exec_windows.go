package exec

import (
	"fmt"
	"os/exec"
)

func ExecSyscall(command string, env []string) error {
	log.Fatal("same-process not available on windows")
	return nil
}

func BuildCommand(command string) *exec.Cmd {
	return exec.Command("cmd.exe", "/C", command)
}

func WritePipe(pipe string, contents []byte) {
	log.Fatal("fifos are not available on windows")
}

func GetPipe(dir, filename string) (string, error) {
	return "", fmt.Errorf("fifos are not available on windows")
}

func SwitchUser(username string) {
	log.Fatal("user switching not available on windows")
}
