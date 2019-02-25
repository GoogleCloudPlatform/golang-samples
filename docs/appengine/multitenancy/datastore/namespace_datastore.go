// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sample

// [START using_namespaces_with_the_Datastore]
import (
	"context"
	"io"
	"net/http"

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
