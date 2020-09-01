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

package spanner

// [START spanner_insert_datatypes_data]

import (
	"context"
	"io"

	"cloud.google.com/go/spanner"
)

func writeDatatypesData(w io.Writer, db string) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	venueColumns := []string{"VenueId", "VenueName", "VenueInfo", "Capacity", "AvailableDates",
		"LastContactDate", "OutdoorVenue", "PopularityScore", "LastUpdateTime"}
	m := []*spanner.Mutation{
		spanner.InsertOrUpdate("Venues", venueColumns,
			[]interface{}{4, "Venue 4", []byte("Hello World 1"), 1800,
				[]string{"2020-12-01", "2020-12-02", "2020-12-03"},
				"2018-09-02", false, 0.85543, spanner.CommitTimestamp}),
		spanner.InsertOrUpdate("Venues", venueColumns,
			[]interface{}{19, "Venue 19", []byte("Hello World 2"), 6300,
				[]string{"2020-11-01", "2020-11-05", "2020-11-15"},
				"2019-01-15", true, 0.98716, spanner.CommitTimestamp}),
		spanner.InsertOrUpdate("Venues", venueColumns,
			[]interface{}{42, "Venue 42", []byte("Hello World 3"), 3000,
				[]string{"2020-10-01", "2020-10-07"}, "2018-10-01",
				false, 0.72598, spanner.CommitTimestamp}),
	}
	_, err = client.Apply(ctx, m)
	return err
}

// [END spanner_insert_datatypes_data]
