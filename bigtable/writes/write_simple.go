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

// Package writes contains snippets related to writing to Cloud Bigtable.
package writes

// [START bigtable_writes_simple]
import (
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"bytes"

	"cloud.google.com/go/bigtable"
)

func writeSimple(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "mobile-time-series"

	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewClient: %v", err)
	}
	defer client.Close()
	tbl := client.Open(tableName)
	columnFamilyName := "stats_summary"
	timestamp := bigtable.Now()

	mut := bigtable.NewMutation()
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, int64(1))

	mut.Set(columnFamilyName, "connected_cell", timestamp, buf.Bytes())
	mut.Set(columnFamilyName, "connected_wifi", timestamp, buf.Bytes())
	mut.Set(columnFamilyName, "os_build", timestamp, []byte("PQ2A.190405.003"))

	rowKey := "phone#4c410523#20190501"
	if err := tbl.Apply(ctx, rowKey, mut); err != nil {
		return fmt.Errorf("Apply: %v", err)
	}

	fmt.Fprintf(w, "Successfully wrote row: %s\n", rowKey)
	return nil
}

// [END bigtable_writes_simple]
