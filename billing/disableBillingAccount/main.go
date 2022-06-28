// https://cloud.google.com/billing/docs/how-to/notify#cap_disable_billing_to_stop_usage
// Note that associating a project with a *closed* billing account will have much the same effect
// as disabling billing on the project: any paid resources used by the project will be shut down.

package main

import (
	"context"
	"fmt"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	project := os.Getenv("GCP_PROJECT_ID")
	//project := "PROJECT_ID"
	p := cloudbilling.ProjectBillingInfo{
		BillingAccountName: "", // disable BillingAccount Project if empty
		BillingEnabled:     false,
		Name:               project,
		ProjectId:          project,
		ServerResponse:     googleapi.ServerResponse{},
		ForceSendFields:    nil,
		NullFields:         nil,
	}

	token := os.Getenv("GCP_ACCESS_TOKEN") //export GCP_ACCESS_TOKEN=$(gcloud auth print-access-token)
	//token := "AIza..."
	cloudBillingService, err := cloudbilling.NewService(ctx, option.WithAPIKey(os.Getenv(token)))
	if err != nil {
		log.Fatal(err)
	}
	ProjectPath := fmt.Sprintf("projects/%s", project)
	b, err := cloudBillingService.Projects.UpdateBillingInfo(ProjectPath, &p).Do()
	if err != nil {
		log.Panic(err)
	}
	log.Println(b)
}
