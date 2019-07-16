package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/schollz/progressbar/v2"

	"github.com/pierrec/cmdflag"
	"github.com/pierrec/lz4"
)

// Uncompress uncompresses a set of files or from stdin to stdout.
func Uncompress(_ *flag.FlagSet) cmdflag.Handler {
	return func(args ...string) (int, error) {
		zr := lz4.NewReader(nil)

		// Use stdin/stdout if no file provided.
		if len(args) == 0 {
			zr.Reset(os.Stdin)
			_, err := io.Copy(os.Stdout, zr)
			return 0, err
		}

		for fidx, zfilename := range args {
			// Input file.
			zfile, err := os.Open(zfilename)
			if err != nil {
				return fidx, err
			}
			zinfo, err := zfile.Stat()
			if err != nil {
				return fidx, err
			}
			mode := zinfo.Mode() // use the same mode for the output file

			// Output file.
			filename := strings.TrimSuffix(zfilename, lz4.Extension)
			file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, mode)
			if err != nil {
				return fidx, err
			}
			zr.Reset(zfile)

			zfinfo, err := zfile.Stat()
			if err != nil {
				return fidx, err
			}
			var (
				size  int
				out   io.Writer = file
				zsize           = zfinfo.Size()
				bar   *progressbar.ProgressBar
			)
			if zsize > 0 {
				bar = progressbar.NewOptions64(zsize,
					// File transfers are usually slow, make sure we display the bar at 0%.
					progressbar.OptionSetRenderBlankState(true),
					// Display the filename.
					progressbar.OptionSetDescription(filename),
					progressbar.OptionClearOnFinish(),
				)
				out = io.MultiWriter(out, bar)
				zr.OnBlockDone = func(n int) {
					size += n
				}
			}

			// Uncompress.
			_, err = io.Copy(out, zr)
			if err != nil {
				return fidx, err
			}
			for _, c := range []io.Closer{zfile, file} {
				err := c.Close()
				if err != nil {
					return fidx, err
				}
			}

			if bar != nil {
				_ = bar.Clear()
				fmt.Printf("%s %d\n", zfilename, size)
			}
		}

		return len(args), nil
	}
}
