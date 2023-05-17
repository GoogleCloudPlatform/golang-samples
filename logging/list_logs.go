package snippets

// [START logging_list_logs]

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"
)

// listLogs lists all available logs in the project.
func listLogs(w io.Writer, projectID string) error {
	ctx := context.Background()

	client, err := logadmin.NewClient(ctx, projectID)
	if err != nil {
		return err
	}
	defer client.Close()

	iter := client.Logs(ctx)
	for {
		log, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("List logs failed: %w", err)
		}
		fmt.Fprintf(w, "%s\n", log)
	}
	return nil
}

// [END logging_list_logs]
