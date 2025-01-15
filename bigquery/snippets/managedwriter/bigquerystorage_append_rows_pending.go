// Copyright 2022 Google LLC
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

package managedwriter

// [START bigquerystorage_write_pending_complexschema]

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"cloud.google.com/go/bigquery/storage/apiv1/storagepb"
	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/managedwriter/exampleproto"
	"google.golang.org/protobuf/proto"
)

// generateExampleMessages generates a slice of serialized protobuf messages using a statically defined
// and compiled protocol buffer file, and returns the binary serialized representation.
func generateExampleMessages(numMessages int) ([][]byte, error) {
	msgs := make([][]byte, numMessages)
	for i := 0; i < numMessages; i++ {

		random := rand.New(rand.NewSource(time.Now().UnixNano()))

		// Our example data embeds an array of structs, so we'll construct that first.
		sList := make([]*exampleproto.SampleStruct, 5)
		for i := 0; i < int(random.Int63n(5)+1); i++ {
			sList[i] = &exampleproto.SampleStruct{
				SubIntCol: proto.Int64(random.Int63()),
			}
		}

		m := &exampleproto.SampleData{
			BoolCol:    proto.Bool(true),
			BytesCol:   []byte("some bytes"),
			Float64Col: proto.Float64(3.14),
			Int64Col:   proto.Int64(123),
			StringCol:  proto.String("example string value"),

			// These types require special encoding/formatting to transmit.

			// DATE values are number of days since the Unix epoch.

			DateCol: proto.Int32(int32(time.Now().UnixNano() / 86400000000000)),

			// DATETIME uses the literal format.
			DatetimeCol: proto.String("2022-01-01 12:13:14.000000"),

			// GEOGRAPHY uses Well-Known-Text (WKT) format.
			GeographyCol: proto.String("POINT(-122.350220 47.649154)"),

			// NUMERIC and BIGNUMERIC can be passed as string, or more efficiently
			// using a packed byte representation.
			NumericCol:    proto.String("99999999999999999999999999999.999999999"),
			BignumericCol: proto.String("578960446186580977117854925043439539266.34992332820282019728792003956564819967"),

			// TIME also uses literal format.
			TimeCol: proto.String("12:13:14.000000"),

			// TIMESTAMP uses microseconds since Unix epoch.
			TimestampCol: proto.Int64(time.Now().UnixNano() / 1000),

			// Int64List is an array of INT64 types.
			Int64List: []int64{2, 4, 6, 8},

			// This is a required field, and thus must be present.
			RowNum: proto.Int64(23),

			// StructCol is a single nested message.
			StructCol: &exampleproto.SampleStruct{
				SubIntCol: proto.Int64(random.Int63()),
			},

			// StructList is a repeated array of a nested message.
			StructList: sList,
		}

		b, err := proto.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("error generating message %d: %w", i, err)
		}
		msgs[i] = b
	}
	return msgs, nil
}

// appendToPendingStream demonstrates using the managedwriter package to write some example data
// to a pending stream, and then committing it to a table.
func appendToPendingStream(w io.Writer, projectID, datasetID, tableID string) error {
	// projectID := "myproject"
	// datasetID := "mydataset"
	// tableID := "mytable"

	ctx := context.Background()
	// Instantiate a managedwriter client to handle interactions with the service.
	client, err := managedwriter.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("managedwriter.NewClient: %w", err)
	}
	// Close the client when we exit the function.
	defer client.Close()

	// Create a new pending stream.  We'll use the stream name to construct a writer.
	pendingStream, err := client.CreateWriteStream(ctx, &storagepb.CreateWriteStreamRequest{
		Parent: fmt.Sprintf("projects/%s/datasets/%s/tables/%s", projectID, datasetID, tableID),
		WriteStream: &storagepb.WriteStream{
			Type: storagepb.WriteStream_PENDING,
		},
	})
	if err != nil {
		return fmt.Errorf("CreateWriteStream: %w", err)
	}

	// We need to communicate the descriptor of the protocol buffer message we're using, which
	// is analagous to the "schema" for the message.  Both SampleData and SampleStruct are
	// two distinct messages in the compiled proto file, so we'll use adapt.NormalizeDescriptor
	// to unify them into a single self-contained descriptor representation.
	m := &exampleproto.SampleData{}
	descriptorProto, err := adapt.NormalizeDescriptor(m.ProtoReflect().Descriptor())
	if err != nil {
		return fmt.Errorf("NormalizeDescriptor: %w", err)
	}

	// Instantiate a ManagedStream, which manages low level details like connection state and provides
	// additional features like a future-like callback for appends, etc.  NewManagedStream can also create
	// the stream on your behalf, but in this example we're being explicit about stream creation.
	managedStream, err := client.NewManagedStream(ctx, managedwriter.WithStreamName(pendingStream.GetName()),
		managedwriter.WithSchemaDescriptor(descriptorProto))
	if err != nil {
		return fmt.Errorf("NewManagedStream: %w", err)
	}
	defer managedStream.Close()

	// First, we'll append a single row.
	rows, err := generateExampleMessages(1)
	if err != nil {
		return fmt.Errorf("generateExampleMessages: %w", err)
	}

	// We'll keep track of the current offset in the stream with curOffset.
	var curOffset int64
	// We can append data asyncronously, so we'll check our appends at the end.
	var results []*managedwriter.AppendResult

	result, err := managedStream.AppendRows(ctx, rows, managedwriter.WithOffset(0))
	if err != nil {
		return fmt.Errorf("AppendRows first call error: %w", err)
	}
	results = append(results, result)

	// Advance our current offset.
	curOffset = curOffset + 1

	// This time, we'll append three more rows in a single request.
	rows, err = generateExampleMessages(3)
	if err != nil {
		return fmt.Errorf("generateExampleMessages: %w", err)
	}
	result, err = managedStream.AppendRows(ctx, rows, managedwriter.WithOffset(curOffset))
	if err != nil {
		return fmt.Errorf("AppendRows second call error: %w", err)
	}
	results = append(results, result)

	// Advance our offset again.
	curOffset = curOffset + 3

	// Finally, we'll append two more rows.
	rows, err = generateExampleMessages(2)
	if err != nil {
		return fmt.Errorf("generateExampleMessages: %w", err)
	}
	result, err = managedStream.AppendRows(ctx, rows, managedwriter.WithOffset(curOffset))
	if err != nil {
		return fmt.Errorf("AppendRows third call error: %w", err)
	}
	results = append(results, result)

	// Now, we'll check that our batch of three appends all completed successfully.
	// Monitoring the results could also be done out of band via a goroutine.
	for k, v := range results {
		// GetResult blocks until we receive a response from the API.
		recvOffset, err := v.GetResult(ctx)
		if err != nil {
			return fmt.Errorf("append %d returned error: %w", k, err)
		}
		fmt.Fprintf(w, "Successfully appended data at offset %d.\n", recvOffset)
	}

	// We're now done appending to this stream.  We now mark pending stream finalized, which blocks
	// further appends.
	rowCount, err := managedStream.Finalize(ctx)
	if err != nil {
		return fmt.Errorf("error during Finalize: %w", err)
	}

	fmt.Fprintf(w, "Stream %s finalized with %d rows.\n", managedStream.StreamName(), rowCount)

	// To commit the data to the table, we need to run a batch commit.  You can commit several streams
	// atomically as a group, but in this instance we'll only commit the single stream.
	req := &storagepb.BatchCommitWriteStreamsRequest{
		Parent:       managedwriter.TableParentFromStreamName(managedStream.StreamName()),
		WriteStreams: []string{managedStream.StreamName()},
	}

	resp, err := client.BatchCommitWriteStreams(ctx, req)
	if err != nil {
		return fmt.Errorf("client.BatchCommit: %w", err)
	}
	if len(resp.GetStreamErrors()) > 0 {
		return fmt.Errorf("stream errors present: %v", resp.GetStreamErrors())
	}

	fmt.Fprintf(w, "Table data committed at %s\n", resp.GetCommitTime().AsTime().Format(time.RFC3339Nano))

	return nil
}

// [END bigquerystorage_write_pending_complexschema]
