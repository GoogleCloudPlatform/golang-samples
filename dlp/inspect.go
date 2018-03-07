package main

import (
	"context"
	"fmt"
	"io"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func inspect(w io.Writer, client *dlp.Client, project, s string) {
	rcr := &dlppb.InspectContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes: []*dlppb.InfoType{
				{
					Name: "US_SOCIAL_SECURITY_NUMBER",
				},
			},
			MinLikelihood: dlppb.Likelihood_LIKELIHOOD_UNSPECIFIED,
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.InspectContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetResult())
}
