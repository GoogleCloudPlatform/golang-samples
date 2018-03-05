// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// dlp is an example of using the DLP API.
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	dlp "cloud.google.com/go/dlp/apiv2beta1"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2beta1"
)

func redact(w io.Writer, client *dlp.Client) {
	rcr := &dlppb.RedactContentRequest{
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{
					Name: "US_SOCIAL_SECURITY_NUMBER",
				},
			},
			MinLikelihood: dlppb.Likelihood_LIKELIHOOD_UNSPECIFIED,
		},
		ReplaceConfigs: []*dlppb.RedactContentRequest_ReplaceConfig{
			{
				InfoType:    &dlppb.InfoType{Name: "US_SOCIAL_SECURITY_NUMBER"},
				ReplaceWith: "[redacted]",
			},
		},
		Items: []*dlppb.ContentItem{
			{
				Type: "text/plain",
				DataItem: &dlppb.ContentItem_Data{
					Data: []byte("My SSN is 500112233"),
				},
			},
		},
	}
	r, err := client.RedactContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	i := r.GetItems()[0]
	fmt.Fprint(w, string(i.GetData()))
}

func main() {
	client, err := dlp.NewClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	redact(os.Stdout, client)
}
