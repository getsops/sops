// Copyright 2018, OpenCensus Authors
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

package jaeger_test

import (
	"log"

	"go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

func ExampleNewExporter_collector() {
	// Register the Jaeger exporter to be able to retrieve
	// the collected spans.
	exporter, err := jaeger.NewExporter(jaeger.Options{
		Endpoint: "http://localhost:14268",
		Process: jaeger.Process{
			ServiceName: "trace-demo",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
}

func ExampleNewExporter_agent() {
	// Register the Jaeger exporter to be able to retrieve
	// the collected spans.
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "localhost:6831",
		Process: jaeger.Process{
			ServiceName: "trace-demo",
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
}

// ExampleNewExporter_processTags shows how to set ProcessTags
// on a Jaeger exporter. These tags will be added to the exported
// Jaeger process.
func ExampleNewExporter_processTags() {
	// Register the Jaeger exporter to be able to retrieve
	// the collected spans.
	exporter, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint: "localhost:6831",
		Process: jaeger.Process{
			ServiceName: "trace-demo",
			Tags: []jaeger.Tag{
				jaeger.StringTag("ip", "127.0.0.1"),
				jaeger.BoolTag("demo", true),
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	trace.RegisterExporter(exporter)
}
