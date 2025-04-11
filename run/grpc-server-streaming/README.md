# Cloud Run gRPC Server Streaming sample application

For a general introduction to gRPC on Go, see the [official gRPC
quickstart](https://grpc.io/docs/quickstart/go/).

This sample presents:

- `./server`: a gRPC server application with an RPC that streams the current
  time in the response (written in Go)
- `./client`: a small program to query the server and show the
  response messages (written in Go)

## Deploy gRPC server to Cloud Run

Use the following button to deploy the application to Cloud Run on your GCP
project:

[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run?dir=run/grpc-server-streaming)

Once deployed successfully, note the domain name (hostname) of the service.

## Make a request using the client

1. Ensure Go 1.13 (or higher) is installed on your machine.

2. Clone this repository.

3. Run the client with the hostname of the deployed service:

    ```sh
   cd grpc-server-streaming
   ```

   ```sh
   go run ./client -duration 5 -server <HOSTNAME>:443
    ```

4. Observe the output, it should be receiving and printing a message
   every second as the server sends them.

    ```sh
   rpc established to timeserver, starting to stream
   received message: current_timestamp: 2020-01-15T01:12:29Z
   received message: current_timestamp: 2020-01-15T01:12:30Z
   received message: current_timestamp: 2020-01-15T01:12:31Z
   received message: current_timestamp: 2020-01-15T01:12:32Z
   received message: current_timestamp: 2020-01-15T01:12:33Z
   end of stream
    ```

## Cleanup

Remove the `grpc-server-streaming` Service you deployed from Cloud Run
using the [Cloud Console](https://console.cloud.google.com/run).
