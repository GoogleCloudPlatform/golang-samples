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

// [START bigtable_functions_quickstart]

// Package bigtable contains an example of using Bigtable from a Cloud Function.
package bigtable

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"cloud.google.com/go/bigtable"
)

// client is a global Bigtable client, to avoid initializing a new client for
// every request.
var client *bigtable.Client
var clientOnce sync.Once

// BigtableRead is an example of reading Bigtable from a Cloud Function.
func BigtableRead(w http.ResponseWriter, r *http.Request) {
	clientOnce.Do(func() {
		// Declare a separate err variable to avoid shadowing client.
		var err error
		client, err = bigtable.NewClient(context.Background(), r.Header.Get("projectID"), r.Header.Get("instanceId"))
		if err != nil {
			http.Error(w, "Error initializing client", http.StatusInternalServerError)
			log.Printf("bigtable.NewClient: %v", err)
			return
		}
	})

	tbl := client.Open(r.Header.Get("tableID"))
	err := tbl.ReadRows(r.Context(), bigtable.PrefixRange("phone#"),
		func(row bigtable.Row) bool {
			osBuild := ""
			for _, col := range row["stats_summary"] {
				if col.Column == "stats_summary:os_build" {
					osBuild = string(col.Value)
				}
			}

			fmt.Fprintf(w, "Rowkey: %s, os_build:  %s\n", row.Key(), osBuild)
			return true
		})

	if err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		log.Printf("tbl.ReadRows(): %v", err)
	}
}

// [END bigtable_functions_quickstart]
