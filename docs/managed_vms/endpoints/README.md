# Google Cloud Endpoints sample

## Deploying the backend

### Update the `swagger.yaml` configuration file

Edit the `swagger.yaml` file and replace `YOUR-PROJECT-ID` with your project ID.

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

Send an echo request:
```
go run client/main.go -api-key=AIza.... -host=https://my-app.appspot.com -echo message
```

Send a JWT authed request:
```
go run client/main.go -api-key=AIza....  -host=https://my-app.appspot.com -service-account=path_to_service_account.json
```
