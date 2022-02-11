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

// This sample application leverages the BigQuery storage Write API to bulk
// ingest a CSV file directly.
package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"sync"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"cloud.google.com/go/storage"
	storagepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1"

	"github.com/GoogleCloudPlatform/golang-samples/bigquery/bigquery_managedwriter_csvimporter/converter"
)

// This is the main execution loop.
func main() {
	// Define some flags governing the behavior of the application.
	var (
		projectID      = flag.String("project_id", "", "cloud project ID")
		datasetID      = flag.String("dataset", "", "destination BigQuery dataset name")
		tableID        = flag.String("table", "", "destination BigQuery table name")
		bucketID       = flag.String("bucket", "cloud-samples-data", "Cloud Storage Bucket with the CSV")
		objectID       = flag.String("object", "bigquery/flights/bq-flights-lax.csv", "Cloud Storage Object path")
		streamFanout   = flag.Int("fanout", 20, "stream fanout used to write data")
		optimalReqSize = flag.Int64("optimal_request_size", 8*1e6, "optimal request size in bytes")
		verbose        = flag.Bool("verbose", true, "whether to log verbosely during execution")
	)
	flag.Parse()
	if err := validateFlags(*projectID, *datasetID, *tableID, *streamFanout, *optimalReqSize); err != nil {
		log.Fatalf("flag error: %v", err)
	}

	ctx := context.Background()

	// Instantiate a client for interacting with the storage API.
	mwClient, err := managedwriter.NewClient(ctx, *projectID)
	if err != nil {
		log.Fatalf("failed to create managedwriter client: %v", err)
	}
	defer mwClient.Close()

	converter, err := setupConverter(ctx, *projectID, *datasetID, *tableID)
	if err != nil {
		log.Fatalf("setupConverter: %v", err)
	}

	// tableName is the format for a table in the BigQuery storage API.
	tableName := fmt.Sprintf("projects/%s/datasets/%s/tables/%s", *projectID, *datasetID, *tableID)

	// Create the managed streams we'll use for writing.
	managedStreams, err := setupManagedStreams(ctx, mwClient, *streamFanout, tableName, converter, *verbose)
	if err != nil {
		log.Fatalf("setupStreams: %v", err)
	}

	// dataChan is used to distribute serialized rows from the CSV reader to the storage writers.
	dataChan := make(chan [][]byte, 100)
	// readerErr will collect the error if GCS reading fails.
	var readerErr error
	// writerErr will collect errors from each of the individual writers.
	writerErr := make([]error, *streamFanout)

	var wg sync.WaitGroup

	// Start the GCS reader.

	wg.Add(1)
	go func() {
		readerErr = processGCSObject(ctx, *bucketID, *objectID, *optimalReqSize, converter, dataChan, *verbose)
		wg.Done()
	}()

	// Start the individual stream writers.
	var finalizedRows int64
	for i := 0; i < *streamFanout; i++ {
		wg.Add(1)
		id := i
		go func() {
			ms := managedStreams[id]
			if *verbose {
				log.Printf("writer ID %d assigned stream %s", id, ms.StreamName())
			}
			rows, err := processWrites(ctx, id, ms, dataChan, *verbose)
			if err != nil {
				writerErr[id] = err
			} else {
				finalizedRows = finalizedRows + rows
			}
			wg.Done()
		}()
	}

	// wait for reader and writers to complete before verifying.
	wg.Wait()

	if readerErr != nil {
		log.Fatalf("aborting commit, reading failed: %v", err)
	}

	var writeFailures int
	for i := 0; i < *streamFanout; i++ {
		err = writerErr[i]
		if err != nil {
			writeFailures = writeFailures + 1
			if *verbose {
				log.Printf("writer ID %d failed: %v", i, err)
			}
		}
	}

	if writeFailures > 0 {
		log.Fatalf("aborting commit, %d writers observed failures", writeFailures)
	}

	if *verbose {
		log.Printf("finalized %d rows", finalizedRows)
	}

	if err := commitStreams(ctx, mwClient, managedStreams); err != nil {
		log.Fatalf("commit failed: %v", err)
	}

	if *verbose {
		log.Printf("committed data successfully")
	}

}

// validateFlags does some minimal validation of flag values.
func validateFlags(project, dataset, table string, fanout int, reqSize int64) error {
	if project == "" {
		return errors.New("no cloud project ID specified.")
	}
	if dataset == "" {
		return errors.New("no destination dataset specified.")
	}
	if table == "" {
		return errors.New("no destination table specified.")
	}
	if fanout < 1 || fanout > 100 {
		return fmt.Errorf("stream fanout factor outside allowed range 1..100: %d specified", fanout)
	}
	if reqSize < 1024 || reqSize > 9*1e6 {
		return fmt.Errorf("optimal request size needs to be between 1KB and 9MB: %d bytes specified", reqSize)
	}
	return nil
}

func processGCSObject(ctx context.Context, bucket, object string, reqSize int64, converter *converter.CSVConverter, dataChan chan<- [][]byte, verbose bool) error {
	defer close(dataChan)
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("storage NewReader (%q %q): %v", bucket, object, err)
	}
	defer rc.Close()

	rdr := csv.NewReader(rc)
	// We expect the first line of the CSV to have the column names we can map to BigQuery
	// column names.
	header, err := rdr.Read()
	if err == io.EOF {
		return errors.New("csv appears to be empty")
	}
	if err != nil {
		return fmt.Errorf("failure reading header row: %v", err)
	}

	if err := converter.Validate(header); err != nil {
		return fmt.Errorf("header validation failed: %v", err)
	}
	// dataLine tracks what line of the CSV file we're on.
	var dataLine int64

	// pendingRows retains rows until we've read enough bytes to reach desired size, or we finish reading.
	var pendingRows [][]byte
	var pendingBytes int64

	// CSV processing loop.
	for {
		row, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read data line %d: %v", dataLine, err)
		}
		dataLine = dataLine + 1
		// convert list to map of colnames->values.
		rowMap := make(map[string]string)
		for i, col := range header {
			if i < len(row) {
				rowMap[col] = row[i]
			}
		}
		// convert rowmap to serialized message.
		b, err := converter.Convert(rowMap)
		if err != nil {
			return fmt.Errorf("failed to serialize data line %d for writing: %v", dataLine, err)
		}
		testSize := int64(len(b)) + pendingBytes
		if testSize < reqSize {
			// new row is still within limit.  retain it as part of the pending row.
			pendingRows = append(pendingRows, b)
			pendingBytes = testSize
		} else {
			if verbose {
				log.Printf("reader releasing %d rows (%d KB)", len(pendingRows), pendingBytes/1024)
			}
			dataChan <- pendingRows
			// Reset pendingRows and tracking.
			pendingRows = [][]byte{b}
			pendingBytes = int64(len(b))
		}
	}
	// We've reached the end of the CSV.  If there's any data still buffered, release it for appending.
	if len(pendingRows) > 0 {
		if verbose {
			log.Printf("reader final release of %d rows (%d KB) for appending", len(pendingRows), pendingBytes/1024)
		}
		dataChan <- pendingRows
	}
	if verbose {
		log.Printf("reader summary: scanned %d lines", dataLine)
	}
	return nil
}

// processWrites handles writing data to a given managed stream.  It returns the number of rows finalized in the stream or error.
func processWrites(ctx context.Context, id int, ms *managedwriter.ManagedStream, ch <-chan [][]byte, verbose bool) (int64, error) {
	var offset int64
	var result *managedwriter.AppendResult
	var err error
	for data := range ch {
		if len(data) > 0 {
			result, err = ms.AppendRows(ctx, data, managedwriter.WithOffset(offset))
			if err != nil {
				return 0, fmt.Errorf("failed to send append at offset %d: %v", offset, err)
			}
			if verbose {
				log.Printf("writer id %d sent %d rows at offset %d", id, len(data), offset)
			}
			offset = offset + int64(len(data))
		}
	}
	// Channel closed, so wait until we get acknowledgement of the final write then return.
	_, err = result.GetResult(ctx)
	if err != nil {
		return 0, fmt.Errorf("writer %d received error checking status of final write: %v", id, err)
	}

	// proceed to finalize the stream.
	numRows, err := ms.Finalize(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to finalize: %v", err)
	}
	if verbose {
		log.Printf("writer id %d finalized with %d rows (%s)", id, numRows, ms.StreamName())
	}
	return numRows, nil
}

// setupConverter creates a data converter suitable for writing data into a given table, and returns both
func setupConverter(ctx context.Context, project, dataset, table string) (*converter.CSVConverter, error) {
	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		return nil, fmt.Errorf("couldn't instantiate BigQuery client: %v", err)
	}
	defer client.Close()

	tableRef := client.Dataset(dataset).Table(table)

	meta, err := tableRef.Metadata(ctx)
	if err != nil {
		return nil, fmt.Errorf("couldn't get table information: %v", err)
	}

	tableSchema, err := adapt.BQSchemaToStorageTableSchema(meta.Schema)
	if err != nil {
		return nil, fmt.Errorf("failed to convert schema: %v", err)
	}

	converter, err := converter.NewCSVConverter(tableSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to create converter: %v", err)
	}

	return converter, nil
}

// setupManagedStreams creates the necessary managed stream instances.
func setupManagedStreams(ctx context.Context, client *managedwriter.Client, numStreams int, tableName string, converter *converter.CSVConverter, verbose bool) ([]*managedwriter.ManagedStream, error) {

	// Get the descriptor for the proto we'll be writing from the converter.
	protoSchema, err := converter.ProtoSchema()
	if err != nil {
		return nil, fmt.Errorf("converter failed to generate proto schema descriptor: %v", err)
	}

	streams := make([]*managedwriter.ManagedStream, numStreams)

	// Construct any remaining streams implicitly.
	for i := 0; i < numStreams; i++ {
		ms, err := client.NewManagedStream(ctx,
			managedwriter.WithDestinationTable(tableName),
			managedwriter.WithSchemaDescriptor(protoSchema),
			managedwriter.WithType(managedwriter.PendingStream))
		if err != nil {
			return nil, fmt.Errorf("failed to create managed stream %d: %v", i, err)
		}
		streams[i] = ms
		if verbose {
			log.Printf("created managed stream, stream ID: %s", ms.StreamName())
		}
	}
	return streams, nil
}

func commitStreams(ctx context.Context, client *managedwriter.Client, streams []*managedwriter.ManagedStream) error {

	var names []string
	for _, s := range streams {
		names = append(names, s.StreamName())
		s.Close()
	}

	req := &storagepb.BatchCommitWriteStreamsRequest{
		Parent:       managedwriter.TableParentFromStreamName(names[0]),
		WriteStreams: names,
	}
	resp, err := client.BatchCommitWriteStreams(ctx, req)
	if err != nil {
		return fmt.Errorf("client.BatchCommit: %v", err)
	}
	if len(resp.GetStreamErrors()) > 0 {
		return fmt.Errorf("stream errors present: %v", resp.GetStreamErrors())
	}
	return nil
}
