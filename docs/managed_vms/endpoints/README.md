# Google Cloud Endpoints sample

## Deploying the backend

Follow the [instructions][deploy] in the docs.

## Running the client

Send an echo request:
```
go run client/main.go -api-key=AIza.... -host=https://my-app.appspot.com -echo message
```

Send a JWT authed request:
```
go run client/main.go -api-key=AIza....  -host=https://my-app.appspot.com -service-account=path_to_service_account.json
```

[deploy]: https://cloud.google.com/appengine/docs/flexible/go/testing-and-deploying-your-app#deploying_your_program
