package sample

// [START tasks_within_transactions]
import (
	"net/url"

	"golang.org/x/net/context"

	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/taskqueue"
)

func f(ctx context.Context) {
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		t := taskqueue.NewPOSTTask("/worker", url.Values{
		// ...
		})
		// Use the transaction's context when invoking taskqueue.Add.
		_, err := taskqueue.Add(ctx, t, "")
		if err != nil {
			// Handle error
		}
		// ...
		return nil
	}, nil)
	if err != nil {
		// Handle error
	}
	// ...
}

// [END tasks_within_transactions]

func example() {
	var ctx context.Context

	// [START deleting_tasks]
	// Purge entire queue...
	err := taskqueue.Purge(ctx, "queue1")

	// Delete an individual task...
	t := &taskqueue.Task{Name: "foo"}
	err = taskqueue.Delete(ctx, t, "queue1")
	// [END deleting_tasks]
	_ = err
}
