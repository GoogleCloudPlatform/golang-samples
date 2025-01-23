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

// [START gae_creating_namespaces_on_a_per_user_basis]
import (
	"context"

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

// [END gae_creating_namespaces_on_a_per_user_basis]
