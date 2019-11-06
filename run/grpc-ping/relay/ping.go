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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping/pkg/api/v1"
	ptypes "github.com/golang/protobuf/ptypes"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	out, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll: %q", err)
	}
	if len(out) == 0 {
		out = []byte("Hello World!")
	}

	p := &pb.Request{
		Message: string(out) + " (HTTP relay)",
	}
	url := "https://" + os.Getenv("GRPC_PING_HOST")
	resp, err := PingRequest(Conn, p, url, os.Getenv("GRPC_PING_UNAUTHENTICATED") == "")
	if err != nil {
		log.Printf("pingRequest: %q", err)
	}

	log.Print("received upstream pong")
	out, _ = json.Marshal(resp.Pong)

	w.Write(out)
}

type pingService struct{}

func (s *pingService) Send(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	p := &pb.Request{
		Message: req.GetMessage() + " (gRPC relay)",
	}
	url := "https://" + os.Getenv("GRPC_PING_HOST")
	resp, err := PingRequest(Conn, p, url, os.Getenv("GRPC_PING_UNAUTHENTICATED") == "")
	if err != nil {
		log.Printf("PingRequest: %q", err)
		return nil, fmt.Errorf("could not reach ping service")
	}

	log.Print("received upstream pong")
	return &pb.Response{
		Pong: resp.Pong,
	}, nil
}

// TODO: How to forward the gRPC stream to the ping service?
func (s *pingService) SendStream(stream pb.PingService_SendStreamServer) error {
	var i int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Client disconnected")
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream.Recv: %v", err)
		}

		m := req.GetMessage()
		i++
		log.Printf("Replying to send[%d]: %+v", i, m)

		err = stream.Send(&pb.Response{
			Pong: &pb.Pong{
				Index:      i,
				Message:    m,
				ReceivedOn: ptypes.TimestampNow(),
			},
		})

		if err != nil {
			return fmt.Errorf("stream.Send: failed to send pong: %v", err)
		}
	}
}
