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

func stagesExpressionsExample(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_stages_expressions_example]
	nowMillis := firestore.ConstantOf(1712404800000) // Example timestamp
	trailing30Days := nowMillis.UnixMillisToTimestamp().TimestampSubtract("day", 30)

	snapshot := client.Pipeline().
		Collection("productViews").
		Where(firestore.FieldOf("viewedAt").GreaterThan(trailing30Days)).
		Aggregate(firestore.Accumulators(firestore.CountDistinct("productId").As("uniqueProductViews"))).
		Execute(ctx)
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	// [END firestore_stages_expressions_example]
	fmt.Fprintln(w, results)
	return nil
}
