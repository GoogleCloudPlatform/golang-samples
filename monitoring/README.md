# Cloud Monitoring Sample

Simple command-line program to demonstrate connecting to the Google
Monitoring API to retrieve API data.

`listresources` demonstrates how to read the Google Cloud Monitoring v3 environment such as
 Monitored Resources and Metric Descriptors.

`custommetric` demonstrates how to create a custom metric, write a timeseries value to it,
and read it back.

## Prerequisites to run locally:

Go to the [Google Developers Console](https://console.developer.google.com).

    * Go to API Manager -> Credentials
    * Click 'New Credentials', and create a Service Account or [click  here](https://console.developers.google.com/project/_/apiui/credential/serviceaccount)
     Download the JSON for this service account, and set the `GOOGLE_APPLICATION_CREDENTIALS`
     environment variable to point to the file containing the JSON credentials.

    ```
    export GOOGLE_APPLICATION_CREDENTIALS=~/Downloads/<project-id>-0123456789abcdef.json
    ```

# Set Up Your Local Dev Environment

go get the code and change into the directory:

    go get -u github.com/GoogleCloudPlatform/golang-samples/monitoring/...
    cd $GOPATH/src/github.com/GoogleCloudPlatform/golang-samples/monitoring

To run the example that prints the environment, run:

    go run listresources/*.go <your-project-id>

To run the example that creates a custom metric and writes to it, run:

    go run custommetric/*.go <your-project-id>
