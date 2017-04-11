# Google Cloud Endpoints Sample for Go using gRPC

This sample demonstrates how to use Google Cloud Endpoints using Go and gRPC.

For a complete walkthrough, see the following guides:

* [Quickstart with gRPC on Container Engine](https://cloud.google.com/endpoints/docs/quickstart-grpc-container-engine)
* [Quickstart with gRPC on Compute Engine](https://cloud.google.com/endpoints/docs/quickstart-grpc-compute-engine-docker)

## Test the code locally (optional)

Run the backend using `go run`:

```bash
$ go run server/main.go
```

Send a request from another terminal:

```bash
$ go run client/main.go
2017/03/30 17:08:32 Greeting: Hello world
```

## Deploying service config

1. First, generate `out.pb` from the proto file:

    ```bash
    protoc --include_imports --include_source_info helloworld/helloworld.proto --descriptor_set_out out.pb
    ```

1. Edit `api_config.yaml`. Replace `YOUR_PROJECT_ID`:

    ```yaml
    name: hellogrpc.endpoints.YOUR_PROJECT_ID.cloud.goog
    ```

1. Deploy your service:

    ```bash
    gcloud service-management deploy out.pb api_config.yaml
    ```

    Your config ID should be printed out, it looks like `2017-03-30r0`.
    Take a note of it, you'll need it later.

    You can list the config IDs using this command:

    ```bash
    gcloud service-management configs list --service hellogrpc.endpoints.YOUR_PROJECT_ID.cloud.goog
    ```

## Building the server's Docker container

Build and tag your gRPC server, storing it in your private container registry:

```bash
gcloud container builds submit --tag gcr.io/YOUR_PROJECT_ID/go-grpc-hello:1.0 .
```

## Deploy to GCE or GKE

### Deploy to GCE

1. Create an instance and SSH into it:

    ```bash
    gcloud compute instances create grpc-host --image-family gci-stable --image-project google-containers --tags=http-server
    gcloud compute ssh grpc-host
    ```

1. Set some environment variables (you'll need to manually set the service config ID):

    ```bash
    GCLOUD_PROJECT=$(curl -s "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google")
    SERVICE_NAME=hellogrpc.endpoints.${GCLOUD_PROJECT}.cloud.goog
    SERVICE_CONFIG_ID=<Your Config ID>
    ```

1. Pull your credentials to access your private container registry:

    ```bash
    /usr/share/google/dockercfg_update.sh
    ```

1. Run your gRPC server's container:

    ```bash
    docker run -d --name=grpc-hello gcr.io/${GCLOUD_PROJECT}/go-grpc-hello:1.0
    ```

1. Run Endpoints proxy:

    ```bash
    docker run --detach --name=esp \
        -p 80:9000 \
        --link=grpc-hello:grpc-hello \
        gcr.io/endpoints-release/endpoints-runtime:1 \
        -s ${SERVICE_NAME} \
        -v ${SERVICE_CONFIG_ID} \
        -P 9000 \
        -a grpc://grpc-hello:50051
    ```

1. Get the IP address of your secured gRPC server:

    ```bash
    gcloud compute instances list --filter=grpc-host
    ```

1. Send a request to the API server (see "Running the client" below)

### Deploy to GKE

If you haven't got a cluster, first [create one](https://cloud.google.com/container-engine/docs/clusters/operations).

1. Edit `container-engine.yaml`. Replace `<YOUR_PROJECT_ID>` and `<SERVICE_CONFIG_ID>`.

1. Create the deployment and service:

    ```
    kubectl apply -f container-engine.yaml
    ```

1. Wait until the load balancer is active:

    ```
    kubectl get svc grpc-hello --watch
    ```

1. Send a request to the API server (see "Running the client" below)

## Running the client

1. First, [create a project API key](https://console.developers.google.com/apis/credentials).

1. Then, after you have your server's IP address (via GKE's `kubectl get svc` or your GCE instance's IP), run:

    ```bash
    go run client/main.go -api-key=AIza.... -addr=YOUR_SERVER_IP_ADDRESS:80 [message]
    ```

[1]: https://cloud.google.com/endpoints/docs/quickstarts
