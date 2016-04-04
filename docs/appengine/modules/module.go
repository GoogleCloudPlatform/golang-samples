// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"net/http"

	"google.golang.org/appengine/log"
)

// [START communication_between_modules_1]
import "google.golang.org/appengine"

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	module := appengine.ModuleName(ctx)
	instance := appengine.InstanceID()

	log.Infof(ctx, "Received on module %s; instance %s", module, instance)
}

// [END communication_between_modules_1]

func handler2(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	// [START communication_between_modules_2]
	hostname, err := appengine.ModuleHostname(ctx, "my-backend", "", "")
	if err != nil {
		// ...
	}
	url := "http://" + hostname + "/"
	// ...
	// [END communication_between_modules_2]

	_ = url
}
