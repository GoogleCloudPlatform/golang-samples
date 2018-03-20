package main

import (
	"context"
	"fmt"
	"io"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func redact(w io.Writer, client *dlp.Client, minLikelihood dlppb.Likelihood, infoTypes []string, project, s string) {
	var i []*dlppb.InfoType
	for _, it := range infoTypes {
		i = append(i, &dlppb.InfoType{Name: it})
	}
	rcr := &dlppb.DeidentifyContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:     i,
			MinLikelihood: minLikelihood,
		},
		DeidentifyConfig: &dlppb.DeidentifyConfig{
			Transformation: &dlppb.DeidentifyConfig_InfoTypeTransformations{
				InfoTypeTransformations: &dlppb.InfoTypeTransformations{
					Transformations: []*dlppb.InfoTypeTransformations_InfoTypeTransformation{
						{
							InfoTypes: []*dlppb.InfoType{},
							PrimitiveTransformation: &dlppb.PrimitiveTransformation{
								Transformation: &dlppb.PrimitiveTransformation_RedactConfig{
									RedactConfig: &dlppb.RedactConfig{},
								},
							},
						},
					},
				},
			},
		},
		Item: &dlppb.ContentItem{
			DataItem: &dlppb.ContentItem_Value{
				Value: s,
			},
		},
	}
	r, err := client.DeidentifyContent(context.Background(), rcr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, r.GetItem())
}
