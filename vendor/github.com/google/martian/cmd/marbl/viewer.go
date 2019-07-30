// Copyright 2018 Google Inc. All rights reserved.
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

// Command-line tool to view .marbl files. This tool reads all headers from provided .marbl
// file and prints them to stdout. Bodies of request/response are not printed to stdout,
// instead they are saved into individual files in form of "marbl_ID_TYPE" where
// ID is the ID of request or response and TYPE is "request" or "response".
//
// Command line arguments:
//   --file  Path to the .marbl file to view.
//   --out   Optional, folder where this tool will save request/response bodies.
//           uses current folder by default.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/google/martian/v3/marbl"
)

var (
	file = flag.String("file", "", ".marbl file to show contents of")
	out  = flag.String("out", "", "folder to write request/response bodies to. Folder must exist.")
)

func main() {
	flag.Parse()

	if *file == "" {
		fmt.Println("--file flag is required")
		return
	}

	file, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}

	reader := marbl.NewReader(file)

	// Iterate through all frames in .marbl file.
	for {
		frame, err := reader.ReadFrame()
		if frame == nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("reader.ReadFrame(): got %v, want no error or io.EOF\n", err)
			break
		}

		// Print current frame to stdout.
		if frame.FrameType() == marbl.HeaderFrame {
			fmt.Print("Header ")
		} else {
			fmt.Print("Data ")
		}
		fmt.Println(frame.String())

		// If frame is Data then we write it into separate
		// file that can be inspected later.
		if frame.FrameType() == marbl.DataFrame {
			df := frame.(marbl.Data)
			var t string
			if df.MessageType == marbl.Request {
				t = "request"
			} else if df.MessageType == marbl.Response {
				t = "response"
			} else {
				t = fmt.Sprintf("unknown_%d", df.MessageType)
			}
			fout := fmt.Sprintf("marbl_%s_%s", df.ID, t)
			if *out != "" {
				fout = *out + "/" + fout
			}
			fmt.Printf("Appending data to file %s\n", fout)

			// Append data to the file. Note that body can be split
			// into multiple frames so we have to append and not overwrite.
			f, err := os.OpenFile(fout, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			if _, err := f.Write(df.Data); err != nil {
				log.Fatal(err)
			}
			if err := f.Close(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
