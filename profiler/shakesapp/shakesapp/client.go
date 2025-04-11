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

package shakesapp

import (
	"context"
	"fmt"
	"math/rand"

	"google.golang.org/grpc"
)

type query struct {
	query          string
	wantMatchCount int64
}

var queries = []query{
	{"hello", 349},
	{"world", 728},
	{"to be, or not to be", 1},
	{"insolence", 14},
}

// SimulateClient creates a client which will send load to the server.
func SimulateClient(ctx context.Context, addr string, numReqs, reqsInFlight int) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewShakespeareServiceClient(conn)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	respErrs := make(chan error)
	inFlightCh := make(chan bool, reqsInFlight)
	for i := 0; i < numReqs; i++ {
		go func() {
			inFlightCh <- true
			defer func() { <-inFlightCh }()
			respErrs <- func() error {
				q := queries[rand.Intn(len(queries))]
				resp, err := client.GetMatchCount(ctx, &ShakespeareRequest{Query: q.query})
				if err != nil {
					return err
				}
				if resp.MatchCount != q.wantMatchCount {
					return fmt.Errorf("GetMatchCount(%q): got %d matches, want %d", q.query, resp.MatchCount, q.wantMatchCount)
				}
				return nil
			}()
		}()
	}
	var retErr error
	for i := 0; i < numReqs; i++ {
		if err := <-respErrs; err != nil {
			retErr = err
		}
	}
	return retErr
}
