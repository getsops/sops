package prefixer

import (
	"bufio"
	"io"
)

// Prefixer implements io.Reader and io.WriterTo. It reads
// data from the underlying reader and prepends every line
// with a given string.
type Prefixer struct {
	reader *bufio.Reader
	prefix []byte
	unread []byte
	eof    bool
}

// New creates a new instance of Prefixer.
func New(r io.Reader, prefix string) *Prefixer {
	return &Prefixer{
		reader: bufio.NewReader(r),
		prefix: []byte(prefix),
	}
}

// Read implements io.Reader. It reads data into p from the
// underlying reader and prepends every line with a prefix.
// It does not block if no data is available yet.
// It returns the number of bytes read into p.
func (r *Prefixer) Read(p []byte) (n int, err error) {
	for {
		// Write unread data from previous read.
		if len(r.unread) > 0 {
			m := copy(p[n:], r.unread)
			n += m
			r.unread = r.unread[m:]
			if len(r.unread) > 0 {
				return n, nil
			}
		}

		// The underlying Reader already returned EOF, do not read again.
		if r.eof {
			return n, io.EOF
		}

		// Read new line, including delim.
		r.unread, err = r.reader.ReadBytes('\n')

		if err == io.EOF {
			r.eof = true
		}

		// No new data, do not block.
		if len(r.unread) == 0 {
			return n, err
		}

		// Some new data, prepend prefix.
		// TODO: We could write the prefix to r.unread buffer just once
		//       and re-use it instead of prepending every time.
		r.unread = append(r.prefix, r.unread...)

		if err != nil {
			if err == io.EOF && len(r.unread) > 0 {
				// The underlying Reader already returned EOF, but we still
				// have some unread data to send, thus clear the error.
				return n, nil
			}
			return n, err
		}
	}
	panic("unreachable")
}

func (r *Prefixer) WriteTo(w io.Writer) (n int64, err error) {
	for {
		// Write unread data from previous read.
		if len(r.unread) > 0 {
			m, err := w.Write(r.unread)
			n += int64(m)
			if err != nil {
				return n, err
			}
			r.unread = r.unread[m:]
			if len(r.unread) > 0 {
				return n, nil
			}
		}

		// The underlying Reader already returned EOF, do not read again.
		if r.eof {
			return n, io.EOF
		}

		// Read new line, including delim.
		r.unread, err = r.reader.ReadBytes('\n')

		if err == io.EOF {
			r.eof = true
		}

		// No new data, do not block.
		if len(r.unread) == 0 {
			return n, err
		}

		// Some new data, prepend prefix.
		// TODO: We could write the prefix to r.unread buffer just once
		//       and re-use it instead of prepending every time.
		r.unread = append(r.prefix, r.unread...)

		if err != nil {
			if err == io.EOF && len(r.unread) > 0 {
				// The underlying Reader already returned EOF, but we still
				// have some unread data to send, thus clear the error.
				return n, nil
			}
			return n, err
		}
	}
	panic("unreachable")
}
