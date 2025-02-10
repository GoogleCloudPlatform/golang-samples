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

	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")

	// Sets the name for the new dataset.
	datasetName := "my_new_dataset"

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("bigquery.NewClient: %v", err)
	}
	defer client.Close()

	dataset := client.Dataset(datasetName)

	metaData, err := dataset.Metadata(ctx)
	if err != nil {
		log.Fatalf("dataset.Metadata: %v", err)
	}

	for _, val := range metaData.Access {
		fmt.Printf("Role %s : %s\n", val.Role, val.Entity)
	}
}
