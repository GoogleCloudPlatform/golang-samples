// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

// [START using_namespaces_with_the_Datastore]
import (
	"io"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

type Counter struct {
	Count int64
}

func incrementCounter(ctx context.Context, name string) error {
	key := datastore.NewKey(ctx, "Counter", name, 0, nil)
	return datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		var ctr Counter
		err := datastore.Get(ctx, key, &ctr)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		ctr.Count++
		_, err = datastore.Put(ctx, key, &ctr)
		return err
	}, nil)
}

func someHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	err := incrementCounter(ctx, "SomeRequest")
	if err != nil {
		// ... handle err
	}

	// temporarily use a new namespace
	{
		ctx, err := appengine.Namespace(ctx, "-global-")
		if err != nil {
			// ... handle err
		}
		err = incrementCounter(ctx, "SomeRequest")
		if err != nil {
			// ... handle err
		}
	}

	io.WriteString(w, "Updated counters.\n")
}

func init() {
	http.HandleFunc("/some_url", someHandler)
}

// [END using_namespaces_with_the_Datastore]
