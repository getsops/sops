package lz4

import (
	"bytes"
	"github.com/pierrec/lz4"
	"io"
)

// Fuzz function for the Reader and Writer.
func Fuzz(data []byte) int {
	var (
		r      = bytes.NewReader(data)
		w      = new(bytes.Buffer)
		pr, pw = io.Pipe()
		zr     = lz4.NewReader(pr)
		zw     = lz4.NewWriter(pw)
	)
	// Compress.
	go func() {
		_, err := io.Copy(zw, r)
		if err != nil {
			panic(err)
		}
		err = zw.Close()
		if err != nil {
			panic(err)
		}
		err = pw.Close()
		if err != nil {
			panic(err)
		}
	}()
	// Decompress.
	_, err := io.Copy(w, zr)
	if err != nil {
		panic(err)
	}
	// Check that the data is valid.
	if !bytes.Equal(data, w.Bytes()) {
		panic("not equal")
	}
	return 1
}
