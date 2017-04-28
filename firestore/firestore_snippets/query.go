// Copyright 2017 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

// [START fs_dependencies]
import (
	"golang.org/x/net/context"

	"cloud.google.com/go/firestore"
)

// [END fs_dependencies]

func prepareQuery(ctx context.Context, client *firestore.Client) error {
	// [START fs_query_create_examples]
	cities := []struct {
		id string
		c  City
	}{
		{id: "SF", c: City{Name: "San Francisco", State: "CA", Country: "USA", Capital: false, Population: 860000}},
		{id: "LA", c: City{Name: "Los Angeles", State: "CA", Country: "USA", Capital: false, Population: 3900000}},
		{id: "DC", c: City{Name: "Washington D.C.", Country: "USA", Capital: false, Population: 680000}},
		{id: "TOK", c: City{Name: "Tokyo", Country: "Japan", Capital: true, Population: 9000000}},
		{id: "BJ", c: City{Name: "Beijing", Country: "China", Capital: true, Population: 21500000}},
	}
	for _, c := range cities {
		if _, err := client.Collection("cities").Doc(c.id).Set(ctx, c.c); err != nil {
			return err
		}
	}
	// [END fs_query_create_examples]
	return nil
}

func createQuery(client *firestore.Client) {
	// [START fs_create_query]
	query := client.Collection("cities").Where("capital", "==", true)
	// [END fs_create_query]
	_ = query
}

func createQueryTwo(client *firestore.Client) {
	// [START fs_create_query_two]
	query := client.Collection("cities").Where("state", "==", "CA")
	// [END fs_create_qiery_two]
	_ = query
}

func createSimpleQueries(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START fs_simple_queries]
	countryQuery := cities.Where("state", "==", "CA")
	popQuery := cities.Where("population", "<", 1000000)
	cityQuery := cities.Where("name", ">=", "San Francisco")
	// [END fs_simple_queries]

	_ = countryQuery
	_ = popQuery
	_ = cityQuery
}

func createChainedQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START fs_chained_query]
	denverQuery := cities.Where("name", "==", "Denver").Where("state", "==", "CO")
	caliQuery := cities.Where("state", "==", "CA").Where("population", "<=", 1000000)
	// [END fs_chained_query]

	_ = denverQuery
	_ = caliQuery
}

func createInvalidChainedQuery(client *firestore.Client) {
	// Note: this is an instance of a currently unsupported chained query
	cities := client.Collection("cities")
	// [START fs_invalid_chained_query]
	query := cities.Where("country", "==", "USA").Where("population", ">", 5000000)
	// [END fs_invalid_chained_query]

	_ = query
}

func createRangeQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START fs_range_query]
	stateQuery := cities.Where("state", ">=", "CA").Where("state", "<", "IN")
	populationQuery := cities.Where("state", "==", "CA").Where("population", ">", 1000000)
	// [END fs_range_query]

	_ = stateQuery
	_ = populationQuery
}

func createInvalidRangeQuery(client *firestore.Client) {
	// Note: This is an invalid range query: range operators
	// are limited to a single field.
	cities := client.Collection("cities")
	// [START fs_invalid_range_query]
	query := cities.Where("state", ">=", "CA").Where("population", ">", 1000000)
	// [END fs_invalid_range_query]

	_ = query
}

func createOrderByNameLimitQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START fs_order_by_name_limit_query]
	query := cities.OrderBy("name", firestore.Asc).Limit(3)
	// [END fs_order_by_name_limit_query]

	_ = query
}

func createOrderByNameDescLimitQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START fs_order_by_name_desc_limit_query]
	query := cities.OrderBy("name", firestore.Desc).Limit(3)
	// [END fs_order_by_name_desc_limit_query]

	_ = query
}

func createMultipleOrderByQuery(client *firestore.Client) {
	// [START fs_order_by_multiple]
	query := client.Collection("cities").OrderBy("state", firestore.Asc).OrderBy("population", firestore.Desc)
	// [END fs_order_by_multiple]
	_ = query
}

func createRangeWithOrderByAndLimitQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START fs_where_order_by_limit_query]
	query := cities.Where("population", ">", 2500000).OrderBy("population", firestore.Desc).Limit(2)
	// [END fs_where_order_by_limit_query]

	_ = query
}

func createRangeWithOrderByQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START fs_range_order_by_query]
	query := cities.Where("population", ">", 2500000).OrderBy("population", firestore.Asc)
	// [END fs_range_order_by_query]

	_ = query
}

func createInvalidRangeWithOrderByQuery(client *firestore.Client) {
	cities := client.Collection("cities")
	// [START fs_invalid_range_order_by_query]
	// Note: This is an invalid query. It violates the constraint that range
	// and order by are required to be on the same field.
	query := cities.Where("population", ">", 2500000).OrderBy("country", firestore.Asc)
	// [END fs_invalid_range_order_by_query]

	_ = query
}

func createSimpleStartAtQuery(client *firestore.Client) {
	// [START fs_simple_start_at]
	query := client.Collection("cities").OrderBy("population", firestore.Asc).StartAt(1000000)
	// [END fs_simple_start_at]
	_ = query
}

func createSimpleEndtAtQuery(client *firestore.Client) {
	// [START fs_simple_end_at]
	query := client.Collection("cities").OrderBy("population", firestore.Asc).EndAt(1000000)
	// [END fs_simple_end_at]
	_ = query
}

func paginateCursor(ctx context.Context, client *firestore.Client) error {
	// [START fs_paginate_cursor]
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
	// [END fs_paginate_cursor]
	_ = secondPage
	return nil
}

func createMultipleStartAtQuery(client *firestore.Client) {
	// [START fs_start_at_multiple]
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
	// [END fs_start_at_multiple]
}
