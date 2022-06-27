package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

var mt *mtest.T

type mockManager struct{}

func (m mockManager) load(t []trainer) error {
	return nil
}

func (m mockManager) setCollection(c *mongo.Collection) {
}

func (m mockManager) list() ([]*trainer, error) {
	trainers := []*trainer{
		{"Ash", 10, "Pallet Town"},
		{"Misty", 10, "Cerulean City"},
		{"Brock", 15, "Pewter City"},
	}

	return trainers, nil
}

func newMockManager() mockManager {
	m := mockManager{}
	return m
}

func newTestManager() mongoManager {
	m := &trainerManager{}
	return m
}

// Testing mongo using mocks is explained here:
// https://medium.com/@victor.neuret/mocking-the-official-mongo-golang-driver-5aad5b226a78
func TestLoad(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tm := newTestManager()

	mt.Run("success", func(mt *mtest.T) {
		tm.setCollection(mt.Coll)
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

		err := tm.load(trainers)
		assert.Nil(t, err)
	})

	mt.Run("err", func(mt *mtest.T) {
		want := "error inserting records to mongo: must provide at least one element in input slice"
		tm.setCollection(mt.Coll)
		trainers := []trainer{}
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		err := tm.load(trainers)
		if !reflect.DeepEqual(want, err.Error()) {
			t.Fatalf("expected: %v, got: %v", want, err.Error())
		}
	})
}

func TestList(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	tm := newTestManager()

	want := []trainer{
		{"Ash", 10, "Pallet Town"},
		{"Misty", 10, "Cerulean City"},
		{"Brock", 15, "Pewter City"},
	}

	mt.Run("success", func(mt *mtest.T) {
		tm.setCollection(mt.Coll)

		first := mtest.CreateCursorResponse(1, "test.trainers", mtest.FirstBatch, bson.D{
			{"name", "Ash"},
			{"age", 10},
			{"city", "Pallet Town"},
		})
		second := mtest.CreateCursorResponse(1, "test.trainers", mtest.NextBatch, bson.D{
			{"name", "Misty"},
			{"age", 10},
			{"city", "Cerulean City"},
		})
		third := mtest.CreateCursorResponse(1, "test.trainers", mtest.NextBatch, bson.D{
			{"name", "Brock"},
			{"age", 15},
			{"city", "Pewter City"},
		})
		killCursors := mtest.CreateCursorResponse(0, "test.trainers", mtest.NextBatch)
		mt.AddMockResponses(first, second, third, killCursors)

		got, err := tm.list()
		if err != nil {
			t.Fatalf("expected: no error, got: %v", err)
		}

		for i, v := range got {
			fmt.Printf("comp want %+v got %+v\n", want[i], *v)

			if !reflect.DeepEqual(want[i], *v) {
				t.Fatalf("expected: %v, got: %v", want[i], *v)
			}
		}
	})
}

func TestHealthCheckHandler(t *testing.T) {
	tm = newMockManager()

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(listHandler)
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
}
