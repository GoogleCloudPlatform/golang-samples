// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Command dataproc uses the Cloud Dataproc API to create a cluster,
// submit a sample job, and delete the cluster when finished.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	dataproc "google.golang.org/api/dataproc/v1"
	storage "cloud.google.com/go/storage"
)

const (
	// Filename of example job (Python) file uploaded to Cloud Storage
	exampleJobName = "test-file"

	// Allows for full access to Google Cloud Platform products
	scope = "https://www.googleapis.com/auth/cloud-platform"

	// Compute URI base
	computeUriBase = "https://www.googleapis.com/compute/v1/projects/%s/zones/%s"
)

var (
	bucketName  = flag.String("bucket", "", "GCS Bucket for storage of job file")
	clusterName = flag.String("cluster-name", "", "Name of cluster")
	projectID   = flag.String("project", "", "Your cloud project ID.")
	pysparkFile = flag.String("pyspark-file", "", "PySpark Job file")
	region      = flag.String("region", "global", "Your cloud project region.")
	zoneID      = flag.String("zone", "", "Cloud Platform zone")
)

// clusterConfig defines the confuration of a cluster including its project, name, and region.
type clusterConfig struct{
		bucket,
    project,
		region,
		state,
		name,
		uuid,
		zone string
}

func main() {
	// Parse command line arguments
	flag.Parse()

	// Check for required flags (golang flags does not support?)
	if *bucketName == "" || *clusterName == "" || *projectID == "" || *pysparkFile == "" || *region == "" || *zoneID == "" {
		log.Fatal("Incorrect arguments specified see 'go-dataproc -help' for help")
	}

	// Create a new Cloud Dataproc service
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, scope)
	if err != nil {
		log.Fatal(err)
	}
	service, err := dataproc.New(client)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Google Clout Storage service
	storageService, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new clusterConfig to hold details about this cluster
	cluster := clusterConfig{project: *projectID, region: *region, name: *clusterName, zone: *zoneID}

	// Create a cluster
	if _, err := createCluster(service, &cluster); err != nil {
		log.Fatal(err)
	}
	log.Println("Cluster created")

	// Wait on the cluster to become active
	log.Println("Waiting for cluster to be ready")
	if _, err = waitForCluster(service, &cluster); err != nil {
		log.Fatal(err)
	}

	// Submit a job to the cluster
	log.Println("Submitting job")
	jobID, err := submitJob(service, storageService, *pysparkFile, *bucketName, &cluster)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the job to complete
	log.Println("Waiting for job to complete")
	if _, err = waitForJob(service, jobID, &cluster); err != nil {
		log.Fatal(err)
	}
	log.Println("Job is finished")

	output, err := getJobOutput(storageService, jobID, &cluster)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Job output:\n%s\n", output)

	// Delete the cluster
	if _, err = deleteCluster(service, &cluster); err != nil {
		log.Fatal(err)
	}
	log.Println("Cluster deleted")
}

// createCluster creates a Cloud Dataproc cluster with the given name and region.
func createCluster(service *dataproc.Service, cluster *clusterConfig) (response string, err error) {
	// Create a gceConfig object for the cluster
	gceConfig := dataproc.GceClusterConfig{
		ZoneUri: fmt.Sprintf(computeUriBase, cluster.project, cluster.zone),
	}

	// Create a (Dataproc API) clusterConfig for the cluster
	clusterConfig := dataproc.ClusterConfig{
		GceClusterConfig: &gceConfig,
	}

	// Create a cluster object
	clusterSpec := dataproc.Cluster{
		ClusterName: cluster.name,
		ProjectId:   cluster.project,
		Config:      &clusterConfig,
	}

	// Create the cluster
	res, err := service.Projects.Regions.Clusters.Create(cluster.project, cluster.region, &clusterSpec).Do()

	return fmt.Sprintf("%s", res), err
}

// deleteCluster deletes the Cloud Dataproc cluster with the given project, region, and name.
func deleteCluster(service *dataproc.Service, cluster *clusterConfig) (response string, err error) {
	// Delete the cluster
	res, err := service.Projects.Regions.Clusters.Delete(cluster.project, cluster.region, cluster.name).Do()

	return fmt.Sprintf("%s", res), err
}

// getJobOutput returns the text from the job (raw driver output) with the given project, cluser name, bucket id, and job id.
func getJobOutput(storageClient *storage.Client, job string, cluster *clusterConfig) (output []byte, err error) {

	// Format the object name based on the Cloud Dataproc service's GCS logging
	// see https://cloud.google.com/dataproc/concepts/driver-output for details
	object := fmt.Sprintf("google-cloud-dataproc-metainfo/%s/jobs/%s/driveroutput.000000000", cluster.uuid, job)

	// Read the file
	rc, err := storageClient.Bucket(cluster.bucket).Object(object).NewReader(context.Background())
	output, err = ioutil.ReadAll(rc)
	rc.Close()

	return output, err
}

// getJobStatus returns the status of the job with the given job id, project, and region.
func getJobStatus(service *dataproc.Service, jobID string, project string, region string) (status string, err error) {
	// Get the Job's status
	res, err := service.Projects.Regions.Jobs.Get(project, region, jobID).Do()
	if err != nil {
		return "", err
	}
	return res.Status.State, nil
}

// listClusters lists all clusters in the current project.
func listClusters(service *dataproc.Service, project string, region string) (clusters [][]string, err error) {
	// List all clusters in a project for a given region
	res, err := service.Projects.Regions.Clusters.List(project, region).Do()
	if err != nil {
		return nil, err
	}

	for _, c := range res.Clusters {
		clusterDetails := []string{
			c.ClusterName,
			c.ClusterUuid,
			c.Status.State,
			c.Config.ConfigBucket,
		}
		clusters = append(clusters, clusterDetails)
	}

	return clusters, err
}

// submitJob submits a PySpark job with the given file path to a PySpark file, project, bucket, and cluster.
func submitJob(service *dataproc.Service, storageClient *storage.Client, filepath string, bucket string, cluster *clusterConfig) (jobID string, err error) {
	// Read the file from disk
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	// Upload the file to GCS
	job_filename_parts := strings.Split(filepath, "/")
	filename := job_filename_parts[len(job_filename_parts)-1]

	ctx := context.Background()
	wc := storageClient.Bucket(bucket).Object(filename).NewWriter(ctx)
	wc.ContentType = "text/plain"
	wc.ACL = []storage.ACLRule{{storage.AllUsers, storage.RoleReader}}
	if _, err := wc.Write(file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Submit the PySpark job
	placement := dataproc.JobPlacement{
		ClusterName: cluster.name,
	}
	pySparkJob := dataproc.PySparkJob{
		MainPythonFileUri: "gs://" + bucket + "/" + filename,
	}
	jobID = "test-job-" + fmt.Sprintf("%v", time.Now().Unix())
	jobReference := dataproc.JobReference{
		JobId:     jobID,
		ProjectId: cluster.project,
	}
	job := dataproc.Job{
		Placement:  &placement,
		PysparkJob: &pySparkJob,
		Reference:  &jobReference,
	}
	jobRequest := dataproc.SubmitJobRequest{
		Job: &job,
	}

	_, err = service.Projects.Regions.Jobs.Submit(cluster.project, cluster.region, &jobRequest).Do()
	return jobID, err
}

// updateClusterMetadata geta metadata and updates the given clusterConfig.
func updateClusterMetadata(service *dataproc.Service, cluster *clusterConfig) (err error) {
	// Get a list of clusters
	clusters, err := listClusters(service, cluster.project, cluster.region)
	if err != nil {
		return err
	}

	// Find the cluster requested in the list
	for _, c := range clusters {
		if c[0] == cluster.name {
			cluster.bucket = c[3]
			cluster.state = c[2]
			cluster.uuid = c[1]
			return nil
		}
	}

	return errors.New("cluster not found")
}

// waitForCluster waits for a cluster transition from "starting" to "running" with the given name.
func waitForCluster(service *dataproc.Service, cluster *clusterConfig) (running bool, err error) {
	for {
		if err := updateClusterMetadata(service, cluster); err != nil {
			return false, err
		}
		if cluster.state == "RUNNING" {
			return true, nil
		}

		// Sleep for one second
		time.Sleep(1000 * time.Millisecond)
	}
}

// waitForJob waits for a job to finish with the given job id, project, and region.
func waitForJob(service *dataproc.Service, jobID string, cluster *clusterConfig) (finished bool, err error) {
	for {
		jobStatus, err := getJobStatus(service, jobID, cluster.project, cluster.region)
		if err != nil {
			return false, err
		}
		if jobStatus == "DONE" {
			return true, nil
		} else if jobStatus == "ERROR" {
			return true, errors.New("Job errored")
		}

		time.Sleep(time.Second)
	}

	return finished, err
}
