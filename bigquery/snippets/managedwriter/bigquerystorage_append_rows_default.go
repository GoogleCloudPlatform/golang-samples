// Copyright 2023 Google LLC
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

// [START bigquerystorage_write_default_complexschema]

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"

	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"github.com/GoogleCloudPlatform/golang-samples/bigquery/snippets/managedwriter/exampleproto"
	"google.golang.org/protobuf/proto"
)

// generateExampleMessages generates a slice of serialized protobuf messages using a statically defined
// and compiled protocol buffer file, and returns the binary serialized representation.
func generateExampleDefaultMessages(numMessages int) ([][]byte, error) {
	msgs := make([][]byte, numMessages)
	for i := 0; i < numMessages; i++ {

		// instantiate a new random source.
		random := rand.New(
			rand.NewSource(time.Now().UnixNano()),
		)

		// Our example data embeds an array of structs, so we'll construct that first.
		sl := make([]*exampleproto.SampleStruct, 5)
		for i := 0; i < int(random.Int63n(5)+1); i++ {
			sl[i] = &exampleproto.SampleStruct{
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

			// This is a required field in the schema, and thus must be present.
			RowNum: proto.Int64(23),

			// StructCol is a single nested message.
			StructCol: &exampleproto.SampleStruct{
				SubIntCol: proto.Int64(random.Int63()),
			},

			// StructList is a repeated array of a nested message.
			StructList: sl,
		}

		b, err := proto.Marshal(m)
		if err != nil {
			return nil, fmt.Errorf("error generating message %d: %w", i, err)
		}
		msgs[i] = b
	}
	return msgs, nil
}

// appendToDefaultStream demonstrates using the managedwriter package to write some example data
// to a default stream.
func appendToDefaultStream(w io.Writer, projectID, datasetID, tableID string) error {
	// projectID := "myproject"
	// datasetID := "mydataset"
	// tableID := "mytable"

	ctx := context.Background()
	// Instantiate a managedwriter client to handle interactions with the service.
	client, err := managedwriter.NewClient(ctx, projectID,
		managedwriter.WithMultiplexing(), // Enables connection sharing.
	)
	if err != nil {
		return fmt.Errorf("managedwriter.NewClient: %w", err)
	}
	// Close the client when we exit the function.
	defer client.Close()

	// We need to communicate the descriptor of the protocol buffer message we're using, which
	// is analagous to the "schema" for the message.  Both SampleData and SampleStruct are
	// two distinct messages in the compiled proto file, so we'll use adapt.NormalizeDescriptor
	// to unify them into a single self-contained descriptor representation.
	var m *exampleproto.SampleData
	descriptorProto, err := adapt.NormalizeDescriptor(m.ProtoReflect().Descriptor())
	if err != nil {
		return fmt.Errorf("NormalizeDescriptor: %w", err)
	}

	// Build the formatted reference to the destination table.
	tableReference := managedwriter.TableParentFromParts(projectID, datasetID, tableID)

	// Instantiate a ManagedStream, which manages low level details like connection state and provides
	// additional features like a future-like callback for appends, etc.  Default streams are provided by
	// the system, so there's no need to create them.
	managedStream, err := client.NewManagedStream(ctx,
		managedwriter.WithType(managedwriter.DefaultStream),
		managedwriter.WithDestinationTable(tableReference),
		managedwriter.WithSchemaDescriptor(descriptorProto),
	)
	if err != nil {
		return fmt.Errorf("NewManagedStream: %w", err)
	}
	// Automatically close the writer when we're done.
	defer managedStream.Close()

	// First, we'll append a single row.
	rows, err := generateExampleDefaultMessages(1)
	if err != nil {
		return fmt.Errorf("generateExampleMessages: %w", err)
	}

	// We can append data asyncronously, so we'll check our appends at the end.
	var results []*managedwriter.AppendResult

	result, err := managedStream.AppendRows(ctx, rows)
	if err != nil {
		return fmt.Errorf("AppendRows first call error: %w", err)
	}
	results = append(results, result)

	// This time, we'll append three more rows in a single request.
	rows, err = generateExampleMessages(3)
	if err != nil {
		return fmt.Errorf("generateExampleMessages: %w", err)
	}
	result, err = managedStream.AppendRows(ctx, rows)
	if err != nil {
		return fmt.Errorf("AppendRows second call error: %w", err)
	}
	results = append(results, result)

	// Finally, we'll append two more rows.
	rows, err = generateExampleMessages(2)
	if err != nil {
		return fmt.Errorf("generateExampleMessages: %w", err)
	}
	result, err = managedStream.AppendRows(ctx, rows)
	if err != nil {
		return fmt.Errorf("AppendRows third call error: %w", err)
	}
	results = append(results, result)

	// We've been collecting references to our status callbacks to allow us to append in a faster
	// asynchronous fashion.  Normally you could do this in another goroutine or similar, but for
	// this example we'll now iterate through those results and verify they were all successful.
	for k, v := range results {
		// GetResult blocks until we receive a response from the API.
		recvOffset, err := v.GetResult(ctx)
		if err != nil {
			return fmt.Errorf("append %d returned error: %w", k, err)
		}
		fmt.Fprintf(w, "Successfully appended data at offset %d.\n", recvOffset)
	}

	// This stream is a default stream, which means it doesn't require any form of finalization
	// or commit.  The rows were automatically committed to the table.
	return nil
}

// [END bigquerystorage_write_default_complexschema]
