// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sample

// [START asserting_identity_to_Google_APIs]
import (
	"net/http"

	"google.golang.org/appengine/urlfetch"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	urlshortener "google.golang.org/api/urlshortener/v1"
)

// shortenURL returns a short URL which redirects to the provided url,
// using Google's urlshortener API.
func shortenURL(ctx context.Context, url string) (string, error) {
	transport := &oauth2.Transport{
		Source: google.AppEngineTokenSource(ctx, urlshortener.UrlshortenerScope),
		Base:   &urlfetch.Transport{Context: ctx},
	}
	client := &http.Client{Transport: transport}

	svc, err := urlshortener.New(client)
	if err != nil {
		return "", err
	}

	resp, err := svc.Url.Insert(&urlshortener.Url{LongUrl: url}).Do()
	if err != nil {
		return "", err
	}
	return resp.Id, nil
}

// [END asserting_identity_to_Google_APIs]
