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
	"os"
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

func main() {
	// Parse command line arguments
	flag.Parse()

	// Check for required flags (golang flags does not support?)
	if *bucketName == "" || *clusterName == "" || *projectID == "" || *pysparkFile == "" || *region == "" || *zoneID == "" {
		log.Fatal("Incorrect arguments specified see 'go-dataproc -help' for help")
	}

	// Create a new Dataproc service
	service, err := getDataprocService()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new storage service
	storageService, err := getStorageService()
	if err != nil {
		log.Fatal(err)
	}

	// Create a cluster
	_, err = createCluster(service, *clusterName, *region, *projectID, *zoneID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cluster created")

	// Wait on the cluster to become active
	fmt.Println("Waiting for cluster to be ready")
	_, err = waitForCluster(service, *clusterName)
	if err != nil {
		log.Fatal(err)
	}

	// Submit a job to the cluster
	fmt.Println("Submitting job")
	jobId, err := submitJob(service, storageService, *pysparkFile, *projectID, *bucketName, *clusterName)
	if err != nil {
		log.Fatal(err)
	}

	// Wait for the job to complete
	fmt.Println("Waiting for job to complete")
	_, err = waitForJob(service, jobId, *projectID, *region)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Job is finished")

	// Get the job output
	clusterData, err := getClusterDataByName(service, *clusterName)
	if err != nil {
		log.Fatal(err)
	}
	output, err := getJobOutput(storageService, *projectID, clusterData[1], clusterData[3], jobId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Job output:\n%s\n", string(output))

	// Delete the cluster
	_, err = deleteCluster(service, *projectID, *region, *clusterName)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Cluster deleted")
	fmt.Printf("Done...")
}

// getDataprocService creates and returns a Cloud Dataproc service object.
func getDataprocService() (service *dataproc.Service, err error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, scope)

	// Create a new Dataproc service
	service, err = dataproc.New(client)

	return service, err
}

// getStorageService creates and returns a Cloud Storage service object.
func getStorageService() (service *storage.Client, err error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	return client, err
}

// createCluster creates a Cloud Dataproc cluster with the given name and region.
func createCluster(service *dataproc.Service, name string, region string,
	project string, zone string) (response string, err error) {
	// Create a gceConfig object for the cluster
	gceConfig := dataproc.GceClusterConfig{
		ZoneUri: fmt.Sprintf(computeUriBase, project, zone),
	}

	// Create a cluserConfig for the cluster
	clusterConfig := dataproc.ClusterConfig{
		GceClusterConfig: &gceConfig,
	}

	// Create a cluster object
	cluster := dataproc.Cluster{
		ClusterName: name,
		ProjectId:   project,
		Config:      &clusterConfig,
	}

	// Create the cluster
	res, err := service.Projects.Regions.Clusters.Create(project, region, &cluster).Do()
	return fmt.Sprintf("%s", res), err
}

// deleteCluster deletes the Cloud Dataproc cluster with the given project, region, and name.
func deleteCluster(service *dataproc.Service, project string, region string, name string) (response string, err error) {
	// Delete the cluster
	res, err := service.Projects.Regions.Clusters.Delete(project, region, name).Do()

	return fmt.Sprintf("%s", res), err
}

// getClusterDataByName returns metadata about the cluster with the given name.
func getClusterDataByName(service *dataproc.Service, name string) (clusterData []string, err error) {
	// Get a list of clusters
	clusters, err := listClusters(service)
	if err != nil {
		return nil, err
	}

	// Find the cluster requested in the list
	for _, item := range clusters {
		if item[0] == name {
			return item, err
		}
	}

	return nil, errors.New("Cluster not found")
}

// getJobOutput returns the text from the job (raw driver output) with the given project, cluser name, bucket id, and job id.
func getJobOutput(storageClient *storage.Client, project string, cluster string, bucket string, job string) (output []byte, err error) {

	// Format the object name
	object := fmt.Sprintf("google-cloud-dataproc-metainfo/%s/jobs/%s/driveroutput.000000000", cluster, job)

	// Read the file
	rc, err := storageClient.Bucket(bucket).Object(object).NewReader(context.Background())
	output, err = ioutil.ReadAll(rc)
	rc.Close()

	return output, err
}

// getJobStatus returns the status of the job with the given job id, project, and region.
func getJobStatus(service *dataproc.Service, jobId string, project string, region string) (status string, err error) {
	//Get the Job's status
	res, err := service.Projects.Regions.Jobs.Get(project, region, jobId).Do()

	return res.Status.State, err
}

// listClusters lists all clusters in the current project.
func listClusters(service *dataproc.Service) (clusters [][]string, err error) {
	// List all clusters in a project for a given region
	res, err := service.Projects.Regions.Clusters.List(*projectID, *region).Do()

	for _, item := range res.Clusters {
		clusterDetails := []string{
			item.ClusterName,
			item.ClusterUuid,
			item.Status.State,
			item.Config.ConfigBucket,
		}
		clusters = append(clusters, clusterDetails)
	}

	return clusters, err
}

// submitJob submits a PySpark job with the given file path to a PySpark file, project, bucket, and cluster.
func submitJob(service *dataproc.Service, storageClient *storage.Client, filepath string, project string, bucket string, cluster string) (jobId string, err error) {
	// Error if the file does not exist
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", err
	}

	// Upload the file to GCS
	job_filename_parts := strings.Split(filepath, "/")
	filename := job_filename_parts[len(job_filename_parts)-1]

	ctx := context.Background()
	wc := storageClient.Bucket(bucket).Object(filename).NewWriter(ctx)
	wc.ContentType = "text/plain"
	wc.ACL = []storage.ACLRule{{storage.AllUsers, storage.RoleReader}}
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	if _, err := wc.Write(file); err != nil {
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Submit the PySpark job
	placement := dataproc.JobPlacement{
		ClusterName: cluster,
	}
	pySparkJob := dataproc.PySparkJob{
		MainPythonFileUri: "gs://" + bucket + "/" + filename,
	}
	jobID := "test-job-" + fmt.Sprintf("%v", time.Now().Unix())
	jobReference := dataproc.JobReference{
		JobId:     jobID,
		ProjectId: project,
	}
	job := dataproc.Job{
		Placement:  &placement,
		PysparkJob: &pySparkJob,
		Reference:  &jobReference,
	}
	jobRequest := dataproc.SubmitJobRequest{
		Job: &job,
	}

	_, err = service.Projects.Regions.Jobs.Submit(*projectID, "global", &jobRequest).Do()
	return jobID, err
}

// waitForCluster waits for a cluster transition from "starting" t0 "running" with the given name.
func waitForCluster(service *dataproc.Service, name string) (running bool, err error) {
	for running == false {
		clusterData, err := getClusterDataByName(service, name)
		if err != nil {
			return false, err
		}
		if clusterData[2] == "RUNNING" {
			running = true
		}

		// Sleep for one second
		time.Sleep(1000 * time.Millisecond)
	}
	return running, err
}

// waitForJob waits for a job to finish with the given job id, project, and region.
func waitForJob(service *dataproc.Service, jobId string, project string, region string) (finished bool, err error) {
	for finished == false {
		jobStatus, err := getJobStatus(service, jobId, project, region)
		if err != nil {
			return false, err
		}
		if jobStatus == "DONE" {
			finished = true
		} else if jobStatus == "ERROR" {
			finished = true
			err = errors.New("Job finished with an error")
		}

		// Sleep for one second
		time.Sleep(1000 * time.Millisecond)
	}

	return finished, err
}
