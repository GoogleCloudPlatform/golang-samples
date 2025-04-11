# Google Cloud Endpoints Sample for Go

This sample demonstrates how to use Google Cloud Endpoints using Go.

For a complete walkthrough showing how to run this sample in different
environments, see the [Google Cloud Endpoints Quickstarts][1].

## Running the backend locally

Run the backend using `go run`:

```bash
go run app.go
```

## Deploying the backend to AppEngine Flex

First, edit the `app.yaml` configuration file, setting `endpoints_api_service.name` to your service domain name.

Then, deploy the app with:

```bash
gcloud app deploy
```

## Running the client

### Send an echo request using an API key

First, [create a project API key](https://console.developers.google.com/apis/credentials).

Then, run:

```bash
go run client/main.go -api-key=AIza.... -host=https://my-app.appspot.com -echo message
```

### Send a request using JWT authentication

First, [download a Service Account JSON key file](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#creatinganaccount).

Then, run:

```bash
go run client/main.go -host=https://my-app.appspot.com -service-account=path_to_service_account.json
```

[1]: https://cloud.google.com/endpoints/docs/quickstarts
