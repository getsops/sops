/*
 * Copyright 2019 gRPC authors.
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
 */

// Package lrs implements load reporting service for xds balancer.
package lrs

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/xds/internal"
	basepb "google.golang.org/grpc/balancer/xds/internal/proto/envoy/api/v2/core/base"
	loadreportpb "google.golang.org/grpc/balancer/xds/internal/proto/envoy/api/v2/endpoint/load_report"
	lrsgrpc "google.golang.org/grpc/balancer/xds/internal/proto/envoy/service/load_stats/v2/lrs"
	lrspb "google.golang.org/grpc/balancer/xds/internal/proto/envoy/service/load_stats/v2/lrs"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal/backoff"
)

const negativeOneUInt64 = ^uint64(0)

// Store defines the interface for a load store. It keeps loads and can report
// them to a server when requested.
type Store interface {
	CallDropped(category string)
	CallStarted(l internal.Locality)
	CallFinished(l internal.Locality, err error)
	ReportTo(ctx context.Context, cc *grpc.ClientConn)
}

type rpcCountData struct {
	// Only atomic accesses are allowed for the fields.
	succeeded  *uint64
	errored    *uint64
	inProgress *uint64
}

func newRPCCountData() *rpcCountData {
	return &rpcCountData{
		succeeded:  new(uint64),
		errored:    new(uint64),
		inProgress: new(uint64),
	}
}

// lrsStore collects loads from xds balancer, and periodically sends load to the
// server.
type lrsStore struct {
	node         *basepb.Node
	backoff      backoff.Strategy
	lastReported time.Time

	drops            sync.Map // map[string]*uint64
	localityRPCCount sync.Map // map[internal.Locality]*rpcCountData
}

// NewStore creates a store for load reports.
func NewStore(serviceName string) Store {
	return &lrsStore{
		node: &basepb.Node{
			Metadata: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					internal.GrpcHostname: {
						Kind: &structpb.Value_StringValue{StringValue: serviceName},
					},
				},
			},
		},
		backoff: backoff.Exponential{
			MaxDelay: 120 * time.Second,
		},
		lastReported: time.Now(),
	}
}

// Update functions are called by picker for each RPC. To avoid contention, all
// updates are done atomically.

// CallDropped adds one drop record with the given category to store.
func (ls *lrsStore) CallDropped(category string) {
	p, ok := ls.drops.Load(category)
	if !ok {
		tp := new(uint64)
		p, _ = ls.drops.LoadOrStore(category, tp)
	}
	atomic.AddUint64(p.(*uint64), 1)
}

func (ls *lrsStore) CallStarted(l internal.Locality) {
	p, ok := ls.localityRPCCount.Load(l)
	if !ok {
		tp := newRPCCountData()
		p, _ = ls.localityRPCCount.LoadOrStore(l, tp)
	}
	atomic.AddUint64(p.(*rpcCountData).inProgress, 1)
}

func (ls *lrsStore) CallFinished(l internal.Locality, err error) {
	p, ok := ls.localityRPCCount.Load(l)
	if !ok {
		// The map is never cleared, only values in the map are reset. So the
		// case where entry for call-finish is not found should never happen.
		return
	}
	atomic.AddUint64(p.(*rpcCountData).inProgress, negativeOneUInt64) // atomic.Add(x, -1)
	if err == nil {
		atomic.AddUint64(p.(*rpcCountData).succeeded, 1)
	} else {
		atomic.AddUint64(p.(*rpcCountData).errored, 1)
	}
}

func (ls *lrsStore) buildStats(clusterName string) []*loadreportpb.ClusterStats {
	var (
		totalDropped  uint64
		droppedReqs   []*loadreportpb.ClusterStats_DroppedRequests
		localityStats []*loadreportpb.UpstreamLocalityStats
	)
	ls.drops.Range(func(category, countP interface{}) bool {
		tempCount := atomic.SwapUint64(countP.(*uint64), 0)
		if tempCount == 0 {
			return true
		}
		totalDropped += tempCount
		droppedReqs = append(droppedReqs, &loadreportpb.ClusterStats_DroppedRequests{
			Category:     category.(string),
			DroppedCount: tempCount,
		})
		return true
	})
	ls.localityRPCCount.Range(func(locality, countP interface{}) bool {
		tempLocality := locality.(internal.Locality)
		tempCount := countP.(*rpcCountData)

		tempSucceeded := atomic.SwapUint64(tempCount.succeeded, 0)
		tempInProgress := atomic.LoadUint64(tempCount.inProgress) // InProgress count is not clear when reading.
		tempErrored := atomic.SwapUint64(tempCount.errored, 0)
		if tempSucceeded == 0 && tempInProgress == 0 && tempErrored == 0 {
			return true
		}

		localityStats = append(localityStats, &loadreportpb.UpstreamLocalityStats{
			Locality: &basepb.Locality{
				Region:  tempLocality.Region,
				Zone:    tempLocality.Zone,
				SubZone: tempLocality.SubZone,
			},
			TotalSuccessfulRequests: tempSucceeded,
			TotalRequestsInProgress: tempInProgress,
			TotalErrorRequests:      tempErrored,
			LoadMetricStats:         nil, // TODO: populate for user loads.
			UpstreamEndpointStats:   nil, // TODO: populate for per endpoint loads.
		})
		return true
	})

	dur := time.Since(ls.lastReported)
	ls.lastReported = time.Now()

	var ret []*loadreportpb.ClusterStats
	ret = append(ret, &loadreportpb.ClusterStats{
		ClusterName:           clusterName,
		UpstreamLocalityStats: localityStats,

		TotalDroppedRequests: totalDropped,
		DroppedRequests:      droppedReqs,
		LoadReportInterval:   ptypes.DurationProto(dur),
	})

	return ret
}

// ReportTo makes a streaming lrs call to cc and blocks.
//
// It retries the call (with backoff) until ctx is canceled.
func (ls *lrsStore) ReportTo(ctx context.Context, cc *grpc.ClientConn) {
	c := lrsgrpc.NewLoadReportingServiceClient(cc)
	var (
		retryCount int
		doBackoff  bool
	)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if doBackoff {
			backoffTimer := time.NewTimer(ls.backoff.Backoff(retryCount))
			select {
			case <-backoffTimer.C:
			case <-ctx.Done():
				backoffTimer.Stop()
				return
			}
			retryCount++
		}

		doBackoff = true
		stream, err := c.StreamLoadStats(ctx)
		if err != nil {
			grpclog.Infof("lrs: failed to create stream: %v", err)
			continue
		}
		if err := stream.Send(&lrspb.LoadStatsRequest{
			Node: ls.node,
		}); err != nil {
			grpclog.Infof("lrs: failed to send first request: %v", err)
			continue
		}
		first, err := stream.Recv()
		if err != nil {
			grpclog.Infof("lrs: failed to receive first response: %v", err)
			continue
		}
		interval, err := ptypes.Duration(first.LoadReportingInterval)
		if err != nil {
			grpclog.Infof("lrs: failed to convert report interval: %v", err)
			continue
		}
		if len(first.Clusters) != 1 {
			grpclog.Infof("lrs: received multiple clusters %v, expect one cluster", first.Clusters)
			continue
		}
		if first.ReportEndpointGranularity {
			// TODO: fixme to support per endpoint loads.
			grpclog.Infof("lrs: endpoint loads requested, but not supported by current implementation")
			continue
		}

		// No backoff afterwards.
		doBackoff = false
		retryCount = 0
		ls.sendLoads(ctx, stream, first.Clusters[0], interval)
	}
}

func (ls *lrsStore) sendLoads(ctx context.Context, stream lrsgrpc.LoadReportingService_StreamLoadStatsClient, clusterName string, interval time.Duration) {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
		case <-ctx.Done():
			return
		}
		if err := stream.Send(&lrspb.LoadStatsRequest{
			Node:         ls.node,
			ClusterStats: ls.buildStats(clusterName),
		}); err != nil {
			grpclog.Infof("lrs: failed to send report: %v", err)
			return
		}
	}
}
