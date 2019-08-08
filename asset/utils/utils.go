package utils

import (
	"context"
    "log"
	"strconv"
	
	asset "cloud.google.com/go/asset/apiv1p2beta1"
    assetpb "google.golang.org/genproto/googleapis/cloud/asset/v1p2beta1"
    cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

func CleanUp(ctx context.Context, client *asset.Client, feedName string) {
    req := &assetpb.DeleteFeedRequest{
    	Name: feedName,
    }
    client.DeleteFeed(ctx, req)
}

func GetProjectNumberById(projectId string) string {
	ctx := context.Background()
	cloudresourcemanagerClient, err := cloudresourcemanager.NewService(ctx)
    if err != nil {
            log.Fatal(err)
    }

    project, err := cloudresourcemanagerClient.Projects.Get(projectId).Do()
    if err != nil {
            log.Fatal(err)
    }
    projectNumber := strconv.FormatInt(project.ProjectNumber, 10)
    return projectNumber
}
