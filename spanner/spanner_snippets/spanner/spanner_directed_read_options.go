// Copyright 2024 Google LLC
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

// [START spanner_directed_read]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/spanner"
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/api/iterator"
)

//	Shows how to run a query with directed read options.
//	Only one of ExcludeReplicas or IncludeReplicas can be set
//	Each accepts a list of ReplicaSelections which contains Location and Type
//	* `location` - The location must be one of the regions within the
//	multi-region configuration of your database.
//	* `type` - The type of the replica
//	Some examples of using replica_selectors are:
//	* `location:us-east1` --> The "us-east1" replica(s) of any available type
//		will be used to process the request.
//	* `type:READ_ONLY`    --> The "READ_ONLY" type replica(s) in nearest
//	available location will be used to process the
//	request.
//	* `location:us-east1 type:READ_ONLY` --> The "READ_ONLY" type replica(s)
//	in location "us-east1" will be used to process the request.
//		IncludeReplicas also contains an option for AutoFailoverDisabled which when set
//	Spanner will not route requests to a replica outside the
//	IncludeReplicas list when all the specified replicas are unavailable
//	or unhealthy. The default value is `false`

func directedReadOptions(w io.Writer, db string) error {
	// db = `projects/<project>/instances/<instance-id>/database/<database-id>`
	ctx := context.Background()
	directedReadOptionsForClient := &sppb.DirectedReadOptions{
		Replicas: &sppb.DirectedReadOptions_ExcludeReplicas_{
			ExcludeReplicas: &sppb.DirectedReadOptions_ExcludeReplicas{
				ReplicaSelections: []*sppb.DirectedReadOptions_ReplicaSelection{
					{
						Location: "us-east4",
					},
				},
			},
		},
	}
	// DirectedReadOptions can be set at client level and will be used in all read-only transaction requests
	client, err := spanner.NewClientWithConfig(ctx, db, spanner.ClientConfig{DirectedReadOptions: directedReadOptionsForClient})
	if err != nil {
		return err
	}
	defer client.Close()

	// DirectedReadOptions set at Request level will override the options set at Client level.
	directedReadOptionsForRequest := &sppb.DirectedReadOptions{
		Replicas: &sppb.DirectedReadOptions_IncludeReplicas_{
			IncludeReplicas: &sppb.DirectedReadOptions_IncludeReplicas{
				ReplicaSelections: []*sppb.DirectedReadOptions_ReplicaSelection{
					{
						Type: sppb.DirectedReadOptions_ReplicaSelection_READ_ONLY,
					},
				},
				AutoFailoverDisabled: true,
			},
		},
	}

	statement := spanner.Statement{SQL: `SELECT SingerId, AlbumId, AlbumTitle FROM Albums`}
	// // Read rows while passing directedReadOptions directly to the query.
	iter := client.Single().QueryWithOptions(ctx, statement, spanner.QueryOptions{DirectedReadOptions: directedReadOptionsForRequest})
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			return nil
		}
		if err != nil {
			return err
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
		fmt.Fprintf(w, "%d %d %s\n", singerID, albumID, albumTitle)
	}
}

// [END spanner_directed_read]
