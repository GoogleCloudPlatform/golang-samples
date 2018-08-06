# Google Cloud Tasks Pull Queue Samples

[![Open in Cloud Shell][shell_img]][shell_link]

[shell_img]: http://gstatic.com/cloudssh/images/open-btn.png
[shell_link]: https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=cloudtasks/README.md

Sample command-line program for interacting with the Google Cloud Tasks API
using pull queues.

Pull queues let you add tasks to a queue, then programatically remove and
interact with them. Tasks can be added or processed in any environment,
such as on Google App Engine or Google Compute Engine.

`tasks.go` is a simple command-line program to demonstrate listing queues,
 creating tasks, and pulling and acknowledging tasks.

## Before you begin

1. Follow the installation steps in the [client library's README](library).
2. To set up authentication, please refer to our [authentication getting started guide](authentication).

[library]: https://github.com/GoogleCloudPlatform/google-cloud-go
[authentication]: https://cloud.google.com/docs/authentication/getting-started

## Creating a queue

To create a queue using the Cloud SDK, use the following gcloud command:

    gcloud beta tasks queues create-pull-queue my-pull-queue

## Runnning the Tests

1. Enable the Cloud Tasks API for your project
2. Pre-create a Cloud Tasks Pull Queue
3. Set the project ID environment variable:

```
export GOLANG_SAMPLES_PROJECT_ID=my-project-id
```

With those steps in place, `go test -v ./...`.
