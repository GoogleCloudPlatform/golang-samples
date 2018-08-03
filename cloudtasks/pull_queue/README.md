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

## Running the Samples

Set the environment variables:

First, your project ID:

    export PROJECT_ID=my-project-id

Then the queue ID, as specified at queue creation time. Queue IDs already
created can be listed with `gcloud beta tasks queues list`.

    export QUEUE_ID=my-pull-queue

And finally the location ID, which can be discovered with
`gcloud beta tasks queues describe $QUEUE_ID`, with the location embedded in
the "name" value (for instance, if the name is
"projects/my-project/locations/us-central1/queues/my-pull-queue", then the
location is "us-central1").

    export LOCATION_ID=us-central1

## Sample CLI Operations

Create a task for a queue:

    go run main.go create $PROJECT_ID $LOCATION_ID $QUEUE_ID

Pull a task:

    go run main.go pull $PROJECT_ID $LOCATION_ID $QUEUE_ID

Acknowledge task:

    go run main.go acknowledge <task>

* where task is the output from pull task, example:  
`'{"name":"projects/my-project-id/locations/us-central1/queues/my-queue/tasks/1234","scheduleTime":"2017-11-01T22:27:
  53.628279Z"}'`

## Runnning the Tests

Use the built-in environment variable to set the Project ID.

```
export GOLANG_SAMPLES_PROJECT_ID=my-project-id
```

