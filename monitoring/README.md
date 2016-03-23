# Cloud Monitoring Sample

Simple command-line program to demonstrate connecting to the Google
Monitoring API to retrieve API data.

`listresources` demonstrates how to read the Google Cloud Monitoring v3 environment such as
 Monitored Resources and Metric Descriptors.

`custommetric` demonstrates how to create a custom metric, write a timeseries value to it,
and read it back.


## Prerequisites to run locally:

Create local credentials by running the following command and following the oauth2 flow:

    gcloud beta auth application-default login

# Set Up Your Local Dev Environment

go get the code and change into the directory:

    go get -u github.com/GoogleCloudPlatform/golang-samples/monitoring/...
    cd $GOPATH/src/github.com/GoogleCloudPlatform/golang-samples/monitoring

To run the example that prints the environment, run:

    go run listresources/*.go <your-project-id>

To run the example that creates a custom metric and writes to it, run:

    go run custommetric/*.go <your-project-id>


## Running on GCE, GAE, or other environments

On Google App Engine, the credentials should be found automatically.

On Google Compute Engine, the credentials should be found automatically, but require that
you create the instance with the correct scopes. 

    gcloud compute instances create --scopes="https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/compute,https://www.googleapis.com/auth/compute.readonly" test-instance

If you did not create the instance with the right scopes, you can still upload a JSON service 
account and set GOOGLE_APPLICATION_CREDENTIALS as described below.


## Using a Service Account

In non-Google Cloud environments, GCE instances created without the correct scopes, or local
workstations if the `gcloud beta auth application-default login` command fails, use a Service 
Account by doing the following:

* Go to API Manager -> Credentials
* Click 'New Credentials', and create a Service Account or [click  here](https://console.cloud.google
.com/project/_/apiui/credential/serviceaccount)
 Download the JSON for this service account, and set the `GOOGLE_APPLICATION_CREDENTIALS`
 environment variable to point to the file containing the JSON credentials.


    export GOOGLE_APPLICATION_CREDENTIALS=~/Downloads/<project-id>-0123456789abcdef.json
