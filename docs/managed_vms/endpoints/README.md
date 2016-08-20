# Google Cloud Endpoints sample

WARNING: This sample does not work with Go 1.7 yet. Please use Go 1.6.

## Updating the project ID

The project ID in the sample must be updated to match your project ID.

Edit the `swagger.yaml` file and replace `YOUR-PROJECT-ID` with your project ID.

## Deploying the backend

To deploy a Go app, use the [aedeploy
tool](https://godoc.org/google.golang.org/appengine/cmd/aedeploy), which will
correctly assemble your app's dependencies in the same way that the go tool
does. From this sample's directory, install aedeploy by running:

    export GOPATH=`pwd`
    go get -u -d github.com/GoogleCloudPlatform/golang-samples/docs/managed_vms/endpoints/...
    go get google.golang.org/appengine/cmd/aedeploy

Deploy your app by appending the `gcloud` command after `aedeploy`:

    bin/aedeploy gcloud beta app deploy

## Running the client

Send an echo request:
```
go run client/main.go -api-key=AIza.... -host=https://my-app.appspot.com -echo message
```

Send a JWT authed request:
```
go run client/main.go -api-key=AIza....  -host=https://my-app.appspot.com -service-account=path_to_service_account.json
```

For more details about auth, check out the documentation for [authenticating
users](https://cloud.google.com/endpoints/docs/authenticating-users).
