// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Sample grpc server with a streaming response.
package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-server-streaming/pkg/api/v1"
)

const responseInterval = time.Second

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("timeserver: starting on port %s", port)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	svc := new(timeService)
	server := grpc.NewServer()
	pb.RegisterTimeServiceServer(server, svc)
	if err = server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

type timeService struct{}

func (timeService) StreamTime(req *pb.Request, resp pb.TimeService_StreamTimeServer) error {
	durationSeconds := req.GetDurationSecs()
	finish := time.Now().Add(time.Second * time.Duration(durationSeconds))

	for time.Now().Before(finish) {
		if err := resp.Send(&pb.TimeResponse{
			CurrentTime: ptypes.TimestampNow()}); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		select {
		case <-time.After(responseInterval):
		case <-resp.Context().Done():
			log.Printf("response context closed, exiting response")
			return resp.Context().Err()

		}
	}
	return nil
}
