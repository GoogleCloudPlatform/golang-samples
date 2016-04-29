// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Sample analytics demonstrates Google Analytics calls from App Engine flexible environment.
package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/satori/go.uuid"

	"google.golang.org/appengine"
)

var gaPropertyID = mustGetenv("GA_TRACKING_ID")

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func main() {
	http.HandleFunc("/", handle)

	appengine.Main()
}

func handle(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if err := trackEvent(r, "Example", "Test action", "label", nil); err != nil {
		fmt.Fprintf(w, "Event did not track: %v", err)
		return
	}
	fmt.Fprint(w, "Event tracked.")
}

func trackEvent(r *http.Request, category, action, label string, value *uint) error {
	if gaPropertyID == "" {
		return errors.New("analytics: GA_TRACKING_ID environment variable is missing")
	}
	if category == "" || action == "" {
		return errors.New("analytics: category and action are required")
	}

	v := url.Values{
		"v":   {"1"},
		"tid": {gaPropertyID},
		// Anonymously identifies a particular user. See the parameter guide for
		// details:
		// https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#cid
		//
		// Depending on your application, this might want to be associated with the
		// user in a cookie.
		"cid": {uuid.NewV4().String()},
		"t":   {"event"},
		"ec":  {category},
		"ea":  {action},
		"ua":  {r.UserAgent()},
	}

	if label != "" {
		v.Set("el", label)
	}

	if value != nil {
		v.Set("ev", fmt.Sprintf("%d", *value))
	}

	if remoteIP, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		v.Set("uip", remoteIP)
	}

	// NOTE: Google Analytics returns a 200, even if the request is malformed.
	_, err := http.PostForm("https://www.google-analytics.com/collect", v)
	return err
}
