package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

func quickstart(project string) {
	// [START dlp_quickstart]
	ctx := context.Background()

	// Instantiates a DLP client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Fatalf("error creating DLP client: %v", err)
	}

	// The string to inspect
	input := "Robert Frost"

	// The minimum likelihood required before returning a match.
	minLikelihood := dlppb.Likelihood_POSSIBLE

	// The maximum number of findings to report (0 = server maximum).
	maxFindings := int32(0)

	// Whether to include the matching string.
	includeQuote := true

	// The infoTypes of information to match.
	infoTypes := []*dlppb.InfoType{
		{
			Name: "PERSON_NAME",
		},
		{
			Name: "US_STATE",
		},
	}

	// Construct item to inspect.
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: input,
		},
	}

	// Construct request.
	req := &dlppb.InspectContentRequest{
		Parent: "projects/" + project,
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:     infoTypes,
			MinLikelihood: minLikelihood,
			Limits: &dlppb.InspectConfig_FindingLimits{
				MaxFindingsPerRequest: maxFindings,
			},
			IncludeQuote: includeQuote,
		},
		Item: item,
	}

	// Run request.
	r, err := client.InspectContent(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fs := r.GetResult().GetFindings()
	if len(fs) == 0 {
		fmt.Println("No findings.")
	}
	fmt.Println("Findings:")
	for _, f := range fs {
		if includeQuote {
			fmt.Println("\tQuote: ", f.GetQuote())
		}
		fmt.Println("\tInfo type: ", f.GetInfoType().GetName())
		fmt.Println("\tLikelihood: ", f.GetLikelihood())
	}
	// [END dlp_quickstart]
}
