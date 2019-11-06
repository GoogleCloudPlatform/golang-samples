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

// Client grpc-ping makes requests to the grpc-ping service.
package main

import (
	"context"
	"crypto/tls"
	"flag"
	"io"
	"log"
	"os"
	"time"

	ptypes "github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/GoogleCloudPlatform/golang-samples/run/grpc-ping/pkg/api/v1"
)

var (
	logger     = log.New(os.Stdout, "", 0)
	serverAddr = flag.String("server", "", "Server address (host:port)")
	serverHost = flag.String("server-host", "", "Host name to which server IP should resolve")
	insecure   = flag.Bool("insecure", false, "Skip SSL validation? [false]")
	skipVerify = flag.Bool("skip-verify", false, "Skip server hostname verification in SSL validation [false]")
	streamSize = flag.Int("stream", 0, "Number of messages to stream [0]")
	message    = flag.String("message", "Hi there", "The body of the content sent to server")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	if *serverHost != "" {
		opts = append(opts, grpc.WithAuthority(*serverHost))
	}
	if *insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		cred := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: *skipVerify,
		})
		opts = append(opts, grpc.WithTransportCredentials(cred))
	}

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		logger.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewPingServiceClient(conn)

	if *streamSize == 0 {
		send(client)
	} else {
		sendStream(client)
	}
}

func send(client pb.PingServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	resp, err := client.Send(ctx, &pb.Request{
		Message: *message,
	})
	if err != nil {
		logger.Fatalf("Error while executing Send: %v", err)
	}

	respMessage := resp.Pong.GetMessage()
	timestamp := ptypes.TimestampString(resp.Pong.GetReceivedOn())
	logger.Println("Unary Request/Unary Response")
	logger.Printf("  Sent Ping: %s", *message)
	logger.Printf("  Received:\n    Pong: %s\n    Server Time: %s", respMessage, timestamp)
}

func sendStream(client pb.PingServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	stream, err := client.SendStream(ctx)
	if err != nil {
		logger.Fatalf("Error client (%v) PingStream: %v", client, err)
	}

	waitCh := make(chan struct{})
	go func() {
		logger.Println("Unary Request/Stream Response")
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitCh)
				return
			}
			if err != nil {
				logger.Fatalf("Failed to receive a response: %v", err)
			}
			c := in.GetPong()
			logger.Printf("  Stream[%d] - Server Time: %s",
				c.GetIndex(), ptypes.TimestampString(c.GetReceivedOn()))
		}
	}()

	i := 0
	for i < *streamSize {
		if err := stream.Send(&pb.Request{
			Message: *message,
		}); err != nil {
			logger.Fatalf("Failed to Send: %v", err)
		}
		i++
	}
	stream.CloseSend()
	<-waitCh
}
