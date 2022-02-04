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

package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"strings"

	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/storage"
)

// This is the main execution loop.
func main() {
	// Define some flags governing the behavior of the application.
	var (
		projectID      = flag.String("project_id", "", "Cloud project ID used for running the sample")
		datasetID      = flag.String("dataset", "dataset", "destination BigQuery dataset name")
		tableID        = flag.String("table", "datatable", "destination BigQuery table name")
		bucketID       = flag.String("bucket", "cloud-samples-data", "Cloud Storage Bucket with the CSV")
		objectID       = flag.String("object", "bigquery/flights/bq-flights-lax.csv", "Cloud Storage Object path")
		streamFanout   = flag.Int("fanout", 20, "stream fanout used to write data")
		optimalReqSize = flag.Int64("optimal_request_size", 4*1e6, "optimal request size in bytes")
		verbose        = flag.Bool("verbose", true, "whether to log verbosely during execution")
	)
	flag.Parse()
	if err := validateFlags(*projectID, *streamFanout, *optimalReqSize); err != nil {
		log.Fatalf("flag error: %v", err)
	}

	ctx := context.Background()
	if err := processGCSObject(ctx, *bucketID, *objectID, *optimalReqSize, fakeConverter, *verbose); err != nil {
		log.Fatalf("error processing csv: %v", err)
	}
	setupManagedStreams(ctx, *streamFanout, *projectID, *datasetID, *tableID, nil, *verbose)
}

// validateFlags does some minimal validation of flag values.
func validateFlags(project string, fanout int, reqSize int64) error {
	if project == "" {
		return errors.New("no cloud project ID specified.")
	}
	if fanout < 1 || fanout > 100 {
		return fmt.Errorf("stream fanout factor outside allowed range 1..100: %d specified", fanout)
	}
	if reqSize < 1024 || reqSize > 9*1e6 {
		return fmt.Errorf("optimal request size needs to be between 1KB and 9MB: %d bytes specified", reqSize)
	}
	return nil
}

func processGCSObject(ctx context.Context, bucket, object string, reqSize int64, converterF func([]string) ([]byte, error), verbose bool) error {
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
	if verbose {
		log.Printf("CSV headers:\n%s", strings.Join(header, "\n"))
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
		// Invoke the converter function to serialize the CSV data as a binary protocol buffer message.
		b, err := converterF(row)
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
				log.Printf("releasing %d rows (%d KB) for appending", len(pendingRows), pendingBytes/1024)
			}
			// TODO release to chan
			// dataChan <- pendingRows

			// Reset pendingRows and tracking.
			pendingRows = [][]byte{b}
			pendingBytes = int64(len(b))
		}
	}
	// We've reached the end of the CSV.  If there's any data still buffered, release it for appending.
	if len(pendingRows) > 0 {
		if verbose {
			log.Printf("final release of %d rows (%d KB) for appending", len(pendingRows), pendingBytes/1024)
		}
		// dataChan <- pendingRows
	}
	if verbose {
		log.Printf("CSV scanned %d lines", dataLine)
	}
	return nil
}

// setupManagedStreams sets up a set of managed stream instances, and a converter function.
func setupManagedStreams(ctx context.Context, numStreams int, project, dataset, table string, csvHeaders []string, verbose bool) ([]*managedwriter.ManagedStream, func([]string) ([]byte, error), error) {
	return nil, nil, fmt.Errorf("unimplemented")
}

func fakeConverter(data []string) ([]byte, error) {
	return []byte(strings.Join(data, "\t")), nil
}
