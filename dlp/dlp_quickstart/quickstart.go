// dlp_quickstart is a basic example of using the Data Loss Prevention API.
package main

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

// [START dlp_quickstart]

func main() {
	ctx := context.Background()

	projectID := "PROJECT_ID"

	// Creates a DLP client.
	client, err := dlp.NewClient(ctx)
	if err != nil {
		log.Fatalf("error creating DLP client: %v", err)
	}

	// The string to inspect.
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
		Parent: "projects/" + projectID,
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
	resp, err := client.InspectContent(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	findings := resp.GetResult().GetFindings()
	if len(findings) == 0 {
		fmt.Println("No findings.")
	}
	fmt.Println("Findings:")
	for _, f := range findings {
		if includeQuote {
			fmt.Println("\tQuote: ", f.GetQuote())
		}
		fmt.Println("\tInfo type: ", f.GetInfoType().GetName())
		fmt.Println("\tLikelihood: ", f.GetLikelihood())
	}
}

// [END dlp_quickstart]
