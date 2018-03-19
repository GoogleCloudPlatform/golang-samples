// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_list_info_types]

// infoTypes returns the info types in the given language and matching the given filter.
func infoTypes(w io.Writer, client *dlp.Client, languageCode, filter string) {
	req := &dlppb.ListInfoTypesRequest{
		LanguageCode: languageCode,
		Filter:       filter,
	}
	r, err := client.ListInfoTypes(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	for _, it := range r.GetInfoTypes() {
		fmt.Fprintln(w, it.GetName())
	}
}

// [END dlp_list_info_types]
