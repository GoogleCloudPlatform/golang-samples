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

// [START error_reporting_setup_go]
// [START error_reporting_quickstart]

// Sample errorreporting_quickstart contains is a quickstart
// example for the Google Cloud Error Reporting API.
package main

import (
	"context"
	"errors"
	"log"
	"os"

	"cloud.google.com/go/errorreporting"
)

var errorClient *errorreporting.Client

func main() {
	// Set your Google Cloud Platform project ID via environment or explicitly
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	args := os.Args[1:]
	if len(args) > 0 && args[0] != "" {
		projectID = args[0]
	}

	ctx := context.Background()
	var err error
	errorClient, err = errorreporting.NewClient(ctx, projectID, errorreporting.Config{
		ServiceName:    "errorreporting_quickstart",
		ServiceVersion: "0.0.0",
		OnError: func(err error) {
			log.Printf("Could not report the error: %v", err)
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer errorClient.Close()

	err = errors.New("something went wrong")
	if err != nil {
		logAndPrintError(err)
		return
	}
}

func logAndPrintError(err error) {
	/// Client autopopulates the error context of the error. For more details about the context see:
	/// https://cloud.google.com/error-reporting/reference/rest/v1beta1/ErrorContext
	errorClient.Report(errorreporting.Entry{
		Error: err,
	})
	log.Print(err)
}

// [END error_reporting_quickstart]
// [END error_reporting_setup_go]
