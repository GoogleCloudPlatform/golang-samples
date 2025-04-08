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

package firestore

// [START firestore_deps]
import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
)

// [END firestore_deps]

func prepareQuery(ctx context.Context, client *firestore.Client) error {
	// [START firestore_query_filter_dataset]
	cities := []struct {
		id string
		c  City
	}{
		{
			id: "SF",
			c: City{Name: "San Francisco", State: "CA", Country: "USA",
				Capital: false, Population: 860000,
				Regions: []string{"west_coast", "norcal"}},
		},
		{
			id: "LA",
			c: City{Name: "Los Angeles", State: "CA", Country: "USA",
				Capital: false, Population: 3900000,
				Regions: []string{"west_coast", "socal"}},
		},
		{
			id: "DC",
			c: City{Name: "Washington D.C.", Country: "USA",
				Capital: true, Population: 680000,
				Regions: []string{"east_coast"}},
		},
		{
			id: "TOK",
			c: City{Name: "Tokyo", Country: "Japan",
				Capital: true, Population: 9000000,
				Regions: []string{"kanto", "honshu"}},
		},
		{
			id: "BJ",
			c: City{Name: "Beijing", Country: "China",
				Capital: true, Population: 21500000,
				Regions: []string{"jingjinji", "hebei"}},
		},
	}
	for _, c := range cities {
		if _, err := client.Collection("cities").Doc(c.id).Set(ctx, c.c); err != nil {
			return err
		}
	}
	// [END firestore_query_filter_dataset]
	return nil
}

func createQuery(client *firestore.Client) {
	// [START firestore_query_filter_eq_boolean]
	query := client.Collection("cities").Where("capital", "==", true)
	// [END firestore_query_filter_eq_boolean]
	_ = query
}

func createQueryTwo(client *firestore.Client) {
	// [START firestore_query_filter_eq_string]
	query := client.Collection("cities").Where("state", "==", "CA")
	// [END firestore_query_filter_eq_string]
	_ = query
}

func createSimpleQueries(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_filter_single_examples]
	countryQuery := cities.Where("state", "==", "CA")
	popQuery := cities.Where("population", "<", 1000000)
	cityQuery := cities.Where("name", ">=", "San Francisco")
	// [END firestore_query_filter_single_examples]

	_ = countryQuery
	_ = popQuery
	_ = cityQuery
}

func createChainedQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_filter_compound_multi_eq]
	denverQuery := cities.Where("name", "==", "Denver").Where("state", "==", "CO")
	caliQuery := cities.Where("state", "==", "CA").Where("population", "<=", 1000000)
	// [END firestore_query_filter_compound_multi_eq]

	_ = denverQuery
	_ = caliQuery
}

func createInvalidChainedQuery(client *firestore.Client) {
	// Note: this is an instance of a currently unsupported chained query
	cities := client.Collection("cities")
	// [START firestore_query_filter_compound_multi_eq]
	query := cities.Where("country", "==", "USA").Where("population", ">", 5000000)
	// [END firestore_query_filter_compound_multi_eq]

	_ = query
}

func createRangeQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_filter_range_valid]
	stateQuery := cities.Where("state", ">=", "CA").Where("state", "<", "IN")
	populationQuery := cities.Where("state", "==", "CA").Where("population", ">", 1000000)
	// [END firestore_query_filter_range_valid]

	_ = stateQuery
	_ = populationQuery
}

func createInvalidRangeQuery(client *firestore.Client) {
	// Note: This is an invalid range query: range operators
	// are limited to a single field.
	cities := client.Collection("cities")
	// [START firestore_query_filter_range_invalid]
	query := cities.Where("state", ">=", "CA").Where("population", ">", 1000000)
	// [END firestore_query_filter_range_invalid]

	_ = query
}

func createOrderByNameLimitQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_order_limit]
	query := cities.OrderBy("name", firestore.Asc).Limit(3)
	// [END firestore_query_order_limit]

	_ = query
}

func createOrderByNameLimitToLastQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_order_limit]
	query := cities.OrderBy("name", firestore.Asc).LimitToLast(3)
	// [END firestore_query_order_limit]

	_ = query
}

func createOrderByNameDescLimitQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_order_desc_limit]
	query := cities.OrderBy("name", firestore.Desc).Limit(3)
	// [END firestore_query_order_desc_limit]

	_ = query
}

func createMultipleOrderByQuery(client *firestore.Client) {
	// [START firestore_query_order_multi]
	query := client.Collection("cities").OrderBy("state", firestore.Asc).OrderBy("population", firestore.Desc)
	// [END firestore_query_order_multi]
	_ = query
}

func createRangeWithOrderByAndLimitQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_order_limit_field_valid]
	query := cities.Where("population", ">", 2500000).OrderBy("population", firestore.Desc).Limit(2)
	// [END firestore_query_order_limit_field_valid]

	_ = query
}

func createRangeWithOrderByQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_order_with_filter]
	query := cities.Where("population", ">", 2500000).OrderBy("population", firestore.Asc)
	// [END firestore_query_order_with_filter]

	_ = query
}

func createInvalidRangeWithOrderByQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START firestore_query_order_field_invalid]
	// Note: This is an invalid query. It violates the constraint that range
	// and order by are required to be on the same field.
	query := cities.Where("population", ">", 2500000).OrderBy("country", firestore.Asc)
	// [END firestore_query_order_field_invalid]

	_ = query
}

func createSimpleStartAtQuery(client *firestore.Client) {
	// [START firestore_query_cursor_start_at_field_value_single]
	query := client.Collection("cities").OrderBy("population", firestore.Asc).StartAt(1000000)
	// [END firestore_query_cursor_start_at_field_value_single]
	_ = query
}

func createSimpleEndtAtQuery(client *firestore.Client) {
	// [START firestore_query_cursor_end_at_field_value_single]
	query := client.Collection("cities").OrderBy("population", firestore.Asc).EndAt(1000000)
	// [END firestore_query_cursor_end_at_field_value_single]
	_ = query
}

func paginateCursor(ctx context.Context, client *firestore.Client) error {
	// [START firestore_query_cursor_pagination]
	cities := client.Collection("cities")

	// Get the first 25 cities, ordered by population.
	firstPage := cities.OrderBy("population", firestore.Asc).Limit(25).Documents(ctx)
	docs, err := firstPage.GetAll()
	if err != nil {
		return err
	}

	// Get the last document.
	lastDoc := docs[len(docs)-1]

	// Construct a new query to get the next 25 cities.
	secondPage := cities.OrderBy("population", firestore.Asc).
		StartAfter(lastDoc.Data()["population"]).
		Limit(25)

	// ...
	// [END firestore_query_cursor_pagination]
	_ = secondPage
	return nil
}

func createMultipleStartAtQuery(client *firestore.Client) {
	// [START firestore_query_cursor_start_at_field_value_multi]
	// Will return all Springfields.
	client.Collection("cities").
		OrderBy("name", firestore.Asc).
		OrderBy("state", firestore.Asc).
		StartAt("Springfield")

	// Will return Springfields where state comes after Wisconsin.
	client.Collection("cities").
		OrderBy("name", firestore.Asc).
		OrderBy("state", firestore.Asc).
		StartAt("Springfield", "Wisconsin")
	// [END firestore_query_cursor_start_at_field_value_multi]
}

func createInQuery(ctx context.Context, client *firestore.Client) error {
	// [START firestore_query_filter_in]
	cities := client.Collection("cities")
	query := cities.Where("country", "in", []string{"USA", "Japan"}).Documents(ctx)
	// [END firestore_query_filter_in]

	_ = query
	return nil
}

func createInQueryWithArray(ctx context.Context, client *firestore.Client) error {
	// [START firestore_query_filter_in_with_array]
	cities := client.Collection("cities")
	query := cities.Where("regions", "in", [][]string{{"west_coast"}, {"east_coast"}}).Documents(ctx)
	// [END firestore_query_filter_in_with_array]

	_ = query
	return nil
}

func createArrayContainsQuery(ctx context.Context, client *firestore.Client) error {
	cities := client.Collection("cities")
	// [START firestore_query_filter_array_contains]
	query := cities.Where("regions", "array-contains", "west_coast").Documents(ctx)
	// [END firestore_query_filter_array_contains]

	_ = query
	return nil
}

func createArrayContainsAnyQuery(ctx context.Context, client *firestore.Client) error {
	// [START firestore_query_filter_array_contains_any]
	cities := client.Collection("cities")
	query := cities.Where("regions", "array-contains-any", []string{"west_coast", "east_coast"}).Documents(ctx)
	// [END firestore_query_filter_array_contains_any]

	_ = query
	return nil
}

func createStartAtDocSnapshotQuery(ctx context.Context, client *firestore.Client) error {
	// [START firestore_query_cursor_start_at_document]
	cities := client.Collection("cities")
	dsnap, err := cities.Doc("SF").Get(ctx)
	if err != nil {
		fmt.Println(err)
	}
	query := cities.OrderBy("population", firestore.Asc).StartAt(dsnap.Data()["population"]).Documents(ctx)
	// [END firestore_query_cursor_start_at_document]

	_ = query
	return nil
}
