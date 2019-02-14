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

// [START gae_go_app_identity]
import (
	"context"
	"net/http"

	"google.golang.org/appengine/urlfetch"

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

// [END gae_go_app_identity]
