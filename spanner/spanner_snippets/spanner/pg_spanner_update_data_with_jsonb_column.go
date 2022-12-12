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

package spanner

// [START spanner_postgresql_jsonb_update_data]

import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
)

// updateDataWithJsonBColumn updates database with JsonB type values
func updateDataWithJsonBColumn(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("updateDataWithJsonBColumn: invalid database id %s", db)
	}

	ctx := context.Background()
	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	type VenueDetails struct {
		Name   spanner.NullString   `json:"name"`
		Rating spanner.NullFloat64  `json:"rating"`
		Open   interface{}          `json:"open"`
		Tags   []spanner.NullString `json:"tags"`
	}

	details_1 := spanner.PGJsonB{Value: []VenueDetails{
		{Name: spanner.NullString{StringVal: "room1", Valid: true}, Open: true},
		{Name: spanner.NullString{StringVal: "room2", Valid: true}, Open: false},
	}, Valid: true}
	details_2 := spanner.PGJsonB{Value: VenueDetails{
		Rating: spanner.NullFloat64{Float64: 9, Valid: true},
		Open:   true,
	}, Valid: true}

	details_3 := spanner.PGJsonB{Value: VenueDetails{
		Name: spanner.NullString{Valid: false},
		Open: map[string]bool{"monday": true, "tuesday": false},
		Tags: []spanner.NullString{{StringVal: "large", Valid: true}, {StringVal: "airy", Valid: true}},
	}, Valid: true}

	cols := []string{"VenueId", "VenueDetails"}
	_, err = client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("Venues", cols, []interface{}{4, details_1}),
		spanner.Update("Venues", cols, []interface{}{19, details_2}),
		spanner.Update("Venues", cols, []interface{}{42, details_3}),
	})

	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Updated data to VenueDetails column\n")

	return nil
}

// [END spanner_postgresql_jsonb_update_data]
