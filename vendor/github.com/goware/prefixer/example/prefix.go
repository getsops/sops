package main

import (
	"io/ioutil"
	"os"

	"github.com/goware/prefixer"
)

func main() {
	// Prefixer accepts anything that implements io.Reader interface
	prefixReader := prefixer.New(os.Stdin, "> ")

	// Read all prefixed lines from STDIN into a buffer
	buffer, _ := ioutil.ReadAll(prefixReader)

	// Write buffer to STDOUT
	os.Stdout.Write(buffer)
}
