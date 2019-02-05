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

package zipkin_test

import (
	"log"

	openzipkin "github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/exporter/zipkin"
	"go.opencensus.io/trace"
)

func Example() {
	// import (
	//     openzipkin "github.com/openzipkin/zipkin-go"
	//     "github.com/openzipkin/zipkin-go/reporter/http"
	//     "go.opencensus.io/exporter/trace/zipkin"
	// )

	localEndpoint, err := openzipkin.NewEndpoint("server", "192.168.1.5:5454")
	if err != nil {
		log.Print(err)
	}
	reporter := http.NewReporter("http://localhost:9411/api/v2/spans")
	exporter := zipkin.NewExporter(reporter, localEndpoint)
	trace.RegisterExporter(exporter)
}
