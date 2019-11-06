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
	"fmt"
	"io"
	"log"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping/pkg/api/v1"
	ptypes "github.com/golang/protobuf/ptypes"
)

type pingService struct{}

func (s *pingService) Send(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Print("sending ping response")
	return &pb.Response{
		Pong: &pb.Pong{
			Index:      1,
			Message:    req.GetMessage(),
			ReceivedOn: ptypes.TimestampNow(),
		},
	}, nil
}

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
