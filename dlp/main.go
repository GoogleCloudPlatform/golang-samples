// Copyright 2018 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// dlp is an example of using the DLP API.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	dlp "cloud.google.com/go/dlp/apiv2beta1"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2beta1"
)

func inspect(w io.Writer, client *dlp.Client, s string) {
	rcr := &dlppb.InspectContentRequest{
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{
					Name: "US_SOCIAL_SECURITY_NUMBER",
				},
			},
			MinLikelihood: dlppb.Likelihood_LIKELIHOOD_UNSPECIFIED,
		},
		Items: []*dlppb.ContentItem{
			{
				Type: "text/plain",
				DataItem: &dlppb.ContentItem_Data{
					Data: []byte(s),
				},
			},
		},
	}
	r, err := client.InspectContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fs := r.GetResults()[0].GetFindings()
	for _, f := range fs {
		fmt.Fprintf(w, "%s\n", f.GetInfoType().GetName())
	}
}

func redact(w io.Writer, client *dlp.Client, s string) {
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
					Data: []byte(s),
				},
			},
		},
	}
	r, err := client.RedactContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "%s\n", r.GetItems()[0].GetData())
}

func main() {
	ctx := context.Background()
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, `Usage: %s CMD "string"\n`, os.Args[0])
		os.Exit(1)
	}
	switch flag.Arg(0) {
	case "inspect":
		inspect(os.Stdout, client, flag.Arg(1))
	case "redact":
		redact(os.Stdout, client, flag.Arg(1))
	}
}
