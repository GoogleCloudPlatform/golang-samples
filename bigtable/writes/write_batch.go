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

package writes

// [START bigtable_writes_batch]
import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"cloud.google.com/go/bigtable"
)

func writeBatch(w io.Writer, projectID, instanceID string, tableName string) error {
	// projectID := "my-project-id"
	// instanceID := "my-instance-id"
	// tableName := "mobile-time-series"

	ctx := context.Background()
	client, err := bigtable.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return fmt.Errorf("bigtable.NewAdminClient: %v", err)
	}
	defer client.Close()
	tbl := client.Open(tableName)
	columnFamilyName := "stats_summary"
	timestamp := bigtable.Now()

	var muts []*bigtable.Mutation

	binary1 := new(bytes.Buffer)
	binary.Write(binary1, binary.BigEndian, int64(1))

	mut := bigtable.NewMutation()
	mut.Set(columnFamilyName, "connected_wifi", timestamp, binary1.Bytes())
	mut.Set(columnFamilyName, "os_build", timestamp, []byte("12155.0.0-rc1"))
	muts = append(muts, mut)

	mut = bigtable.NewMutation()
	mut.Set(columnFamilyName, "connected_wifi", timestamp, binary1.Bytes())
	mut.Set(columnFamilyName, "os_build", timestamp, []byte("12145.0.0-rc6"))
	muts = append(muts, mut)

	rowKeys := []string{"tablet#a0b81f74#20190501", "tablet#a0b81f74#20190502"}
	if _, err := tbl.ApplyBulk(ctx, rowKeys, muts); err != nil {
		return fmt.Errorf("ApplyBulk: %v", err)
	}

	fmt.Fprintf(w, "Successfully wrote 2 rows: %s\n", rowKeys)
	return nil
}

// [END bigtable_writes_batch]
