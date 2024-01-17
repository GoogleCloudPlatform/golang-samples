// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// Testing mongo using mocks is explained here:
// https://medium.com/@victor.neuret/mocking-the-official-mongo-golang-driver-5aad5b226a78
func TestLoad(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("success", func(mt *mtest.T) {
		tm := &trainerManager{
			collection: mt.Coll,
		}
		trainers := []trainer{
			{
				Name: "Ash",
				Age:  10,
				City: "Pallet Town",
			},
			{
				Name: "Misty",
				Age:  10,
				City: "Cerulean City",
			},
			{
				Name: "Brock",
				Age:  15,
				City: "Pewter City",
			},
		}
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := tm.load(context.Background(), trainers)
		assert.Nil(t, err)
	})

	mt.Run("err", func(mt *mtest.T) {
		want := "error inserting records to mongo: must provide at least one element in input slice"
		tm := &trainerManager{
			collection: mt.Coll,
		}
		trainers := []trainer{}
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := tm.load(context.Background(), trainers)
		if want != err.Error() {
			t.Fatalf("expected: %v, got: %v", want, err.Error())
		}
	})
}

func TestList(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	want := []trainer{
		{Name: "Ash", Age: 10, City: "Pallet Town"},
		{Name: "Misty", Age: 10, City: "Cerulean City"},
		{Name: "Brock", Age: 15, City: "Pewter City"},
	}

	mt.Run("success", func(mt *mtest.T) {
		tm := &trainerManager{
			collection: mt.Coll,
		}

		first := mtest.CreateCursorResponse(1, "test.trainers", mtest.FirstBatch, bson.D{
			{Key: "name", Value: "Ash"},
			{Key: "age", Value: 10},
			{Key: "city", Value: "Pallet Town"},
		})
		second := mtest.CreateCursorResponse(1, "test.trainers", mtest.NextBatch, bson.D{
			{Key: "name", Value: "Misty"},
			{Key: "age", Value: 10},
			{Key: "city", Value: "Cerulean City"},
		})
		third := mtest.CreateCursorResponse(1, "test.trainers", mtest.NextBatch, bson.D{
			{Key: "name", Value: "Brock"},
			{Key: "age", Value: 15},
			{Key: "city", Value: "Pewter City"},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.trainers", mtest.NextBatch)
		mt.AddMockResponses(first, second, third, killCursors)

		got, err := tm.list(context.Background())
		if err != nil {
			t.Fatalf("expected: no error, got: %v", err)
		}

		for i, v := range got {
			if !reflect.DeepEqual(want[i], *v) {
				t.Fatalf("expected: %v, got: %v", want[i], *v)
			}
		}
	})
}

func TestListHandler(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("success", func(mt *mtest.T) {
		tm := &trainerManager{
			collection: mt.Coll,
		}

		first := mtest.CreateCursorResponse(1, "test.trainers", mtest.FirstBatch, bson.D{
			{Key: "name", Value: "Ash"},
			{Key: "age", Value: 10},
			{Key: "city", Value: "Pallet Town"},
		})
		second := mtest.CreateCursorResponse(1, "test.trainers", mtest.NextBatch, bson.D{
			{Key: "name", Value: "Misty"},
			{Key: "age", Value: 10},
			{Key: "city", Value: "Cerulean City"},
		})
		third := mtest.CreateCursorResponse(1, "test.trainers", mtest.NextBatch, bson.D{
			{Key: "name", Value: "Brock"},
			{Key: "age", Value: 15},
			{Key: "city", Value: "Pewter City"},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.trainers", mtest.NextBatch)
		mt.AddMockResponses(first, second, third, killCursors)

		// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(listHandler(tm))
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
		// Check the response body is what we expect.
		expected := `[{"Name":"Ash","Age":10,"City":"Pallet Town"},{"Name":"Misty","Age":10,"City":"Cerulean City"},{"Name":"Brock","Age":15,"City":"Pewter City"}]`
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	})
}
