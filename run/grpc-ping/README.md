# Cloud Run gRPC Ping Sample Application

For a general introduction to gRPC on Go, see the [official gRPC quickstart](https://grpc.io/docs/quickstart/go/).

This sample presents a single server which is built to one container image.

To demonstrate service-to-service gRPC requests, this container image is deployed as two services: "ping" and "ping-upstream". ping is made public and ping-upstream is the data provider.

## Deploying to Cloud Run

1. Build & Deploy the gRPC services:

   ```sh
   export GOOGLE_CLOUD_PROJECT=[PROJECT_ID]
   # Build and push container images.
   gcloud builds submit --tag gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-ping

   # Deploy ping service for private access.
   gcloud run deploy ping-upstream --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-ping

   # Get the host name of the ping service.
   PING_URL=$(gcloud run services describe ping-upstream --format='value(status.url)')
   PING_DOMAIN=${PING_URL#https://}

   # Deploy ping-relay service for public access.
   gcloud run deploy ping --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-ping \
       --update-env-vars GRPC_PING_HOST=${PING_DOMAIN}:443 \
       --allow-unauthenticated
   ```

2. Make a ping request:

   Use the client CLI to send a request:

   ```sh
   go run ./client -server [RELAY-SERVICE-DOMAIN]:443 -relay -message "Hello Friend"
   ```

If you later make some code changes, updating is more concise:

```sh
export GOOGLE_CLOUD_PROJECT=[PROJECT_ID]
gcloud run deploy ping --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-ping
gcloud run deploy ping-relay --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-ping
```

See below for instructions on updating the proto.

## Environment Variable Configuration Options

* `GRPC_PING_HOST`: [relay: `example.com:443`; required] Ping upstream service host nanme.
* `GRPC_PING_INSECURE`: [relay: `false`] Use an insecure connection to the ping service. Primarily for local development.
* `GRPC_PING_UNAUTHENTICATED`: [relay: `false`] Make unauthenticated requests to the ping service. Primarily for local development.

## Building Locally

```sh
docker build -t grpc-ping .
```

## Running Locally

### Running client &rArr; server ping

```sh
cd ping
go run .
```

Open another terminal at the grpc-ping directory:

```sh
go run ./client -server localhost:8080 -insecure -message "Hello Friend!"
```

ping-j6jtwetqdq-uc.a.run.app

### Running client &rArr; server &rArr; server ping

1. Start the ping service:

   ```sh
   cd ./ping
   PORT=9090 go run .
   ```

2. In another session, start the relay service:

   ```sh
   cd ./relay
   GRPC_PING_INSECURE=1 GRPC_PING_HOST=localhost:9090 GRPC_PING_UNAUTHENTICATED=1 \
       go run .
   ```

3. From the grpc-ping directory use the grpc client to send a request:

   ```sh
   go run ./client -server localhost:8080 -insecure -relay -message "Hello Relayed Friend!"
   ```

## Updating the Proto

1. Retrieve the protoc plugin for Go:

    ```
    go get -u github.com/golang/protobuf/protoc-gen-go
    ```

2. Modify the Protobuf by editing `api/v1/message.proto`.

3. Regenerate the go code:

    ```
    protoc \
        --proto_path api/v1 \
        --go_out "plugins=grpc:pkg/api/v1" \
        message.proto
    ```
