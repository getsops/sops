package lz4_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/pierrec/lz4"
)

func TestWriter(t *testing.T) {
	goldenFiles := []string{
		"testdata/e.txt",
		"testdata/gettysburg.txt",
		"testdata/Mark.Twain-Tom.Sawyer.txt",
		"testdata/Mark.Twain-Tom.Sawyer_long.txt",
		"testdata/pg1661.txt",
		"testdata/pi.txt",
		"testdata/random.data",
		"testdata/repeat.txt",
	}

	for _, fname := range goldenFiles {
		for _, header := range []lz4.Header{
			{}, // Default header.
			{BlockChecksum: true},
			{NoChecksum: true},
			{BlockMaxSize: 64 << 10}, // 64Kb
			{CompressionLevel: 10},
			{Size: 123},
		} {
			label := fmt.Sprintf("%s/%s", fname, header)
			t.Run(label, func(t *testing.T) {
				fname := fname
				header := header
				t.Parallel()

				raw, err := ioutil.ReadFile(fname)
				if err != nil {
					t.Fatal(err)
				}
				r := bytes.NewReader(raw)

				// Compress.
				var zout bytes.Buffer
				zw := lz4.NewWriter(&zout)
				zw.Header = header
				_, err = io.Copy(zw, r)
				if err != nil {
					t.Fatal(err)
				}
				err = zw.Close()
				if err != nil {
					t.Fatal(err)
				}

				// Uncompress.
				var out bytes.Buffer
				zr := lz4.NewReader(&zout)
				n, err := io.Copy(&out, zr)
				if err != nil {
					t.Fatal(err)
				}

				// The uncompressed data must be the same as the initial input.
				if got, want := int(n), len(raw); got != want {
					t.Errorf("invalid sizes: got %d; want %d", got, want)
				}

				if got, want := out.Bytes(), raw; !reflect.DeepEqual(got, want) {
					t.Fatal("uncompressed data does not match original")
				}
			})
		}
	}
}

func TestIssue41(t *testing.T) {
	r, w := io.Pipe()
	zw := lz4.NewWriter(w)
	zr := lz4.NewReader(r)

	data := "x"
	go func() {
		_, _ = fmt.Fprint(zw, data)
		_ = zw.Flush()
		_ = zw.Close()
		_ = w.Close()
	}()
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(zr)
	if got, want := buf.String(), data; got != want {
		t.Fatal("uncompressed data does not match original")
	}
}

func TestIssue43(t *testing.T) {
	r, w := io.Pipe()
	go func() {
		defer w.Close()

		f, err := os.Open("testdata/issue43.data")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		zw := lz4.NewWriter(w)
		defer zw.Close()

		_, err = io.Copy(zw, f)
		if err != nil {
			t.Fatal(err)
		}
	}()
	_, err := io.Copy(ioutil.Discard, lz4.NewReader(r))
	if err != nil {
		t.Fatal(err)
	}
}
