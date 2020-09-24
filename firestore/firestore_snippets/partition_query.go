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

package main

// [START fs_partition_query]
import (
	"context"
	"fmt"
	"io"

	firestore "cloud.google.com/go/firestore/apiv1"
	"google.golang.org/api/iterator"
	firestorepb "google.golang.org/genproto/googleapis/firestore/v1"
)

// partitionQuery partitions a query by returning partition cursors.
func partitionQuery(ctx context.Context, w io.Writer, parent string) error {
	// projectID := "project-id"
	client, err := firestore.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("firestore.NewClient: %v", err)
	}
	defer client.Close()

	structuredQuery := &firestorepb.StructuredQuery{
		Select: &firestorepb.StructuredQuery_Projection{}}

	req := &firestorepb.PartitionQueryRequest{
		Parent:         parent,
		PartitionCount: 3,
		PageSize:       2,
		QueryType: &firestorepb.PartitionQueryRequest_StructuredQuery{
			StructuredQuery: structuredQuery,
		},
	}
	it := client.PartitionQuery(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "Got partition cursor: %v\n", resp)
	}
	return nil
}

// [END fs_partition_query]
