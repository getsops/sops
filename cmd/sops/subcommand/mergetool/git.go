package mergetool

import (
	"fmt"
	"os"
)

// Reference: https://git-scm.com/docs/git-mergetool

type MergePaths struct {
	Base   string
	Local  string
	Remote string
	Merged string
}

func MergePathsFromEnv() MergePaths {
	return MergePaths{
		Base:   os.Getenv("BASE"),
		Local:  os.Getenv("LOCAL"),
		Remote: os.Getenv("REMOTE"),
		Merged: os.Getenv("MERGED"),
	}
}

func (mp MergePaths) String() string {
	return fmt.Sprintf("BASE=%s LOCAL=%s REMOTE=%s MERGED=%s", mp.Base, mp.Local, mp.Remote, mp.Merged)
}
