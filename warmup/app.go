// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

import (
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

func init() {
	http.HandleFunc("/_ah/warmup", warmupHandler)
}

func warmupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	// Do things using ctx.

	log.Infof(ctx, "warmup done")
}
