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

// [START spanner_set_custom_timeout_and_retry]

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/spanner"
	spapi "cloud.google.com/go/spanner/apiv1"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/grpc/codes"
)

func setCustomTimeoutAndRetry(w io.Writer, db string) error {
	// Set a custom retry setting.
	co := &spapi.CallOptions{
		ExecuteSql: []gax.CallOption{
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    500 * time.Millisecond,
					Max:        16000 * time.Millisecond,
					Multiplier: 1.5,
				})
			}),
		},
	}
	client, err := spanner.NewClientWithConfig(context.Background(), db, spanner.ClientConfig{CallOptions: co})
	if err != nil {
		return err
	}
	defer client.Close()

	// Set timeout to 60s.
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err = client.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT Singers (SingerId, FirstName, LastName)
					VALUES (20, 'George', 'Washington')`,
		}

		rowCount, err := txn.Update(ctx, stmt)
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%d record(s) inserted.\n", rowCount)
		return nil
	})
	return err
}

// [END spanner_set_custom_timeout_and_retry]
