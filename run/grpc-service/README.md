# grpc-service


Simple go [gRPC](https://grpc.io/) service for Cloud Run.

> Note, to keep this readme short, I will be asking you to execute scripts rather than listing here complete commands. You should really review each one of these scripts for content, and, to understand the individual commands so you can use them in the future.

## Pre-requirements

### Setup

If you don't have one already, start by creating new project and configuring your [Google Cloud SDK](https://cloud.google.com/sdk/docs/). In addition, you will also need to install the `alpha` component.

```shell
gcloud components install alpha
```

Also, if you have not done so already, you will have [set up Cloud Run](https://cloud.google.com/run/docs/setup).

Finally, to start, clone this repo locally:

```shell
git clone https://github.com/GoogleCloudPlatform/golang-samples.git
```

And navigate into the sample directory:

```shell
cd golang-samples/run/grpc-service
```

### Config

All the variables used in this sample are defined in the [bin/config](bin/config) file. You can edit these to your preferred values.

* `SERVICE` - (default "grpc-service") - this is the name of the service. It's used to build image, create service account, and in Cloud Run deployment. If you do need to change that, edit it before any other step.
* `SERVICE_REGION` - (default "us-central1") - this is the GCP region to which you want to deploy this service. For complete list of regions where Cloud Run service is currently offered see https://cloud.google.com/about/locations/
* `SERVICE_VERSION` - (default "v1") - this is the version of the service which is used to tag the container image as well as the Cloud Run service

## API Definition

First, you need to define the API, the payload shape and the methods that will be used to communicate between client and the server. A sample `proto` file is located in [api/v1/message.proto](api/v1/message.proto), edit it as necessary. You can learn more about [Protocol Buffers](https://developers.google.com/protocol-buffers/) (Protobuf), language-neutral mechanism for serializing structured data, [here](https://github.com/golang/protobuf).

To generate the `go` code from that `proto` run [bin/api](bin/api) script:

```shell
bin/api
```

> Make sure the `--proto_path` flag in [bin/api](bin/api) script points to your local proto dependencies. More on how to set this up [here](https://github.com/golang/protobuf).

As a result, you should now have a new `go` file titled [pkg/api/v1/message.pb.go](pkg/api/v1/message.pb.go). You can review that file but don't edit it as it will be overwritten the next time you run the [bin/api](bin/api) script.

## Server

### Container Image

> If you made any chances to the [api/v1/message.proto](api/v1/message.proto) file you will have to edit both, the [cmd/server/main.go](cmd/server/main.go) and the [cmd/client/main.go](cmd/client/main.go) files so that they compile.

Next, build the server container image which will be used to deploy Cloud Run service using the [bin/image](bin/image) script:

```shell
bin/image
```

### Service Account

While not necessary, as a good principle to follow, we are going to create a specific service account for this service and assign it all the necessary roles using the [bin/user](bin/user) script:

```shell
bin/user
```

### Service Deployment

Once the container image and service account are ready, you can now deploy the new service using [bin/deploy](bin/deploy) script:

```shell
bin/deploy
```

## Client

## Build Client

To invoke the deployed Cloud Run service, build the gRPC client using the [bin/client](bin/client) script:

```shell
bin/client
```

The resulting CLI will be compiled into the `bin` directory. The output of the [bin/client](bin/client) script will also print out the two ways you can execute that client:

```shell
Client CLI generated.
Usage:
 Unary Request/Unary Response
 bin/cli --server grpc-sample-***-uc.a.run.app:443 --message hi

 Unary Request/Stream Response
 bin/cli --server grpc-sample-***-uc.a.run.app:443 --message hi --stream 5
```

## Run

### Unary Execution

When executing the built CLI in unary way (by not including the `--stream` flag):

```shell
bin/cli --server grpc-sample-***-uc.a.run.app:443 --message hi
```

you will see the details of the sent and received message

```shell
Unary Request/Unary Response
 Sent:
  hi
 Response:
  content:<index:1 message:"hi" received_on:<seconds:1567098976 nanos:535796117 > >
```

### Stream Execution

Where as executing it using stream (with `--stream` number)

```shell
bin/cli --server grpc-sample-***-uc.a.run.app:443 --message hi --stream 5
```

the CLI will print the sent message index and server processing time:

```shell
Unary Request/Stream Response
  Stream[1] - Server time: 2019-08-29T17:16:22.837297811Z
  Stream[2] - Server time: 2019-08-29T17:16:22.837928885Z
  Stream[3] - Server time: 2019-08-29T17:16:22.83794915Z
  Stream[4] - Server time: 2019-08-29T17:16:22.837959711Z
  Stream[5] - Server time: 2019-08-29T17:16:22.837968925Z
```

Thant's it. Hope it was helpful.