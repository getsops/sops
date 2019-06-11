/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

/*
Package main provides benchmark with setting flags.

An example to run some benchmarks with profiling enabled:

go run benchmark/benchmain/main.go -benchtime=10s -workloads=all \
  -compression=gzip -maxConcurrentCalls=1 -trace=off \
  -reqSizeBytes=1,1048576 -respSizeBytes=1,1048576 -networkMode=Local \
  -cpuProfile=cpuProf -memProfile=memProf -memProfileRate=10000 -resultFile=result

As a suggestion, when creating a branch, you can run this benchmark and save the result
file "-resultFile=basePerf", and later when you at the middle of the work or finish the
work, you can get the benchmark result and compare it with the base anytime.

Assume there are two result files names as "basePerf" and "curPerf" created by adding
-resultFile=basePerf and -resultFile=curPerf.
	To format the curPerf, run:
  	go run benchmark/benchresult/main.go curPerf
	To observe how the performance changes based on a base result, run:
  	go run benchmark/benchresult/main.go basePerf curPerf
*/
package main

import (
	"context"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"google.golang.org/grpc"
	bm "google.golang.org/grpc/benchmark"
	"google.golang.org/grpc/benchmark/flags"
	testpb "google.golang.org/grpc/benchmark/grpc_testing"
	"google.golang.org/grpc/benchmark/latency"
	"google.golang.org/grpc/benchmark/stats"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/channelz"
	"google.golang.org/grpc/test/bufconn"
)

var (
	workloads = flags.StringWithAllowedValues("workloads", workloadsAll,
		fmt.Sprintf("Workloads to execute - One of: %v", strings.Join(allWorkloads, ", ")), allWorkloads)
	traceMode = flags.StringWithAllowedValues("trace", toggleModeOff,
		fmt.Sprintf("Trace mode - One of: %v", strings.Join(allToggleModes, ", ")), allToggleModes)
	preloaderMode = flags.StringWithAllowedValues("preloader", toggleModeOff,
		fmt.Sprintf("Preloader mode - One of: %v", strings.Join(allToggleModes, ", ")), allToggleModes)
	channelzOn = flags.StringWithAllowedValues("channelz", toggleModeOff,
		fmt.Sprintf("Channelz mode - One of: %v", strings.Join(allToggleModes, ", ")), allToggleModes)
	compressorMode = flags.StringWithAllowedValues("compression", compModeOff,
		fmt.Sprintf("Compression mode - One of: %v", strings.Join(allCompModes, ", ")), allCompModes)
	networkMode = flags.StringWithAllowedValues("networkMode", networkModeNone,
		"Network mode includes LAN, WAN, Local and Longhaul", allNetworkModes)
	readLatency        = flags.DurationSlice("latency", defaultReadLatency, "Simulated one-way network latency - may be a comma-separated list")
	readKbps           = flags.IntSlice("kbps", defaultReadKbps, "Simulated network throughput (in kbps) - may be a comma-separated list")
	readMTU            = flags.IntSlice("mtu", defaultReadMTU, "Simulated network MTU (Maximum Transmission Unit) - may be a comma-separated list")
	maxConcurrentCalls = flags.IntSlice("maxConcurrentCalls", defaultMaxConcurrentCalls, "Number of concurrent RPCs during benchmarks")
	readReqSizeBytes   = flags.IntSlice("reqSizeBytes", defaultReqSizeBytes, "Request size in bytes - may be a comma-separated list")
	readRespSizeBytes  = flags.IntSlice("respSizeBytes", defaultRespSizeBytes, "Response size in bytes - may be a comma-separated list")
	benchTime          = flag.Duration("benchtime", time.Second, "Configures the amount of time to run each benchmark")
	memProfile         = flag.String("memProfile", "", "Enables memory profiling output to the filename provided.")
	memProfileRate     = flag.Int("memProfileRate", 512*1024, "Configures the memory profiling rate. \n"+
		"memProfile should be set before setting profile rate. To include every allocated block in the profile, "+
		"set MemProfileRate to 1. To turn off profiling entirely, set MemProfileRate to 0. 512 * 1024 by default.")
	cpuProfile          = flag.String("cpuProfile", "", "Enables CPU profiling output to the filename provided")
	benchmarkResultFile = flag.String("resultFile", "", "Save the benchmark result into a binary file")
	useBufconn          = flag.Bool("bufconn", false, "Use in-memory connection instead of system network I/O")
)

const (
	workloadsUnary         = "unary"
	workloadsStreaming     = "streaming"
	workloadsUnconstrained = "unconstrained"
	workloadsAll           = "all"
	// Compression modes.
	compModeOff  = "off"
	compModeGzip = "gzip"
	compModeNop  = "nop"
	compModeAll  = "all"
	// Toggle modes.
	toggleModeOff  = "off"
	toggleModeOn   = "on"
	toggleModeBoth = "both"
	// Network modes.
	networkModeNone  = "none"
	networkModeLocal = "Local"
	networkModeLAN   = "LAN"
	networkModeWAN   = "WAN"
	networkLongHaul  = "Longhaul"

	numStatsBuckets = 10
)

var (
	allWorkloads              = []string{workloadsUnary, workloadsStreaming, workloadsUnconstrained, workloadsAll}
	allCompModes              = []string{compModeOff, compModeGzip, compModeNop, compModeAll}
	allToggleModes            = []string{toggleModeOff, toggleModeOn, toggleModeBoth}
	allNetworkModes           = []string{networkModeNone, networkModeLocal, networkModeLAN, networkModeWAN, networkLongHaul}
	defaultReadLatency        = []time.Duration{0, 40 * time.Millisecond} // if non-positive, no delay.
	defaultReadKbps           = []int{0, 10240}                           // if non-positive, infinite
	defaultReadMTU            = []int{0}                                  // if non-positive, infinite
	defaultMaxConcurrentCalls = []int{1, 8, 64, 512}
	defaultReqSizeBytes       = []int{1, 1024, 1024 * 1024}
	defaultRespSizeBytes      = []int{1, 1024, 1024 * 1024}
	networks                  = map[string]latency.Network{
		networkModeLocal: latency.Local,
		networkModeLAN:   latency.LAN,
		networkModeWAN:   latency.WAN,
		networkLongHaul:  latency.Longhaul,
	}
)

// runModes indicates the workloads to run. This is initialized with a call to
// `runModesFromWorkloads`, passing the workloads flag set by the user.
type runModes struct {
	unary, streaming, unconstrained bool
}

// runModesFromWorkloads determines the runModes based on the value of
// workloads flag set by the user.
func runModesFromWorkloads(workload string) runModes {
	r := runModes{}
	switch workload {
	case workloadsUnary:
		r.unary = true
	case workloadsStreaming:
		r.streaming = true
	case workloadsUnconstrained:
		r.unconstrained = true
	case workloadsAll:
		r.unary = true
		r.streaming = true
		r.unconstrained = true
	default:
		log.Fatalf("Unknown workloads setting: %v (want one of: %v)",
			workloads, strings.Join(allWorkloads, ", "))
	}
	return r
}

func unaryBenchmark(startTimer func(), stopTimer func(uint64), benchFeatures stats.Features, benchTime time.Duration, s *stats.Stats) uint64 {
	caller, cleanup := makeFuncUnary(benchFeatures)
	defer cleanup()
	return runBenchmark(caller, startTimer, stopTimer, benchFeatures, benchTime, s)
}

func streamBenchmark(startTimer func(), stopTimer func(uint64), benchFeatures stats.Features, benchTime time.Duration, s *stats.Stats) uint64 {
	caller, cleanup := makeFuncStream(benchFeatures)
	defer cleanup()
	return runBenchmark(caller, startTimer, stopTimer, benchFeatures, benchTime, s)
}

func unconstrainedStreamBenchmark(benchFeatures stats.Features, warmuptime, benchTime time.Duration) (uint64, uint64) {
	var sender, recver func(int)
	var cleanup func()
	if benchFeatures.EnablePreloader {
		sender, recver, cleanup = makeFuncUnconstrainedStreamPreloaded(benchFeatures)
	} else {
		sender, recver, cleanup = makeFuncUnconstrainedStream(benchFeatures)
	}
	defer cleanup()

	var (
		wg            sync.WaitGroup
		requestCount  uint64
		responseCount uint64
	)
	wg.Add(2 * benchFeatures.MaxConcurrentCalls)

	// Resets the counters once warmed up
	go func() {
		<-time.NewTimer(warmuptime).C
		atomic.StoreUint64(&requestCount, 0)
		atomic.StoreUint64(&responseCount, 0)
	}()

	bmEnd := time.Now().Add(benchTime + warmuptime)
	for i := 0; i < benchFeatures.MaxConcurrentCalls; i++ {
		go func(pos int) {
			for {
				t := time.Now()
				if t.After(bmEnd) {
					break
				}
				sender(pos)
				atomic.AddUint64(&requestCount, 1)
			}
			wg.Done()
		}(i)
		go func(pos int) {
			for {
				t := time.Now()
				if t.After(bmEnd) {
					break
				}
				recver(pos)
				atomic.AddUint64(&responseCount, 1)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	return requestCount, responseCount
}

func makeClient(benchFeatures stats.Features) (testpb.BenchmarkServiceClient, func()) {
	nw := &latency.Network{Kbps: benchFeatures.Kbps, Latency: benchFeatures.Latency, MTU: benchFeatures.Mtu}
	opts := []grpc.DialOption{}
	sopts := []grpc.ServerOption{}
	if benchFeatures.ModeCompressor == compModeNop {
		sopts = append(sopts,
			grpc.RPCCompressor(nopCompressor{}),
			grpc.RPCDecompressor(nopDecompressor{}),
		)
		opts = append(opts,
			grpc.WithCompressor(nopCompressor{}),
			grpc.WithDecompressor(nopDecompressor{}),
		)
	}
	if benchFeatures.ModeCompressor == compModeGzip {
		sopts = append(sopts,
			grpc.RPCCompressor(grpc.NewGZIPCompressor()),
			grpc.RPCDecompressor(grpc.NewGZIPDecompressor()),
		)
		opts = append(opts,
			grpc.WithCompressor(grpc.NewGZIPCompressor()),
			grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
		)
	}
	sopts = append(sopts, grpc.MaxConcurrentStreams(uint32(benchFeatures.MaxConcurrentCalls+1)))
	opts = append(opts, grpc.WithInsecure())

	var lis net.Listener
	if benchFeatures.UseBufConn {
		bcLis := bufconn.Listen(256 * 1024)
		lis = bcLis
		opts = append(opts, grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
			return nw.ContextDialer(func(context.Context, string, string) (net.Conn, error) {
				return bcLis.Dial()
			})(ctx, "", "")
		}))
	} else {
		var err error
		lis, err = net.Listen("tcp", "localhost:0")
		if err != nil {
			grpclog.Fatalf("Failed to listen: %v", err)
		}
		opts = append(opts, grpc.WithContextDialer(func(ctx context.Context, address string) (net.Conn, error) {
			return nw.ContextDialer((&net.Dialer{}).DialContext)(ctx, "tcp", lis.Addr().String())
		}))
	}
	lis = nw.Listener(lis)
	stopper := bm.StartServer(bm.ServerInfo{Type: "protobuf", Listener: lis}, sopts...)
	conn := bm.NewClientConn("" /* target not used */, opts...)
	return testpb.NewBenchmarkServiceClient(conn), func() {
		conn.Close()
		stopper()
	}
}

func makeFuncUnary(benchFeatures stats.Features) (func(int), func()) {
	tc, cleanup := makeClient(benchFeatures)
	return func(int) {
		unaryCaller(tc, benchFeatures.ReqSizeBytes, benchFeatures.RespSizeBytes)
	}, cleanup
}

func makeFuncStream(benchFeatures stats.Features) (func(int), func()) {
	tc, cleanup := makeClient(benchFeatures)

	streams := make([]testpb.BenchmarkService_StreamingCallClient, benchFeatures.MaxConcurrentCalls)
	for i := 0; i < benchFeatures.MaxConcurrentCalls; i++ {
		stream, err := tc.StreamingCall(context.Background())
		if err != nil {
			grpclog.Fatalf("%v.StreamingCall(_) = _, %v", tc, err)
		}
		streams[i] = stream
	}

	return func(pos int) {
		streamCaller(streams[pos], benchFeatures.ReqSizeBytes, benchFeatures.RespSizeBytes)
	}, cleanup
}

func makeFuncUnconstrainedStreamPreloaded(benchFeatures stats.Features) (func(int), func(int), func()) {
	streams, req, cleanup := setupUnconstrainedStream(benchFeatures)

	preparedMsg := make([]*grpc.PreparedMsg, len(streams))
	for i, stream := range streams {
		preparedMsg[i] = &grpc.PreparedMsg{}
		err := preparedMsg[i].Encode(stream, req)
		if err != nil {
			grpclog.Fatalf("%v.Encode(%v, %v) = %v", preparedMsg[i], req, stream, err)
		}
	}

	return func(pos int) {
			streams[pos].SendMsg(preparedMsg[pos])
		}, func(pos int) {
			streams[pos].Recv()
		}, cleanup
}

func makeFuncUnconstrainedStream(benchFeatures stats.Features) (func(int), func(int), func()) {
	streams, req, cleanup := setupUnconstrainedStream(benchFeatures)

	return func(pos int) {
			streams[pos].Send(req)
		}, func(pos int) {
			streams[pos].Recv()
		}, cleanup
}

func setupUnconstrainedStream(benchFeatures stats.Features) ([]testpb.BenchmarkService_StreamingCallClient, *testpb.SimpleRequest, func()) {
	tc, cleanup := makeClient(benchFeatures)

	streams := make([]testpb.BenchmarkService_StreamingCallClient, benchFeatures.MaxConcurrentCalls)
	for i := 0; i < benchFeatures.MaxConcurrentCalls; i++ {
		stream, err := tc.UnconstrainedStreamingCall(context.Background())
		if err != nil {
			grpclog.Fatalf("%v.UnconstrainedStreamingCall(_) = _, %v", tc, err)
		}
		streams[i] = stream
	}

	pl := bm.NewPayload(testpb.PayloadType_COMPRESSABLE, benchFeatures.ReqSizeBytes)
	req := &testpb.SimpleRequest{
		ResponseType: pl.Type,
		ResponseSize: int32(benchFeatures.RespSizeBytes),
		Payload:      pl,
	}

	return streams, req, cleanup
}

func unaryCaller(client testpb.BenchmarkServiceClient, reqSize, respSize int) {
	if err := bm.DoUnaryCall(client, reqSize, respSize); err != nil {
		grpclog.Fatalf("DoUnaryCall failed: %v", err)
	}
}

func streamCaller(stream testpb.BenchmarkService_StreamingCallClient, reqSize, respSize int) {
	if err := bm.DoStreamingRoundTrip(stream, reqSize, respSize); err != nil {
		grpclog.Fatalf("DoStreamingRoundTrip failed: %v", err)
	}
}

func runBenchmark(caller func(int), startTimer func(), stopTimer func(uint64), benchFeatures stats.Features, benchTime time.Duration, s *stats.Stats) uint64 {
	// Warm up connection.
	for i := 0; i < 10; i++ {
		caller(0)
	}
	// Run benchmark.
	startTimer()
	var (
		mu sync.Mutex
		wg sync.WaitGroup
	)
	wg.Add(benchFeatures.MaxConcurrentCalls)
	bmEnd := time.Now().Add(benchTime)
	var count uint64
	for i := 0; i < benchFeatures.MaxConcurrentCalls; i++ {
		go func(pos int) {
			for {
				t := time.Now()
				if t.After(bmEnd) {
					break
				}
				start := time.Now()
				caller(pos)
				elapse := time.Since(start)
				atomic.AddUint64(&count, 1)
				mu.Lock()
				s.Add(elapse)
				mu.Unlock()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	stopTimer(count)
	return count
}

// benchOpts represents all configurable options available while running this
// benchmark. This is built from the values passed as flags.
type benchOpts struct {
	rModes              runModes
	benchTime           time.Duration
	memProfileRate      int
	memProfile          string
	cpuProfile          string
	networkMode         string
	benchmarkResultFile string
	useBufconn          bool
	features            *featureOpts
}

// featureOpts represents options which can have multiple values. The user
// usually provides a comma-separated list of options for each of these
// features through command line flags. We generate all possible combinations
// for the provided values and run the benchmarks for each combination.
type featureOpts struct {
	enableTrace        []bool          // Feature index 0
	readLatencies      []time.Duration // Feature index 1
	readKbps           []int           // Feature index 2
	readMTU            []int           // Feature index 3
	maxConcurrentCalls []int           // Feature index 4
	reqSizeBytes       []int           // Feature index 5
	respSizeBytes      []int           // Feature index 6
	compModes          []string        // Feature index 7
	enableChannelz     []bool          // Feature index 8
	enablePreloader    []bool          // Feature index 9
}

// featureIndex is an enum for the different features that could be configured
// by the user through command line flags.
type featureIndex int

const (
	enableTraceIndex featureIndex = iota
	readLatenciesIndex
	readKbpsIndex
	readMTUIndex
	maxConcurrentCallsIndex
	reqSizeBytesIndex
	respSizeBytesIndex
	compModesIndex
	enableChannelzIndex
	enablePreloaderIndex

	// This is a place holder to indicate the total number of feature indices we
	// have. Any new feature indices should be added above this.
	maxFeatureIndex
)

// makeFeaturesNum returns a slice of ints of size 'maxFeatureIndex' where each
// element of the slice (indexed by 'featuresIndex' enum) contains the number
// of features to be exercised by the benchmark code.
// For example: Index 0 of the returned slice contains the number of values for
// enableTrace feature, while index 1 contains the number of value of
// readLatencies feature and so on.
func makeFeaturesNum(b *benchOpts) []int {
	featuresNum := make([]int, maxFeatureIndex)
	for i := 0; i < len(featuresNum); i++ {
		switch featureIndex(i) {
		case enableTraceIndex:
			featuresNum[i] = len(b.features.enableTrace)
		case readLatenciesIndex:
			featuresNum[i] = len(b.features.readLatencies)
		case readKbpsIndex:
			featuresNum[i] = len(b.features.readKbps)
		case readMTUIndex:
			featuresNum[i] = len(b.features.readMTU)
		case maxConcurrentCallsIndex:
			featuresNum[i] = len(b.features.maxConcurrentCalls)
		case reqSizeBytesIndex:
			featuresNum[i] = len(b.features.reqSizeBytes)
		case respSizeBytesIndex:
			featuresNum[i] = len(b.features.respSizeBytes)
		case compModesIndex:
			featuresNum[i] = len(b.features.compModes)
		case enableChannelzIndex:
			featuresNum[i] = len(b.features.enableChannelz)
		case enablePreloaderIndex:
			featuresNum[i] = len(b.features.enablePreloader)
		default:
			log.Fatalf("Unknown feature index %v in generateFeatures. maxFeatureIndex is %v", i, maxFeatureIndex)
		}
	}
	return featuresNum
}

// sharedFeatures returns a bool slice which acts as a bitmask. Each item in
// the slice represents a feature, indexed by 'featureIndex' enum.  The bit is
// set to 1 if the corresponding feature does not have multiple value, so is
// shared amongst all benchmarks.
func sharedFeatures(featuresNum []int) []bool {
	result := make([]bool, len(featuresNum))
	for i, num := range featuresNum {
		if num <= 1 {
			result[i] = true
		}
	}
	return result
}

// generateFeatures generates all combinations of the provided feature options.
// While all the feature options are stored in the benchOpts struct, the input
// parameter 'featuresNum' is a slice indexed by 'featureIndex' enum containing
// the number of values for each feature.
// For example, let's say the user sets -workloads=all and
// -maxConcurrentCalls=1,100, this would end up with the following
// combinations:
// [workloads: unary, maxConcurrentCalls=1]
// [workloads: unary, maxConcurrentCalls=1]
// [workloads: streaming, maxConcurrentCalls=100]
// [workloads: streaming, maxConcurrentCalls=100]
// [workloads: unconstrained, maxConcurrentCalls=1]
// [workloads: unconstrained, maxConcurrentCalls=100]
func (b *benchOpts) generateFeatures(featuresNum []int) []stats.Features {
	// curPos and initialPos are two slices where each value acts as an index
	// into the appropriate feature slice maintained in benchOpts.features. This
	// loop generates all possible combinations of features by changing one value
	// at a time, and once curPos becomes equal to initialPos, we have explored
	// all options.
	var result []stats.Features
	var curPos []int
	initialPos := make([]int, maxFeatureIndex)
	for !reflect.DeepEqual(initialPos, curPos) {
		if curPos == nil {
			curPos = make([]int, maxFeatureIndex)
		}
		result = append(result, stats.Features{
			// These features stay the same for each iteration.
			NetworkMode: b.networkMode,
			UseBufConn:  b.useBufconn,
			// These features can potentially change for each iteration.
			EnableTrace:        b.features.enableTrace[curPos[enableTraceIndex]],
			Latency:            b.features.readLatencies[curPos[readLatenciesIndex]],
			Kbps:               b.features.readKbps[curPos[readKbpsIndex]],
			Mtu:                b.features.readMTU[curPos[readMTUIndex]],
			MaxConcurrentCalls: b.features.maxConcurrentCalls[curPos[maxConcurrentCallsIndex]],
			ReqSizeBytes:       b.features.reqSizeBytes[curPos[reqSizeBytesIndex]],
			RespSizeBytes:      b.features.respSizeBytes[curPos[respSizeBytesIndex]],
			ModeCompressor:     b.features.compModes[curPos[compModesIndex]],
			EnableChannelz:     b.features.enableChannelz[curPos[enableChannelzIndex]],
			EnablePreloader:    b.features.enablePreloader[curPos[enablePreloaderIndex]],
		})
		addOne(curPos, featuresNum)
	}
	return result
}

// addOne mutates the input slice 'features' by changing one feature, thus
// arriving at the next combination of feature values. 'featuresMaxPosition'
// provides the numbers of allowed values for each feature, indexed by
// 'featureIndex' enum.
func addOne(features []int, featuresMaxPosition []int) {
	for i := len(features) - 1; i >= 0; i-- {
		features[i] = (features[i] + 1)
		if features[i]/featuresMaxPosition[i] == 0 {
			break
		}
		features[i] = features[i] % featuresMaxPosition[i]
	}
}

// processFlags reads the command line flags and builds benchOpts. Specifying
// invalid values for certain flags will cause flag.Parse() to fail, and the
// program to terminate.
// This *SHOULD* be the only place where the flags are accessed. All other
// parts of the benchmark code should rely on the returned benchOpts.
func processFlags() *benchOpts {
	flag.Parse()
	if flag.NArg() != 0 {
		log.Fatal("Error: unparsed arguments: ", flag.Args())
	}

	opts := &benchOpts{
		rModes:              runModesFromWorkloads(*workloads),
		benchTime:           *benchTime,
		memProfileRate:      *memProfileRate,
		memProfile:          *memProfile,
		cpuProfile:          *cpuProfile,
		networkMode:         *networkMode,
		benchmarkResultFile: *benchmarkResultFile,
		useBufconn:          *useBufconn,
		features: &featureOpts{
			enableTrace:        setToggleMode(*traceMode),
			readLatencies:      append([]time.Duration(nil), *readLatency...),
			readKbps:           append([]int(nil), *readKbps...),
			readMTU:            append([]int(nil), *readMTU...),
			maxConcurrentCalls: append([]int(nil), *maxConcurrentCalls...),
			reqSizeBytes:       append([]int(nil), *readReqSizeBytes...),
			respSizeBytes:      append([]int(nil), *readRespSizeBytes...),
			compModes:          setCompressorMode(*compressorMode),
			enableChannelz:     setToggleMode(*channelzOn),
			enablePreloader:    setToggleMode(*preloaderMode),
		},
	}

	// Re-write latency, kpbs and mtu if network mode is set.
	if network, ok := networks[opts.networkMode]; ok {
		opts.features.readLatencies = []time.Duration{network.Latency}
		opts.features.readKbps = []int{network.Kbps}
		opts.features.readMTU = []int{network.MTU}
	}
	return opts
}

func setToggleMode(val string) []bool {
	switch val {
	case toggleModeOn:
		return []bool{true}
	case toggleModeOff:
		return []bool{false}
	case toggleModeBoth:
		return []bool{false, true}
	default:
		// This should never happen because a wrong value passed to this flag would
		// be caught during flag.Parse().
		return []bool{}
	}
}

func setCompressorMode(val string) []string {
	switch val {
	case compModeNop, compModeGzip, compModeOff:
		return []string{val}
	case compModeAll:
		return []string{compModeNop, compModeGzip, compModeOff}
	default:
		// This should never happen because a wrong value passed to this flag would
		// be caught during flag.Parse().
		return []string{}
	}
}

func printThroughput(requestCount uint64, requestSize int, responseCount uint64, responseSize int, benchTime time.Duration) {
	requestThroughput := float64(requestCount) * float64(requestSize) * 8 / benchTime.Seconds()
	responseThroughput := float64(responseCount) * float64(responseSize) * 8 / benchTime.Seconds()
	fmt.Printf("Number of requests:  %v\tRequest throughput:  %v bit/s\n", requestCount, requestThroughput)
	fmt.Printf("Number of responses: %v\tResponse throughput: %v bit/s\n", responseCount, responseThroughput)
	fmt.Println()
}

func main() {
	opts := processFlags()
	before(opts)
	s := stats.NewStats(numStatsBuckets)
	s.SortLatency()
	var memStats runtime.MemStats
	var results testing.BenchmarkResult
	var startAllocs, startBytes uint64
	var startTime time.Time
	var startTimer = func() {
		runtime.ReadMemStats(&memStats)
		startAllocs = memStats.Mallocs
		startBytes = memStats.TotalAlloc
		startTime = time.Now()
	}
	var stopTimer = func(count uint64) {
		runtime.ReadMemStats(&memStats)
		results = testing.BenchmarkResult{
			N:         int(count),
			T:         time.Since(startTime),
			Bytes:     0,
			MemAllocs: memStats.Mallocs - startAllocs,
			MemBytes:  memStats.TotalAlloc - startBytes,
		}
	}

	// Run benchmarks
	resultSlice := []stats.BenchResults{}
	featuresNum := makeFeaturesNum(opts)
	sharedPos := sharedFeatures(featuresNum)
	for _, benchFeature := range opts.generateFeatures(featuresNum) {
		grpc.EnableTracing = benchFeature.EnableTrace
		if benchFeature.EnableChannelz {
			channelz.TurnOn()
		}
		if opts.rModes.unary {
			count := unaryBenchmark(startTimer, stopTimer, benchFeature, opts.benchTime, s)
			s.SetBenchmarkResult("Unary", benchFeature, results.N,
				results.AllocedBytesPerOp(), results.AllocsPerOp(), sharedPos)
			fmt.Println(s.BenchString())
			fmt.Println(s.String())
			printThroughput(count, benchFeature.ReqSizeBytes, count, benchFeature.RespSizeBytes, opts.benchTime)
			resultSlice = append(resultSlice, s.GetBenchmarkResults())
			s.Clear()
		}
		if opts.rModes.streaming {
			count := streamBenchmark(startTimer, stopTimer, benchFeature, opts.benchTime, s)
			s.SetBenchmarkResult("Stream", benchFeature, results.N,
				results.AllocedBytesPerOp(), results.AllocsPerOp(), sharedPos)
			fmt.Println(s.BenchString())
			fmt.Println(s.String())
			printThroughput(count, benchFeature.ReqSizeBytes, count, benchFeature.RespSizeBytes, opts.benchTime)
			resultSlice = append(resultSlice, s.GetBenchmarkResults())
			s.Clear()
		}
		if opts.rModes.unconstrained {
			requestCount, responseCount := unconstrainedStreamBenchmark(benchFeature, time.Second, opts.benchTime)
			fmt.Printf("Unconstrained Stream-%v\n", benchFeature)
			printThroughput(requestCount, benchFeature.ReqSizeBytes, responseCount, benchFeature.RespSizeBytes, opts.benchTime)
		}
	}
	after(opts, resultSlice)
}

func before(opts *benchOpts) {
	if opts.memProfile != "" {
		runtime.MemProfileRate = opts.memProfileRate
	}
	if opts.cpuProfile != "" {
		f, err := os.Create(opts.cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			return
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't start cpu profile: %s\n", err)
			f.Close()
			return
		}
	}
}

func after(opts *benchOpts, data []stats.BenchResults) {
	if opts.cpuProfile != "" {
		pprof.StopCPUProfile() // flushes profile to disk
	}
	if opts.memProfile != "" {
		f, err := os.Create(opts.memProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			os.Exit(2)
		}
		runtime.GC() // materialize all statistics
		if err = pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't write heap profile %s: %s\n", opts.memProfile, err)
			os.Exit(2)
		}
		f.Close()
	}
	if opts.benchmarkResultFile != "" {
		f, err := os.Create(opts.benchmarkResultFile)
		if err != nil {
			log.Fatalf("testing: can't write benchmark result %s: %s\n", opts.benchmarkResultFile, err)
		}
		dataEncoder := gob.NewEncoder(f)
		dataEncoder.Encode(data)
		f.Close()
	}
}

// nopCompressor is a compressor that just copies data.
type nopCompressor struct{}

func (nopCompressor) Do(w io.Writer, p []byte) error {
	n, err := w.Write(p)
	if err != nil {
		return err
	}
	if n != len(p) {
		return fmt.Errorf("nopCompressor.Write: wrote %v bytes; want %v", n, len(p))
	}
	return nil
}

func (nopCompressor) Type() string { return compModeNop }

// nopDecompressor is a decompressor that just copies data.
type nopDecompressor struct{}

func (nopDecompressor) Do(r io.Reader) ([]byte, error) { return ioutil.ReadAll(r) }
func (nopDecompressor) Type() string                   { return compModeNop }
