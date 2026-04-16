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

func addFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_add_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Add(firestore.FieldOf("soldBooks"), firestore.FieldOf("unsoldBooks")).As("totalBooks"),
		)).
		Execute(ctx)
	// [END firestore_add_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func subtractFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_subtract_function]
	storeCredit := 7
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Subtract(firestore.FieldOf("price"), storeCredit).As("totalCost"),
		)).
		Execute(ctx)
	// [END firestore_subtract_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func multiplyFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_multiply_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Multiply(firestore.FieldOf("price"), firestore.FieldOf("soldBooks")).As("revenue"),
		)).
		Execute(ctx)
	// [END firestore_multiply_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func divideFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_divide_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Divide(firestore.FieldOf("ratings"), firestore.FieldOf("soldBooks")).As("reviewRate"),
		)).
		Execute(ctx)
	// [END firestore_divide_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func modFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_mod_function]
	displayCapacity := 1000
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Mod(firestore.FieldOf("unsoldBooks"), displayCapacity).As("warehousedBooks"),
		)).
		Execute(ctx)
	// [END firestore_mod_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func ceilFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_ceil_function]
	booksPerShelf := 100
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Ceil(firestore.Divide(firestore.FieldOf("unsoldBooks"), booksPerShelf)).As("requiredShelves"),
		)).
		Execute(ctx)
	// [END firestore_ceil_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func floorFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_floor_function]
	snapshot := client.Pipeline().
		Collection("books").
		AddFields(firestore.Selectables(
			firestore.Floor(firestore.Divide(firestore.FieldOf("wordCount"), firestore.FieldOf("pages"))).As("wordsPerPage"),
		)).
		Execute(ctx)
	// [END firestore_floor_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func roundFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_round_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Round(firestore.Multiply(firestore.FieldOf("soldBooks"), firestore.FieldOf("price"))).As("partialRevenue"),
		)).
		Aggregate(firestore.Accumulators(
			firestore.Sum("partialRevenue").As("totalRevenue"),
		)).
		Execute(ctx)
	// [END firestore_round_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func powFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_pow_function]
	googleplexLat := 37.4221
	googleplexLng := -122.0853
	snapshot := client.Pipeline().
		Collection("cities").
		AddFields(firestore.Selectables(
			firestore.Pow(firestore.Multiply(firestore.Subtract(firestore.FieldOf("lat"), googleplexLat), 111), 2).As("latitudeDifference"),
			firestore.Pow(firestore.Multiply(firestore.Subtract(firestore.FieldOf("lng"), googleplexLng), 111), 2).As("longitudeDifference"),
		)).
		Select(firestore.Fields(
			firestore.Sqrt(firestore.Add(firestore.FieldOf("latitudeDifference"), firestore.FieldOf("longitudeDifference"))).
				// Inaccurate for large distances or close to poles
				As("approximateDistanceToGoogle"),
		)).
		Execute(ctx)
	// [END firestore_pow_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func sqrtFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_sqrt_function]
	googleplexLat := 37.4221
	googleplexLng := -122.0853
	snapshot := client.Pipeline().
		Collection("cities").
		AddFields(firestore.Selectables(
			firestore.Pow(firestore.Multiply(firestore.Subtract(firestore.FieldOf("lat"), googleplexLat), 111), 2).As("latitudeDifference"),
			firestore.Pow(firestore.Multiply(firestore.Subtract(firestore.FieldOf("lng"), googleplexLng), 111), 2).As("longitudeDifference"),
		)).
		Select(firestore.Fields(
			firestore.Sqrt(firestore.Add(firestore.FieldOf("latitudeDifference"), firestore.FieldOf("longitudeDifference"))).
				// Inaccurate for large distances or close to poles
				As("approximateDistanceToGoogle"),
		)).
		Execute(ctx)
	// [END firestore_sqrt_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func expFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_exp_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Exp(firestore.FieldOf("rating")).As("expRating"),
		)).
		Execute(ctx)
	// [END firestore_exp_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func lnFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_ln_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Ln(firestore.FieldOf("rating")).As("lnRating"),
		)).
		Execute(ctx)
	// [END firestore_ln_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}

func logFunction(w io.Writer, client *firestore.Client) error {
	ctx := context.Background()
	// [START firestore_log_function]
	snapshot := client.Pipeline().
		Collection("books").
		Select(firestore.Fields(
			firestore.Log(firestore.FieldOf("rating"), 2).As("log2Rating"),
		)).
		Execute(ctx)
	// [END firestore_log_function]
	results, err := snapshot.Results().GetAll()
	if err != nil {
		fmt.Fprintf(w, "snapshot.Results().GetAll failed: %v", err)
		return err
	}
	fmt.Fprintln(w, results)
	return nil
}
