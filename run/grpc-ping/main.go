// Copyright 2019 Google LLC
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

// Sample grpc-ping acts as an intermediary to the ping service.
package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping/pkg/api/v1"
)

// [START cloudrun_grpc_server]
func main() {
	log.Printf("grpc-ping: starting server...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPingServiceServer(grpcServer, &pingService{})
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

// [END cloudrun_grpc_server]

// conn holds an open connection to the ping service.
var conn *grpc.ClientConn

func init() {
	if os.Getenv("GRPC_PING_HOST") != "" {
		var err error
		conn, err = NewConn(os.Getenv("GRPC_PING_HOST"), os.Getenv("GRPC_PING_INSECURE") != "")
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Starting without support for SendUpstream: configure with 'GRPC_PING_HOST' environment variable. E.g., example.com:443")
	}
}
