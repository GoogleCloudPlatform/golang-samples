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

package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"

	"google.golang.org/api/iterator"
)

func createDocReference(client *firestore.Client) {
	// [START firestore_data_reference_document]
	alovelaceRef := client.Collection("users").Doc("alovelace")
	// [END firestore_data_reference_document]

	_ = alovelaceRef
}

func createCollectionReference(client *firestore.Client) {
	// [START firestore_data_reference_collection]
	usersRef := client.Collection("users")
	// [END firestore_data_reference_collection]

	_ = usersRef
}

func createDocReferenceFromString(client *firestore.Client) {
	// [START firestore_data_reference_document_path]
	alovelaceRef := client.Doc("users/alovelace")
	// [END firestore_data_reference_document_path]

	_ = alovelaceRef
}

func createSubcollectionReference(client *firestore.Client) {
	// [START firestore_data_reference_subcollection]
	messageRef := client.Collection("rooms").Doc("roomA").
		Collection("messages").Doc("message1")
	// [END firestore_data_reference_subcollection]

	_ = messageRef
}

func prepareRetrieve(ctx context.Context, client *firestore.Client) error {
	// [START firestore_data_get_dataset]
	cities := []struct {
		id string
		c  City
	}{
		{id: "SF", c: City{Name: "San Francisco", State: "CA", Country: "USA", Capital: false, Population: 860000}},
		{id: "LA", c: City{Name: "Los Angeles", State: "CA", Country: "USA", Capital: false, Population: 3900000}},
		{id: "DC", c: City{Name: "Washington D.C.", Country: "USA", Capital: true, Population: 680000}},
		{id: "TOK", c: City{Name: "Tokyo", Country: "Japan", Capital: true, Population: 9000000}},
		{id: "BJ", c: City{Name: "Beijing", Country: "China", Capital: true, Population: 21500000}},
	}
	for _, c := range cities {
		_, err := client.Collection("cities").Doc(c.id).Set(ctx, c.c)
		if err != nil {
			return err
		}
	}
	// [END firestore_data_get_dataset]
	return nil
}

func docAsMap(ctx context.Context, client *firestore.Client) (map[string]interface{}, error) {
	// [START firestore_data_get_as_map]
	dsnap, err := client.Collection("cities").Doc("SF").Get(ctx)
	if err != nil {
		return nil, err
	}
	m := dsnap.Data()
	fmt.Printf("Document data: %#v\n", m)
	// [END firestore_data_get_as_map]
	return m, nil
}

func docAsEntity(ctx context.Context, client *firestore.Client) (*City, error) {
	// [START firestore_data_get_as_custom_type]
	dsnap, err := client.Collection("cities").Doc("BJ").Get(ctx)
	if err != nil {
		return nil, err
	}
	var c City
	dsnap.DataTo(&c)
	fmt.Printf("Document data: %#v\n", c)
	// [END firestore_data_get_as_custom_type]
	return &c, nil
}

func multipleDocs(ctx context.Context, client *firestore.Client) error {
	// [START firestore_data_query]
	fmt.Println("All capital cities:")
	iter := client.Collection("cities").Where("capital", "==", true).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(doc.Data())
	}
	// [END firestore_data_query]
	return nil
}

func allDocs(ctx context.Context, client *firestore.Client) error {
	// [START firestore_data_get_all_documents]
	fmt.Println("All cities:")
	iter := client.Collection("cities").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Println(doc.Data())
	}
	// [END firestore_data_get_all_documents]
	return nil
}

func getCollections(ctx context.Context, client *firestore.Client) error {
	// [START firestore_data_get_sub_collections]
	iter := client.Collection("cities").Doc("SF").Collections(ctx)
	for {
		collRef, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		fmt.Printf("Found collection with id: %s\n", collRef.ID)
	}
	// [END firestore_data_get_sub_collections]
	return nil
}
