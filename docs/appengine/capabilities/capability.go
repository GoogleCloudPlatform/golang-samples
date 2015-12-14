package sample

import (
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/capability"
)

// [START datastore_lookup]
func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if !capability.Enabled(ctx, "datastore_v3", "*") {
		http.Error(w, "This service is currently unavailable.", 503)
		return
	}
	// do Datastore lookup ...
}

// [END datastore_lookup]

func example() {
	var ctx context.Context
	// [START intro]
	if !capability.Enabled(ctx, "datastore_v3", "write") {
		// Datastore is in read-only mode.
	}
	// [END intro]
}
