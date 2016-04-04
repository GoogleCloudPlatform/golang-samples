// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

// [START intro_1]
import (
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	ctx := appengine.NewContext(r)
	u := user.Current(ctx)
	if u == nil {
		url, _ := user.LoginURL(ctx, "/")
		fmt.Fprintf(w, `<a href="%s">Sign in or register</a>`, url)
		return
	}
	url, _ := user.LogoutURL(ctx, "/")
	fmt.Fprintf(w, `Welcome, %s! (<a href="%s">sign out</a>)`, u, url)
}

// [END intro_1]
