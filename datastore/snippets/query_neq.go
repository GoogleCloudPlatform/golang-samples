package datastore_snippets

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/datastore"
)

func queryNotEquals(projectId string) error {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}
	defer client.Close()

	// [START datastore_not_equals_query]
	q := datastore.NewQuery("TaskList")
	q.FilterField("Task", "!=", []string{"sampleTask"})
	// [END datastore_not_equals_query]

	it := client.Run(ctx, q)
	for {
		var dst interface{}
		key, err := it.Next(&dst)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
		fmt.Printf("Key retrieved: %v\n", key)
	}

	return nil
}
