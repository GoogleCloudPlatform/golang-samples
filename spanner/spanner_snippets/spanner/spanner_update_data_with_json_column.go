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

// [START spanner_update_data_with_json_column]
import (
	"context"
	"fmt"
	"io"
	"regexp"

	"cloud.google.com/go/spanner"
)

// updateDataWithJsonColumn updates database with Json type values
func updateDataWithJsonColumn(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	matches := regexp.MustCompile("^(.*)/databases/(.*)$").FindStringSubmatch(db)
	if matches == nil || len(matches) != 3 {
		return fmt.Errorf("addJsonColumn: invalid database id %s", db)
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

	details_1, _ := spanner.NullJSON([]VenueDetails{
		{Name: spanner.NullString{"room1", true}, Open: true},
		{Name: spanner.NullString{"room2", true}, Open: false},
	}, true)
	details_2, _ := spanner.NullJSON(VenueDetails{
		Rating: spanner.NullFloat64{9, true},
		Open:   true,
	}, true)
	details_3, _ := spanner.NullJSON(VenueDetails{
		Name: spanner.NullString{"", false},
		Open: map[string]bool{"monday": true, "tuesday": false},
		Tags: []spanner.NullString{spanner.NullString{"large", true}, spanner.NullString{"airy", true}},
	}, true)

	cols := []string{"VenueId", "VenueDetails"}
	_, err = client.Apply(ctx, []*spanner.Mutation{
		spanner.Update("VenueDetails", cols, []interface{}{4, details_1}),
		spanner.Update("VenueDetails", cols, []interface{}{19, details_2}),
		spanner.Update("VenueDetails", cols, []interface{}{42, details_3}),
	})

	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Updated data to VenueDetails column\n")

	return nil
}

// [END spanner_update_data_with_json_column]
