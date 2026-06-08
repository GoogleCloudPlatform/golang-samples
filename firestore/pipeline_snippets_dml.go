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

func pipelineUpdate(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_pipeline_update]
	snapshot := client.Pipeline().
		CollectionGroup("users").
		Where(firestore.Not(firestore.FieldExists(firestore.FieldOf("preferences.color")))).
		AddFields(firestore.Selectables(
			firestore.ConstantOfNull().As("preferences.color"),
		)).
		RemoveFields(firestore.Fields("color")).
		Update().
		Execute(ctx)
	// [END firestore_pipeline_update]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func pipelineDelete(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_pipeline_delete]
	snapshot := client.Pipeline().
		CollectionGroup("users").
		Where(firestore.FieldOf("address.country").Equal("USA")).
		Where(firestore.FieldOf("__create_time__").Add(firestore.ConstantOf(10)).LessThan(firestore.CurrentTimestamp())).
		Delete().
		Execute(ctx)
	// [END firestore_pipeline_delete]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
