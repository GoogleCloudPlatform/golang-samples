// Copyright 2026 Google LLC
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

package firestore

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/firestore"
)

func unixMicrosToTimestampFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unix_micros_timestamp]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.UnixMicrosToTimestamp(firestore.FieldOf("createdAtMicros")).As("createdAtString"),
		)).
		Execute(ctx)
	// [END firestore_unix_micros_timestamp]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unixMillisToTimestampFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unix_millis_timestamp]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.UnixMillisToTimestamp(firestore.FieldOf("createdAtMillis")).As("createdAtString"),
		)).
		Execute(ctx)
	// [END firestore_unix_millis_timestamp]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func unixSecondsToTimestampFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_unix_seconds_timestamp]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.UnixSecondsToTimestamp(firestore.FieldOf("createdAtSeconds")).As("createdAtString"),
		)).
		Execute(ctx)
	// [END firestore_unix_seconds_timestamp]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func timestampAddFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_timestamp_add]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.TimestampAdd(firestore.FieldOf("createdAt"), "day", 3653).As("expiresAt"),
		)).
		Execute(ctx)
	// [END firestore_timestamp_add]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func timestampSubFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_timestamp_sub]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.TimestampSubtract(firestore.FieldOf("expiresAt"), "day", 14).As("sendWarningTimestamp"),
		)).
		Execute(ctx)
	// [END firestore_timestamp_sub]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func timestampToUnixMicrosFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_timestamp_unix_micros]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.TimestampToUnixMicros(firestore.FieldOf("dateString")).As("unixMicros"),
		)).
		Execute(ctx)
	// [END firestore_timestamp_unix_micros]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func timestampToUnixMillisFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_timestamp_unix_millis]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.TimestampToUnixMillis(firestore.FieldOf("dateString")).As("unixMillis"),
		)).
		Execute(ctx)
	// [END firestore_timestamp_unix_millis]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func timestampToUnixSecondsFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_timestamp_unix_seconds]
	snapshot := client.Pipeline().
		Collection("documents").
		Select(firestore.Fields(
			firestore.TimestampToUnixSeconds(firestore.FieldOf("dateString")).As("unixSeconds"),
		)).
		Execute(ctx)
	// [END firestore_timestamp_unix_seconds]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
