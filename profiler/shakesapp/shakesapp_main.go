// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Shakesapp is a web application which starts up a server that can be
// queried to determine how many times a string appears in the works of
// Shakespeare and then sends requests to that server.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/golang/glog"
	"google.golang.org/grpc"

	"shakesapp/shakesapp"
)

var (
	version      = flag.String("version", "1.0.0", "application version")
	projectID    = flag.String("project_id", "", "project ID to run profiler with; only required when running outside of GCP.")
	port         = flag.Int("port", 7788, "service port")
	numReqs      = flag.Int("num_requests", 20, "number of requests to simulate")
	reqsInFlight = flag.Int("requests_in_flight", 1, "number of requests to run in parallel")
	numRounds    = flag.Int("num_rounds", 0, "number of simulation rounds (0 for infinite)")
)

func main() {
	flag.Parse()

	if err := profiler.Start(profiler.Config{
		Service:        "shakesapp",
		ServiceVersion: *version,
		ProjectID:      *projectID,
		MutexProfiling: true,
	}); err != nil {
		log.Fatalf("failed to start profiler: %v", err)
	}

	server := grpc.NewServer()
	shakesapp.RegisterShakespeareServiceServer(server, shakesapp.NewServer())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		glog.Exitf("failed to listen: %v", err)
	}
	go server.Serve(lis)

	ctx := context.Background()
	for i := 1; *numRounds == 0 || i <= *numRounds; i++ {
		start := time.Now()
		glog.Infof("Simulating client requests, round %d", i)
		if err := shakesapp.SimulateClient(ctx, fmt.Sprintf(":%d", *port), *numReqs, *reqsInFlight); err != nil {
			glog.Exitf("Failed to simulate client requests: %v", err)
		}
		delta := time.Since(start).Round(10 * time.Millisecond)
		glog.Infof("Simulated %d requests in %s, rate of %f reqs / sec", *numReqs, delta, float64(*numReqs)/delta.Seconds())
	}
}
