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

// [START spanner_field_access_on_nested_struct_parameters]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

func queryWithNestedStructField(w io.Writer, db string) error {
	ctx := context.Background()

	client, err := spanner.NewClient(ctx, db)
	if err != nil {
		return err
	}
	defer client.Close()

	type nameType struct {
		FirstName string
		LastName  string
	}
	type songInfoStruct struct {
		SongName    string
		ArtistNames []nameType
	}
	var songInfo = songInfoStruct{
		SongName: "Imagination",
		ArtistNames: []nameType{
			{FirstName: "Elena", LastName: "Campbell"},
			{FirstName: "Hannah", LastName: "Harris"},
		},
	}
	stmt := spanner.Statement{
		SQL: `SELECT SingerId, @songinfo.SongName FROM Singers
			WHERE STRUCT<FirstName STRING, LastName STRING>(FirstName, LastName)
			IN UNNEST(@songinfo.ArtistNames)`,
		Params: map[string]interface{}{"songinfo": songInfo},
	}
	iter := client.Single().Query(ctx, stmt)
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID int64
		var songName string
		if err := row.Columns(&singerID, &songName); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %s\n", singerID, songName)
	}
}

// [END spanner_field_access_on_nested_struct_parameters]
