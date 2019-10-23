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

// [START bigquerystorage_quickstart]

// The bigquery_storage_quickstart application demonstrates usage of the
// BigQuery Storage read API.  It demonstrates API features such as column
// projection (limiting the output to a subset of a table's columns),
// column filtering (using simple predicates to filter records on the server
// side), establishing the snapshot time (reading data from the table at a
// specific point in time), and decoding Avro row blocks using the third party
// "github.com/linkedin/goavro" library.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"sort"
	"sync"
	"time"

	bqStorage "cloud.google.com/go/bigquery/storage/apiv1beta1"
	"github.com/golang/protobuf/ptypes"
	gax "github.com/googleapis/gax-go/v2"
	goavro "github.com/linkedin/goavro/v2"
	bqStoragepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1beta1"
	"google.golang.org/grpc"
)

// rpcOpts is used to configure the underlying gRPC client to accept large
// messages.  The BigQuery Storage API may send message blocks up to 10MB
// in size.
var rpcOpts = gax.WithGRPCOptions(
	grpc.MaxCallRecvMsgSize(1024 * 1024 * 11),
)

// Command-line flags.
var (
	projectID = flag.String("project_id", "",
		"Cloud Project ID, used for session creation.")
	snapshotMillis = flag.Int64("snapshot_millis", 0,
		"Snapshot time to use for reads, represented in epoch milliseconds format.  Default behavior reads current data.")
)

func main() {
	flag.Parse()
	ctx := context.Background()
	bqStorageClient, err := bqStorage.NewBigQueryStorageClient(ctx)
	if err != nil {
		log.Fatalf("NewBigQueryStorageClient: %v", err)
	}
	defer bqStorageClient.Close()

	// Verify we've been provided a parent project which will contain the read session.  The
	// session may exist in a different project than the table being read.
	if *projectID == "" {
		log.Fatalf("No parent project ID specified, please supply using the --project_id flag.")
	}

	// This example uses baby name data from the public datasets.
	readTable := &bqStoragepb.TableReference{
		ProjectId: "bigquery-public-data",
		DatasetId: "usa_names",
		TableId:   "usa_1910_current",
	}

	// We limit the output columns to a subset of those allowed in the table,
	// and set a simple filter to only report names from the state of
	// Washington (WA).
	tableReadOptions := &bqStoragepb.TableReadOptions{
		SelectedFields: []string{"name", "number", "state"},
		RowRestriction: `state = "WA"`,
	}

	readSessionRequest := &bqStoragepb.CreateReadSessionRequest{
		Parent:         fmt.Sprintf("projects/%s", *projectID),
		TableReference: readTable,
		ReadOptions:    tableReadOptions,
		// This API can also deliver data serialized in Apache Arrow format.
		// This example leverages Apache Avro.
		Format: bqStoragepb.DataFormat_AVRO,
		// We use a LIQUID strategy in this example because we only
		// read from a single stream.  Consider BALANCED if you're consuming
		// multiple streams concurrently and want more consistent stream sizes.
		ShardingStrategy: bqStoragepb.ShardingStrategy_LIQUID,
	}

	// Set a snapshot time if it's been specified.
	if *snapshotMillis > 0 {
		ts, err := ptypes.TimestampProto(time.Unix(0, *snapshotMillis*1000))
		if err != nil {
			log.Fatalf("Invalid snapshot millis (%d): %v", *snapshotMillis, err)
		}
		readSessionRequest.TableModifiers = &bqStoragepb.TableModifiers{
			SnapshotTime: ts,
		}
	}

	// Create the session from the request.
	session, err := bqStorageClient.CreateReadSession(ctx, readSessionRequest, rpcOpts)
	if err != nil {
		log.Fatalf("CreateReadSession: %v", err)
	}

	if len(session.GetStreams()) == 0 {
		log.Fatalf("no streams in session.  if this was a small query result, consider writing to output to a named table.")
	}

	// We'll use only a single stream for reading data from the table.  Because
	// of dynamic sharding, this will yield all the rows in the table. However,
	// if you wanted to fan out multiple readers you could do so by having a
	// reader process each individual stream.
	readStream := session.GetStreams()[0]

	ch := make(chan *bqStoragepb.AvroRows)

	// Use a waitgroup to coordinate the reading and decoding goroutines.
	var wg sync.WaitGroup

	// Start the reading in one goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := processStream(ctx, bqStorageClient, readStream, ch); err != nil {
			log.Fatalf("processStream failure: %v", err)
		}
		close(ch)
	}()

	// Start Avro processing and decoding in another goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := processAvro(ctx, session.GetAvroSchema().GetSchema(), ch)
		if err != nil {
			log.Fatalf("Error processing avro: %v", err)
		}
	}()

	// Wait until both the reading and decoding goroutines complete.
	wg.Wait()

}

// printDatum prints the decoded row datum.
func printDatum(d interface{}) {
	m, ok := d.(map[string]interface{})
	if !ok {
		log.Printf("failed type assertion: %v", d)
	}
	// Go's map implementation returns keys in a random ordering, so we sort
	// the keys before accessing.
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("%s: %-20v ", key, valueFromTypeMap(m[key]))
	}
	fmt.Println()
}

// valueFromTypeMap returns the first value/key in the type map.  This function
// is only suitable for simple schemas, as complex typing such as arrays and
// records necessitate a more robust implementation.  See the goavro library
// and the Avro specification for more information.
func valueFromTypeMap(field interface{}) interface{} {
	m, ok := field.(map[string]interface{})
	if !ok {
		return nil
	}
	for _, v := range m {
		// Return the first key encountered.
		return v
	}
	return nil
}

// processStream reads rows from a single storage Stream, and sends the Avro
// data blocks to a channel. This function will retry on transient stream
// failures and bookmark progress to avoid re-reading data that's already been
// successfully transmitted.
func processStream(ctx context.Context, client *bqStorage.BigQueryStorageClient, st *bqStoragepb.Stream, ch chan<- *bqStoragepb.AvroRows) error {
	var offset int64

	// Streams may be long-running.  Rather than using a global retry for the
	// stream, implement a retry that resets once progress is made.
	retryLimit := 3

	for {
		retries := 0
		// Send the initiating request to start streaming row blocks.
		rowStream, err := client.ReadRows(ctx, &bqStoragepb.ReadRowsRequest{
			ReadPosition: &bqStoragepb.StreamPosition{
				Stream: st,
				Offset: offset,
			}}, rpcOpts)
		if err != nil {
			return fmt.Errorf("Couldn't invoke ReadRows: %v", err)
		}

		// Process the streamed responses.
		for {
			r, err := rowStream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				retries++
				if retries >= retryLimit {
					return fmt.Errorf("processStream retries exhausted: %v", err)
				}
			}

			rc := r.GetRowCount()
			if rc > 0 {
				// Bookmark our progress in case of retries and send the rowblock on the channel.
				offset = offset + rc
				// We're making progress, reset retries.
				retries = 0
				ch <- r.GetAvroRows()
			}
		}
	}
}

// processAvro receives row blocks from a channel, and uses the provided Avro
// schema to decode the blocks into individual row messages for printing.  Will
// continue to run until the channel is closed or the provided context is
// cancelled.
func processAvro(ctx context.Context, schema string, ch <-chan *bqStoragepb.AvroRows) error {
	// Establish a decoder that can process blocks of messages using the
	// reference schema. All blocks share the same schema, so the decoder
	// can be long-lived.
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		return fmt.Errorf("couldn't create codec: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			// Context was cancelled.  Stop.
			return nil
		case rows, ok := <-ch:
			if !ok {
				// Channel closed, no further avro messages.  Stop.
				return nil
			}
			undecoded := rows.GetSerializedBinaryRows()
			for len(undecoded) > 0 {
				datum, remainingBytes, err := codec.NativeFromBinary(undecoded)

				if err != nil {
					if err == io.EOF {
						break
					}
					return fmt.Errorf("decoding error with %d bytes remaining: %v", len(undecoded), err)
				}
				printDatum(datum)
				undecoded = remainingBytes
			}
		}
	}
}

// [END bigquerystorage_quickstart]
