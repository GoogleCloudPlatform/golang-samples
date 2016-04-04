// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

// [START creating_namespaces_on_a_per_user_basis]
import (
	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func namespace(ctx context.Context) context.Context {
	// assumes the user is logged in.
	ctx, err := appengine.Namespace(ctx, user.Current(ctx).ID)
	if err != nil {
		// ...
	}
	return ctx
}

// [END creating_namespaces_on_a_per_user_basis]
