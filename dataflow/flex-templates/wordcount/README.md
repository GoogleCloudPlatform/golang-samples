# Dataflow flex templates - Wordcount

[![Open in Cloud Shell](http://gstatic.com/cloudssh/images/open-btn.svg)](https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=dataflow/flex-templates/wordcount/README.md)

📝 Docs: [Using Flex Templates](https://cloud.google.com/dataflow/docs/guides/templates/using-flex-templates)

Samples showing how to create and run an
[Apache Beam](https://beam.apache.org/) template with a custom Docker image on
[Google Cloud Dataflow](https://cloud.google.com/dataflow/docs/).

## Before you begin

Follow the
[Getting started with Google Cloud Dataflow](../../README.md)
page, and make sure you have a Google Cloud project with billing enabled
and a *service account JSON key* set up in your `GOOGLE_APPLICATION_CREDENTIALS`
environment variable.
Additionally, for this sample you need the following:

1. [Enable the APIs](https://console.cloud.google.com/flows/enableapi?apiid=appengine.googleapis.com,cloudbuild.googleapis.com):
    App Engine, Cloud Build.

1. Create a
    [Cloud Storage bucket](https://cloud.google.com/storage/docs/creating-buckets).

    ```sh
    export BUCKET="your-gcs-bucket"
    gcloud storage buckets create gs://$BUCKET
    ```

1. Clone the
    [`golang-samples` repository](https://github.com/GoogleCloudPlatform/golang-samples)
    and navigate to the code sample.

    ```sh
    git clone https://github.com/GoogleCloudPlatform/golang-samples.git
    cd golang-samples/dataflow/flex-templates/wordcount
    ```

## Wordcount sample

This sample shows how to deploy an Apache Beam streaming pipeline that reads
text files from [Google Cloud Storage](https://cloud.google.com/storage), 
counts the occurences of each word in the text, and outputs the results back
to Google Cloud Storage.

* [Dockerfile](Dockerfile)
* [wordcount.go](wordcount.go)
* [metadata.json](metadata.json)

### Compiling the pipeline code

Go flex templates take compiled binaries when built, meaning that the containers remain
small; however, this means that the pipeline code must be compiled for the target environment
rather than the machine being used to write the template. For more information, see the 
[Apache Beam SDK documentation on cross-compilation](https://beam.apache.org/documentation/sdks/go-cross-compilation/).

We will compile the Go binary to execute on a linux-amd64 architecture used by Dataflow workers. 

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o wordcount .
```

### Building a container image

We will build the
[Docker](https://docs.docker.com/engine/docker-overview/)
image for the Apache Beam pipeline.
We are using
[Cloud Build](https://cloud.google.com/cloud-build)
so we don't need a local installation of Docker.

> ℹ️  You can speed up subsequent builds with
> [Kaniko cache](https://cloud.google.com/cloud-build/docs/kaniko-cache)
> in Cloud Build.
>
> ```sh
> # (Optional) Enable to use Kaniko cache by default.
> gcloud config set builds/use_kaniko True
> ```

Cloud Build allows you to
[build a Docker image using a `Dockerfile`](https://cloud.google.com/cloud-build/docs/quickstart-docker#build_using_dockerfile).
and saves it into
[Container Registry](https://cloud.google.com/container-registry/),
where the image is accessible to other Google Cloud products.

```sh
export TEMPLATE_IMAGE="gcr.io/$PROJECT/samples/dataflow/wordcount:latest"

# Build the image into Container Registry, this is roughly equivalent to:
#   gcloud auth configure-docker
#   docker image build -t $TEMPLATE_IMAGE .
#   docker push $TEMPLATE_IMAGE
gcloud builds submit --tag "$TEMPLATE_IMAGE" .
```

Images starting with `gcr.io/PROJECT/` are saved into your project's
Container Registry, where the image is accessible to other Google Cloud products.

### Creating a Flex Template

To run a template, you need to create a *template spec* file containing all the
necessary information to run the job, such as the SDK information and metadata.

The [`metadata.json`](metadata.json) file contains additional information for
the template such as the "name", "description", and input "parameters" field.

The template file must be created in a Cloud Storage location,
and is used to run a new Dataflow job.

```sh
export TEMPLATE_PATH="gs://$BUCKET/samples/dataflow/templates/wordcount.json"

# Build the Flex Template.
gcloud dataflow flex-template build $TEMPLATE_PATH \
  --image "$TEMPLATE_IMAGE" \
  --sdk-language "GO" \
  --metadata-file "metadata.json"
```

The template is now available through the template file in the Cloud Storage
location that you specified.

### Running a Dataflow Flex Template pipeline

You can now run the Apache Beam pipeline in Dataflow by referring to the
template file and passing the template
[parameters](https://cloud.google.com/dataflow/docs/guides/specifying-exec-params#setting-other-cloud-dataflow-pipeline-options)
required by the pipeline. For this pipeline the input is optional and will default to a public storage bucket holding
the text of Shakespeare's King Lear.

```sh
export REGION="us-central1"

# Run the Flex Template.
gcloud dataflow flex-template run "wordcount-`date +%Y%m%d-%H%M%S`" \
    --template-file-gcs-location "$TEMPLATE_PATH" \
    --parameters input="projects/$PROJECT/subscriptions/$SUBSCRIPTION" \
    --parameters output="gs://$BUCKET/counts.txt" \
    --region "$REGION"
```

Check the results in your GCS bucket by downloading the output:

```
gcloud alpha storage cp gs://$BUCKET/counts.txt $LOCAL_PATH
```

### Cleaning up

After you've finished this tutorial, you can clean up the resources you created
on Google Cloud so you won't be billed for them in the future.
The following sections describe how to delete or turn off these resources.

#### Clean up the Flex template resources

1. Delete the template spec file from Cloud Storage.

    ```sh
    gcloud storage rm $TEMPLATE_PATH
    ```

1. Delete the Flex Template container image from Container Registry.

    ```sh
    gcloud container images delete $TEMPLATE_IMAGE --force-delete-tags
    ```

#### Clean up Google Cloud project resources

1. Delete the Cloud Storage bucket, this alone does not incur any charges.

    > ⚠️ The following command also deletes all objects in the bucket.
    > These objects cannot be recovered.
    >
    > ```sh
    > gcloud storage rm --recursive gs://$BUCKET
    > ```

## Limitations

There are certain limitations that apply to Flex Templates jobs. 

📝 [Using Flex Templates](https://cloud.google.com/dataflow/docs/guides/templates/using-flex-templates#limitations)
Google Cloud Dataflow documentation page is the authoritative source for the up-to-date information on that.
