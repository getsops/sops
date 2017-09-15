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
  -compression=on -maxConcurrentCalls=1 -traceMode=false \
  -reqSizeBytes=1,1048576 -respSizeBytes=1,1048576 \
  -latency=0s -kbps=0 -mtu=0 \
  -cpuProfile=cpuProf -memProfile=memProf -memProfileRate=10000
*/
package main

import (
	"errors"
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
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	bm "google.golang.org/grpc/benchmark"
	testpb "google.golang.org/grpc/benchmark/grpc_testing"
	"google.golang.org/grpc/benchmark/latency"
	"google.golang.org/grpc/benchmark/stats"
	"google.golang.org/grpc/grpclog"
)

const (
	compressionOn   = "on"
	compressionOff  = "off"
	compressionBoth = "both"
)

var allCompressionModes = []string{compressionOn, compressionOff, compressionBoth}

const (
	workloadsUnary     = "unary"
	workloadsStreaming = "streaming"
	workloadsAll       = "all"
)

var allWorkloads = []string{workloadsUnary, workloadsStreaming, workloadsAll}

var (
	runMode = []bool{true, true} // {runUnary, runStream}
	// When set the latency to 0 (no delay), the result is slower than the real result with no delay
	// because latency simulation section has extra operations
	ltc                    = []time.Duration{0, 40 * time.Millisecond} // if non-positive, no delay.
	kbps                   = []int{0, 10240}                           // if non-positive, infinite
	mtu                    = []int{0}                                  // if non-positive, infinite
	maxConcurrentCalls     = []int{1, 8, 64, 512}
	reqSizeBytes           = []int{1, 1024, 1024 * 1024}
	respSizeBytes          = []int{1, 1024, 1024 * 1024}
	enableTrace            = []bool{false}
	benchtime              time.Duration
	memProfile, cpuProfile string
	memProfileRate         int
	enableCompressor       []bool
)

func unaryBenchmark(startTimer func(), stopTimer func(int32), benchFeatures bm.Features, benchtime time.Duration, s *stats.Stats) {
	caller, close := makeFuncUnary(benchFeatures)
	defer close()
	runBenchmark(caller, startTimer, stopTimer, benchFeatures, benchtime, s)
}

func streamBenchmark(startTimer func(), stopTimer func(int32), benchFeatures bm.Features, benchtime time.Duration, s *stats.Stats) {
	caller, close := makeFuncStream(benchFeatures)
	defer close()
	runBenchmark(caller, startTimer, stopTimer, benchFeatures, benchtime, s)
}

func makeFuncUnary(benchFeatures bm.Features) (func(int), func()) {
	nw := &latency.Network{Kbps: benchFeatures.Kbps, Latency: benchFeatures.Latency, MTU: benchFeatures.Mtu}
	opts := []grpc.DialOption{}
	sopts := []grpc.ServerOption{}
	if benchFeatures.EnableCompressor {
		sopts = append(sopts,
			grpc.RPCCompressor(nopCompressor{}),
			grpc.RPCDecompressor(nopDecompressor{}),
		)
		opts = append(opts,
			grpc.WithCompressor(nopCompressor{}),
			grpc.WithDecompressor(nopDecompressor{}),
		)
	}
	sopts = append(sopts, grpc.MaxConcurrentStreams(uint32(benchFeatures.MaxConcurrentCalls+1)))
	opts = append(opts, grpc.WithDialer(func(address string, timeout time.Duration) (net.Conn, error) {
		return nw.TimeoutDialer(net.DialTimeout)("tcp", address, timeout)
	}))
	opts = append(opts, grpc.WithInsecure())

	target, stopper := bm.StartServer(bm.ServerInfo{Addr: "localhost:0", Type: "protobuf", Network: nw}, sopts...)
	conn := bm.NewClientConn(target, opts...)
	tc := testpb.NewBenchmarkServiceClient(conn)
	return func(int) {
			unaryCaller(tc, benchFeatures.ReqSizeBytes, benchFeatures.RespSizeBytes)
		}, func() {
			conn.Close()
			stopper()
		}
}

func makeFuncStream(benchFeatures bm.Features) (func(int), func()) {
	fmt.Println(benchFeatures)
	nw := &latency.Network{Kbps: benchFeatures.Kbps, Latency: benchFeatures.Latency, MTU: benchFeatures.Mtu}
	opts := []grpc.DialOption{}
	sopts := []grpc.ServerOption{}
	if benchFeatures.EnableCompressor {
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
	opts = append(opts, grpc.WithDialer(func(address string, timeout time.Duration) (net.Conn, error) {
		return nw.TimeoutDialer(net.DialTimeout)("tcp", address, timeout)
	}))
	opts = append(opts, grpc.WithInsecure())

	target, stopper := bm.StartServer(bm.ServerInfo{Addr: "localhost:0", Type: "protobuf", Network: nw}, sopts...)
	conn := bm.NewClientConn(target, opts...)
	tc := testpb.NewBenchmarkServiceClient(conn)
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
		}, func() {
			conn.Close()
			stopper()
		}
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

func runBenchmark(caller func(int), startTimer func(), stopTimer func(int32), benchFeatures bm.Features, benchtime time.Duration, s *stats.Stats) {
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
	bmEnd := time.Now().Add(benchtime)
	var count int32
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
				atomic.AddInt32(&count, 1)
				mu.Lock()
				s.Add(elapse)
				mu.Unlock()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	stopTimer(count)
}

// Initiate main function to get settings of features.
func init() {
	var (
		workloads, compressorMode, readLatency    string
		readKbps, readMtu, readMaxConcurrentCalls intSliceType
		readReqSizeBytes, readRespSizeBytes       intSliceType
		traceMode                                 bool
	)
	flag.StringVar(&workloads, "workloads", workloadsAll,
		fmt.Sprintf("Workloads to execute - One of: %v", strings.Join(allWorkloads, ", ")))
	flag.BoolVar(&traceMode, "traceMode", false, "Enable gRPC tracing")
	flag.StringVar(&readLatency, "latency", "", "Simulated one-way network latency - may be a comma-separated list")
	flag.DurationVar(&benchtime, "benchtime", time.Second, "Configures the amount of time to run each benchmark")
	flag.Var(&readKbps, "kbps", "Simulated network throughput (in kbps) - may be a comma-separated list")
	flag.Var(&readMtu, "mtu", "Simulated network MTU (Maximum Transmission Unit) - may be a comma-separated list")
	flag.Var(&readMaxConcurrentCalls, "maxConcurrentCalls", "Number of concurrent RPCs during benchmarks")
	flag.Var(&readReqSizeBytes, "reqSizeBytes", "Request size in bytes - may be a comma-separated list")
	flag.Var(&readRespSizeBytes, "respSizeBytes", "Response size in bytes - may be a comma-separated list")
	flag.StringVar(&memProfile, "memProfile", "", "Enables memory profiling output to the filename provided")
	flag.IntVar(&memProfileRate, "memProfileRate", 0, "Configures the memory profiling rate")
	flag.StringVar(&cpuProfile, "cpuProfile", "", "Enables CPU profiling output to the filename provided")
	flag.StringVar(&compressorMode, "compression", compressionOff,
		fmt.Sprintf("Compression mode - One of: %v", strings.Join(allCompressionModes, ", ")))
	flag.Parse()
	if flag.NArg() != 0 {
		log.Fatal("Error: unparsed arguments: ", flag.Args())
	}
	switch workloads {
	case workloadsUnary:
		runMode[0] = true
		runMode[1] = false
	case workloadsStreaming:
		runMode[0] = false
		runMode[1] = true
	case workloadsAll:
		runMode[0] = true
		runMode[1] = true
	default:
		log.Fatalf("Unknown workloads setting: %v (want one of: %v)",
			workloads, strings.Join(allWorkloads, ", "))
	}
	switch compressorMode {
	case compressionOn:
		enableCompressor = []bool{true}
	case compressionOff:
		enableCompressor = []bool{false}
	case compressionBoth:
		enableCompressor = []bool{false, true}
	default:
		log.Fatalf("Unknown compression mode setting: %v (want one of: %v)",
			compressorMode, strings.Join(allCompressionModes, ", "))
	}
	if traceMode {
		enableTrace = []bool{true}
	}
	// Time input formats as (time + unit).
	readTimeFromInput(&ltc, readLatency)
	readIntFromIntSlice(&kbps, readKbps)
	readIntFromIntSlice(&mtu, readMtu)
	readIntFromIntSlice(&maxConcurrentCalls, readMaxConcurrentCalls)
	readIntFromIntSlice(&reqSizeBytes, readReqSizeBytes)
	readIntFromIntSlice(&respSizeBytes, readRespSizeBytes)
}

type intSliceType []int

func (intSlice *intSliceType) String() string {
	return fmt.Sprintf("%v", *intSlice)
}

func (intSlice *intSliceType) Set(value string) error {
	if len(*intSlice) > 0 {
		return errors.New("interval flag already set")
	}
	for _, num := range strings.Split(value, ",") {
		next, err := strconv.Atoi(num)
		if err != nil {
			return err
		}
		*intSlice = append(*intSlice, next)
	}
	return nil
}

func readIntFromIntSlice(values *[]int, replace intSliceType) {
	// If not set replace in the flag, just return to run the default settings.
	if len(replace) == 0 {
		return
	}
	*values = replace
}

func readTimeFromInput(values *[]time.Duration, replace string) {
	if strings.Compare(replace, "") != 0 {
		*values = []time.Duration{}
		for _, ltc := range strings.Split(replace, ",") {
			duration, err := time.ParseDuration(ltc)
			if err != nil {
				log.Fatal(err.Error())
			}
			*values = append(*values, duration)
		}
	}
}

func main() {
	before()
	featuresPos := make([]int, 8)
	// 0:enableTracing 1:ltc 2:kbps 3:mtu 4:maxC 5:reqSize 6:respSize
	featuresNum := []int{len(enableTrace), len(ltc), len(kbps), len(mtu),
		len(maxConcurrentCalls), len(reqSizeBytes), len(respSizeBytes), len(enableCompressor)}
	initalPos := make([]int, len(featuresPos))
	s := stats.NewStats(10)
	var memStats runtime.MemStats
	var results testing.BenchmarkResult
	var startAllocs, startBytes uint64
	var startTime time.Time
	start := true
	var startTimer = func() {
		runtime.ReadMemStats(&memStats)
		startAllocs = memStats.Mallocs
		startBytes = memStats.TotalAlloc
		startTime = time.Now()
	}
	var stopTimer = func(count int32) {
		runtime.ReadMemStats(&memStats)
		results = testing.BenchmarkResult{N: int(count), T: time.Now().Sub(startTime),
			Bytes: 0, MemAllocs: memStats.Mallocs - startAllocs, MemBytes: memStats.TotalAlloc - startBytes}
	}
	// Run benchmarks
	for !reflect.DeepEqual(featuresPos, initalPos) || start {
		start = false
		tracing := "Trace"
		if !enableTrace[featuresPos[0]] {
			tracing = "noTrace"
		}
		benchFeature := bm.Features{
			EnableTrace:        enableTrace[featuresPos[0]],
			Latency:            ltc[featuresPos[1]],
			Kbps:               kbps[featuresPos[2]],
			Mtu:                mtu[featuresPos[3]],
			MaxConcurrentCalls: maxConcurrentCalls[featuresPos[4]],
			ReqSizeBytes:       reqSizeBytes[featuresPos[5]],
			RespSizeBytes:      respSizeBytes[featuresPos[6]],
			EnableCompressor:   enableCompressor[featuresPos[7]],
		}

		grpc.EnableTracing = enableTrace[featuresPos[0]]
		if runMode[0] {
			fmt.Printf("Unary-%s-%s:\n", tracing, benchFeature.String())
			unaryBenchmark(startTimer, stopTimer, benchFeature, benchtime, s)
			fmt.Println(results.String(), results.MemString())
			fmt.Println(s.String())
			s.Clear()
		}

		if runMode[1] {
			fmt.Printf("Stream-%s-%s\n", tracing, benchFeature.String())
			streamBenchmark(startTimer, stopTimer, benchFeature, benchtime, s)
			fmt.Println(results.String(), results.MemString())
			fmt.Println(s.String())
			s.Clear()
		}
		bm.AddOne(featuresPos, featuresNum)
	}
	after()

}

func before() {
	if memProfileRate > 0 {
		runtime.MemProfileRate = memProfileRate
	}
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
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

func after() {
	if cpuProfile != "" {
		pprof.StopCPUProfile() // flushes profile to disk
	}
	if memProfile != "" {
		f, err := os.Create(memProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "testing: %s\n", err)
			os.Exit(2)
		}
		runtime.GC() // materialize all statistics
		if err = pprof.WriteHeapProfile(f); err != nil {
			fmt.Fprintf(os.Stderr, "testing: can't write %s: %s\n", memProfile, err)
			os.Exit(2)
		}
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

func (nopCompressor) Type() string { return "nop" }

// nopDecompressor is a decompressor that just copies data.
type nopDecompressor struct{}

func (nopDecompressor) Do(r io.Reader) ([]byte, error) { return ioutil.ReadAll(r) }
func (nopDecompressor) Type() string                   { return "nop" }
