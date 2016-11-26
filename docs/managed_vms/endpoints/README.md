# Google Cloud Endpoints sample for Google App Engine flexible environment

## Deploying the backend

### Update the `openapi.yaml` configuration file

Edit the `openapi.yaml` file and replace `YOUR-PROJECT-ID` with your project ID.

The backend can be secured with several authentication schemes:

* Firebase Auth
* Google ID token
* Google JWT (e.g. service account)
* Auth0

Each of those require further configuration.
See the documentation (currently under whitelist) for more information.

### Deploy

First, install `aedeploy`:

    go get -u google.golang.org/appengine/cmd/aedeploy

Deploy the application:

    aedeploy gcloud beta app deploy

## Running the client

### Send an echo request using an API key

First, [create a project API key](https://console.developers.google.com/apis/credentials).

Then, run:

```
go run client/main.go -api-key=AIza.... -host=https://my-app.appspot.com -echo message
```

### Send a request using JWT authentication

First, [download a Service Account JSON key file](https://developers.google.com/identity/protocols/OAuth2ServiceAccount#creatinganaccount).

Then, run:

```
go run client/main.go -host=https://my-app.appspot.com -service-account=path_to_service_account.json
```
