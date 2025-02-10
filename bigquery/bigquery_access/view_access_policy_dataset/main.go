package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/bigquery"
)

func main() {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// Creates a client.
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	// Sets the name for the dataset which access list will be asked.
	datasetName := "my_new_dataset"

	// Creates handle for managing dataset
	dataset := client.Dataset(datasetName)

	// Gets dataset's metadata
	metaData, err := dataset.Metadata(ctx)
	if err != nil {
		log.Fatalf("dataset.Metadata: %v", err)
	}

	// Iterate over access permissions
	for _, val := range metaData.Access {
		fmt.Printf("Role %s : %s\n", val.Role, val.Entity)
	}
}
