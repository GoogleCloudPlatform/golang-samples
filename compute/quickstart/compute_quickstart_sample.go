package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	collection *mongo.Collection
	ctx        = context.Background()
	mongohost  string
)

func main() {
	var err error
	mongohost = os.Getenv("HOST")
	uri := fmt.Sprintf("mongodb://%s:%s", mongohost, "80")
	clientOptions := options.Client().ApplyURI(uri)

	if client, err = mongo.Connect(ctx, clientOptions); err != nil {
		log.Fatal(fmt.Printf("error connecting to mongo: %s", err))
	}

	collection = client.Database("test").Collection("trainers")

	trainers := []interface{}{
		trainer{"Ash", 10, "Pallet Town"},
		trainer{"Misty", 10, "Cerulean City"},
		trainer{"Brock", 15, "Pewter City"},
	}

	if _, err := collection.InsertMany(ctx, trainers); err != nil {
		log.Fatal(fmt.Printf("error inserting records to mongo: %s", err))
	}

	http.HandleFunc("/", listHandler)
	log.Fatal(http.ListenAndServe(":80", nil))
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	trainers, err := listTrainers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	j, err := json.Marshal(trainers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)
	return
}

type trainer struct {
	Name string
	Age  int
	City string
}

func listTrainers() ([]*trainer, error) {
	var results []*trainer

	cur, err := collection.Find(ctx, bson.D{{}}, options.Find())
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem trainer
		if err := cur.Decode(&elem); err != nil {
			return nil, err
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)
	return results, nil
}
