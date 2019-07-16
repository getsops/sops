package main

import (
	"flag"
	"fmt"

	"github.com/pierrec/cmdflag"
)

func main() {
	flag.CommandLine.Bool(cmdflag.VersionBoolFlag, false, "print the program version")

	cli := cmdflag.New(nil)
	cli.MustAdd(cmdflag.Application{
		Name:  "compress",
		Args:  "[arguments] [<file name> ...]",
		Descr: "Compress the given files or from stdin to stdout.",
		Err:   flag.ExitOnError,
		Init:  Compress,
	})
	cli.MustAdd(cmdflag.Application{
		Name:  "uncompress",
		Args:  "[arguments] [<file name> ...]",
		Descr: "Uncompress the given files or from stdin to stdout.",
		Err:   flag.ExitOnError,
		Init:  Uncompress,
	})

	if err := cli.Parse(); err != nil {
		fmt.Println(err)
		return
	}
}
