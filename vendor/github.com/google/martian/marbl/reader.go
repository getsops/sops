// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package marbl

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

// Header is either an HTTP header or meta-data pertaining to the request or response.
type Header struct {
	ID          string
	MessageType MessageType
	Name        string
	Value       string
}

// String returns the contents of a Header frame in a format appropriate for debugging and runtime logging.
func (hf Header) String() string {
	return fmt.Sprintf("ID=%s; Type=%d; Name=%s; Value=%s", hf.ID, hf.MessageType, hf.Name, hf.Value)
}

// FrameType returns HeaderFrame
func (hf Header) FrameType() FrameType {
	return HeaderFrame
}

// Data is the payload (body) of the request or response.
type Data struct {
	ID          string
	MessageType MessageType
	Index       uint32
	Terminal    bool
	Data        []byte
}

// String returns the contents of a Data frame in a format appropriate for debugging and runtime logging. The
// contents of the data content slice (df.Data) is not printed, instead the length of Data is printed.
func (df Data) String() string {
	return fmt.Sprintf("ID=%s; Type=%d; Index=%d; Terminal=%t; Data length=%d",
		df.ID, df.MessageType, df.Index, df.Terminal, len(df.Data))
}

// FrameType returns DataFrame
func (df Data) FrameType() FrameType {
	return DataFrame
}

// Frame describes the interface for a frame (either Data or Header).
type Frame interface {
	String() string
	FrameType() FrameType
}

// Reader wraps a buffered Reader that reads from the io.Reader and emits Frames.
type Reader struct {
	r io.Reader
}

// NewReader returns a Reader initialized with a buffered reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

// ReadFrame reads from r, determines the FrameType, and returns either a Header or Data and an error.
func (r *Reader) ReadFrame() (Frame, error) {
	fh := make([]byte, 10)

	if _, err := io.ReadFull(r.r, fh); err != nil {
		return nil, err
	}

	switch FrameType(fh[0]) {
	case HeaderFrame:
		hf := Header{
			ID:          string(fh[2:]),
			MessageType: MessageType(fh[1]),
		}

		lens := make([]byte, 8)
		if _, err := io.ReadFull(r.r, lens); err != nil {
			return nil, err
		}

		nl := binary.BigEndian.Uint32(lens[:4])
		vl := binary.BigEndian.Uint32(lens[4:])

		nv := make([]byte, int(nl+vl))
		if _, err := io.ReadFull(r.r, nv); err != nil {
			return nil, err
		}

		hf.Name = string(nv[:nl])
		hf.Value = string(nv[nl:])

		return hf, nil
	case DataFrame:
		df := Data{
			ID:          string(fh[2:]),
			MessageType: MessageType(fh[1]),
		}

		// Reading 9 bytes:
		// 4 bytes index
		// 1 byte terminal
		// 4 bytes data length
		desc := make([]byte, 9)
		if _, err := io.ReadFull(r.r, desc); err != nil {
			return nil, err
		}

		df.Index = binary.BigEndian.Uint32(desc[:4])
		df.Terminal = desc[4] == 1

		dl := binary.BigEndian.Uint32(desc[5:])


		data := make([]byte, int(dl))
		if _, err := io.ReadFull(r.r, data); err != nil {
			return nil, err
		}

		df.Data = data

		return df, nil
	default:
		return nil, fmt.Errorf("marbl: unknown type of frame")
	}
}
