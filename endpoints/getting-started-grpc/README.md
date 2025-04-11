# Google Cloud Endpoints Sample for Go using gRPC

This sample demonstrates how to use Google Cloud Endpoints using Go and gRPC.

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
    protoc \
        --include_imports \
        --include_source_info \
        --descriptor_set_out out.pb \
        helloworld/helloworld.proto
    ```

1. Edit `api_config.yaml`. Replace `YOUR_PROJECT_ID`:

    ```yaml
    name: hellogrpc.endpoints.YOUR_PROJECT_ID.cloud.goog
    ```

1. Deploy your service:

    ```bash
    gcloud endpoints services deploy out.pb api_config.yaml
    ```

## Building the server's Docker container

Build and tag your gRPC server, storing it in your private container registry:

```bash
gcloud builds submit --tag gcr.io/YOUR_PROJECT_ID/go-grpc-hello:1.0 .
```

## Deploy to GCE or GKE

### Deploy to GCE

1. Create an instance and SSH into it:

    ```bash
    gcloud compute instances create grpc-host --image-family cos-stable --image-project cos-cloud --tags=http-server
    gcloud compute ssh grpc-host
    ```

1. Set some environment variables:

    ```bash
    GOOGLE_CLOUD_PROJECT=$(curl -s "http://metadata.google.internal/computeMetadata/v1/project/project-id" -H "Metadata-Flavor: Google")
    SERVICE_NAME=hellogrpc.endpoints.${GOOGLE_CLOUD_PROJECT}.cloud.goog
    ```

1. Pull your credentials to access your private container registry:

    ```bash
    docker-credential-gcr configure-docker
    ```

1. Create your own container network called esp_net:

    ```bash
    docker network create --driver bridge esp_net
    ```

1. Run your gRPC server's container:

    ```bash
    docker run --detach --net=esp_net --name=grpc-hello gcr.io/${GOOGLE_CLOUD_PROJECT}/go-grpc-hello:1.0
    ```

1. Run Endpoints proxy:

    ```bash
    docker run \
        --detach \
        --name=esp \
        --publish=80:9000 \
        --net=esp_net \
        gcr.io/endpoints-release/endpoints-runtime:1 \
        --service=${SERVICE_NAME} \
        --rollout_strategy=managed \
        --http2_port=9000 \
        --backend=grpc://grpc-hello:50051
    ```

1. On your local machine, get the IP address of your secured gRPC server:

    ```bash
    gcloud compute instances list --filter=grpc-host
    ```

1. Send a request to the API server (see "Running the client" below)

### Deploy to GKE

If you haven't got a cluster, first [create one](https://cloud.google.com/kubernetes-engine/docs/how-to/creating-a-container-cluster).

1. Edit `deployment.yaml`. Replace `<YOUR_PROJECT_ID>`.

1. Create the deployment and service:

    ```
    kubectl apply -f deployment.yaml
    ```

1. Wait until the load balancer is active:

    ```
    kubectl get svc grpc-hello --watch
    ```

1. Send a request to the API server (see "Running the client" below)

## Running the client

1. First, [create a project API key](https://console.developers.google.com/apis/credentials).

1. Then, on your local machine, after you have your server's IP address (via
   GKE's `kubectl get svc` or your GCE instance's IP), run:

    ```bash
    go run client/main.go --api-key=AIza.... --addr=YOUR_SERVER_IP_ADDRESS:80 [message]
    ```

[1]: https://cloud.google.com/endpoints/docs/quickstarts

## Configuring authentication and authenticating requests

### Configuring Authentication

This sample shows how to make requests authenticated by a service account using a signed JWT token.

1. First, [create a service account](https://console.developers.google.com/apis/credentials)

1. Edit `api_config_auth.yaml`. Replace `PROJECT_ID` and `SERVICE-ACCOUNT-ID`.

1. Update the service configuration using `api_config_auth.yaml` instead of `api_config.yaml`:

    ```bash
    gcloud endpoints services deploy out.pb api_config_auth.yaml
    ```

### Authenticating requests

1. First, [create and download a service account key](https://console.developers.google.com/apis/credentials) in JSON format.

1. Then, run:

    ```bash
    go run client/main.go \
        --keyfile=SERVICE_ACCOUNT_KEY.json \
        --audience=hellogrpc.endpoints.PROJECT_ID.cloud.goog \
        --addr=YOUR_SERVER_IP_ADDRESS:80 \
        [message]
    ```
