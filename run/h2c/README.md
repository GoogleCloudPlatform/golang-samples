# Cloud Run HTTP/2 Sample Server

This Go application serves requests using only the HTTP/2 cleartext (h2c)
protocol (it does not support upgrading from HTTP/1.1).

It is provided for testing end-to-end HTTP/2 capabilities on Cloud Run.

## Deploying to Cloud Run

Click this button:

[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)

Alternatively, you can use the `gcloud` SDK:

```sh
git clone https://github.com/GoogleCloudPlatform/golang-samples.git
cd ./golang-samples/run/h2c
gcloud beta run deploy http2-test --use-http2 --source=.
```

**Cleanup:** Remove the `http2-test` Service you deployed from Cloud Run
using the [Cloud Console](https://console.cloud.google.com/run).

## Testing locally

1. Check out this repository and navigate to this directory

1. Start the server locally:

    ```sh
    go run .
    ```

## Verifying h2c

1. After the server starts, make a request using `curl` and use
`--http2-prior-knowledge` to prevent upgrading from HTTP/1.1 (which is
intentionally not supported by this program).

    ```sh
    curl -v --http2-prior-knowledge localhost:8080
    ```

    **Note:** If you have deployed on Cloud Run, you can replace
   `localhost:8080` with your Cloud Run service URL.

1. The `curl` output should indicate HTTP/2 is used:

    ```text
    < HTTP/2 200
    < content-type: text/plain; charset=utf-8
    < content-length: 32
    < date: Fri, 08 Jan 2021 18:13:06 GMT

    This request is served over HTTP/2.0 protocol!
    ```

Since this server intentionally doesn't support upgrading from HTTP/1.1,
performing the same query without the `--http2-prior-knowledge` option makes
curl to use HTTP/1, which fails as expected.
