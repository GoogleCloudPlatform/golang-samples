package utils

import (
	"context"
    "log"
	"strconv"
	
	asset "cloud.google.com/go/asset/apiv1p2beta1"
    assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1p2beta1"
    cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

// CleanUp will clear up test data after test.
func CleanUp(ctx context.Context, client *asset.Client, feedName string) {
    req := &assetpb.DeleteFeedRequest{
    	Name: feedName,
    }
    client.DeleteFeed(ctx, req)
}


// GetProjectNumberByID will get projectNumber from projectID by calling
// cloudresourcemanager api
func GetProjectNumberByID(projectID string) string {
	ctx := context.Background()
	cloudresourcemanagerClient, err := cloudresourcemanager.NewService(ctx)
    if err != nil {
            log.Fatal(err)
    }

    project, err := cloudresourcemanagerClient.Projects.Get(projectID).Do()
    if err != nil {
            log.Fatal(err)
    }
    projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
    return projectNumber
}
