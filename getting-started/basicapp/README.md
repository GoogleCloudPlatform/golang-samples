# Building, Testing, and Optimizing Code on the Google Cloud
This example demonstrates a web app that sends messages from a user to her or
his friends. It demonstrates the use of Google [Cloud
Build](https://cloud.google.com/cloud-build/) for building and testing portable
Go code with an application that can be run locally, in unit tests, in 
[App Engine Flex](https://cloud.google.com/appengine/docs/flexible/), a plain
Docker container, and in Kubernetes. The app is instrumented for 
[Stackdriver Profiler](https://cloud.google.com/profiler/) and 
[Stackdriver Trace](https://cloud.google.com/trace/).


## Project Setup
The prerequisites for running the example locally are 
[Go](https://golang.org/dl/), the 
[gcloud](https://cloud.google.com/sdk/gcloud/) command line tool, and Docker
community edition. If you prefer not to install gcloud and Docker locally you
can use the [Cloud Shell](https://cloud.google.com/shell/docs/) instead. 

To run this example, you will need to create a Google Cloud project with billing
enabled. You will also need to enable the Cloud SQL API, Cloud SQL Admin API,
the Stackdriver Profiling API, and the Stackdriver Trace API.

If you are trying this out in your own workspace you can get the code with the
git clone command
```
git clone https://github.com/GoogleCloudPlatform/golang-samples.git
cd getting-started/basicapp
```

If you want to try changing code then you may want to save your code in a
private [Google Cloud Source
repository](https://cloud.google.com/source-repositories/docs/) with the
commands
```
gcloud source repos create basicapp
git add .
git commit -m "My private copy"
```

## Local development, mocks, and unit testing
You can run locally and unit test using a mock of the messaging service.

### Running locally
To run the app locally commands export an environment variable for use of mocks
instead of a real database:
```
export MESSAGE_SERVICE=mock
go get -d -v ./...
go run basicwebapp.go
```

Check that the application successfully responds to a HTTP request:
```
curl -I http://localhost:8080
```

You can now use curl to send a message to your friend a message like
```
curl "http://localhost:8080/send?user=Friend1&friend=Friend2&text=We+miss+you!"
```

Your friend can also check their messages with a URL like:
```
curl "http://localhost:8080/messages?user=Friend2"
```

With a mock service we lose the messages as soon as the app is stopped. Unset
the environment variable for use of mocks with the command
```
unset MESSAGE_SERVICE
```

## Unit tests
To run unit tests use the two commands, for each of the main and services
packages:
```
go test
go test -v github.com/GoogleCloudPlatform/golang-samples/getting-started/basicapp/services
```

Check that the tests pass. The Go
[httptest](https://golang.org/pkg/net/http/httptest/) package is used to
simulate HTTP requests and writers in some of the unit tests.

Run the build and execute unit tests using Cloud Build
```
gcloud builds submit --config build/cb-unittest.yaml .
```

View the build and test results in the cloud console under the Cloud Build menu.

## Database
### Setup
This section is based on the instructions at [Using Cloud SQL with
Go](https://cloud.google.com/go/getting-started/using-cloud-sql). Use the
commands below from the Cloud Shell. Make sure that the Cloud SQL API and the
Cloud SQL Admin API are enabled.

Create a [Second Generation Cloud SQL 
instance](https://cloud.google.com/sql/docs/mysql/create-instance) using the 
commands below, naming the instance 'basicapp'.
```
INSTANCE_NAME=basicapp
MACHINE_TYPE=db-n1-standard-1
REGION=us-central1
gcloud sql instances create $INSTANCE_NAME --tier=$MACHINE_TYPE \
  --region=$REGION
```

To set the root password, use the cloud console or the command
```
gcloud sql users set-password root % --instance basicapp --password [PASSWORD]  
```

To see the connection details for your instance use the command
```
gcloud sql instances describe $INSTANCE_NAME
```

Look for the ```connectionName```, which as the form 
```$PROJECT:$REGION:$INSTANCE_NAME```. Also, look for the IP address of the
instance.

Connect to the your Cloud SQL instance with the command in the cloud shell
```
INSTANCE_NAME=basicapp
gcloud sql connect $INSTANCE_NAME --user=root
```

Execute the statements in [data/database_setup.sql](data/database_setup.sql).

### Working with the Database in a Local Development Environment (Optional)

To install the Cloud SQL proxy locally follow the steps in [Using Cloud SQL
with Go](https://cloud.google.com/go/getting-started/using-cloud-sql), for
example on Mac, use the commands

```
curl -o cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.darwin.386
chmod +x cloud_sql_proxy
```

Start the proxy with the commands

```
PROJECT=[your project]
REGION=[your region]
INSTANCE_NAME=basicapp
CONNECTION_NAME=$PROJECT:$REGION:$INSTANCE_NAME
./cloud_sql_proxy -instances="$CONNECTION_NAME"=tcp:3306
```

Run the web app in a separate terminal
```
DBPASSWORD=[your password]
DBUSER=proxyuser
DBHOST=127.0.0.1
DATABASE=messagesdb
export MYSQL_CONNECTION: "$DBUSER:$DBPASSWORD@tcp($DBHOST:3306)/$DATABASE"
go run cmd/basicwebapp/basicwebapp.go
```

You can now use the same curl commands as above to send your friend a message.

## Packaging in a Docker Container
The basicapp example can package your application in a Docker container based
on the [golang image](https://hub.docker.com/_/golang/), with the commands
```
docker build -f build/docker/Dockerfile -t basicwebapp-image .
docker run -itd --rm --name basicwebapp \
  -p 8080:8080 \
  --env MESSAGE_SERVICE=mock \
  basicwebapp-image
```

This will run the container with the mock messaging service. We will describe
how to configure the app to connect to Cloud SQL below.

Check that the application is running OK use the same curl commands as above.

Stop the app with the command
```
docker stop basicwebapp
```

Build the Docker container with cloud build:
```
gcloud builds submit --config build/cb-docker.yaml .
```

Run in Docker, test that it works, and stop it
```
gcloud auth configure-docker
PROJECT_ID=[your project]
docker run -itd --rm -p 8080:8080 \
  --name basicwebapp \
  --env MESSAGE_SERVICE=mock \
  gcr.io/$PROJECT_ID/basicwebapp
docker ps
curl -I http://localhost:8080
docker stop basicwebapp
```

## Deploy the app to Flex
To deploy the app to Flex with the database dependency mocked out use the
command
```
gcloud app deploy app-mock.yaml
```

Check that the app is being served by App Engine with the command
```
gcloud app browse
```

To run the app together with the Cloud SQL database follow steps as in [Using 
Cloud SQL for 
MySQL](https://cloud.google.com/appengine/docs/flexible/go/using-cloud-sql) for
Flex. You can do this by editing the app.yaml file, filling in the text for the
password. Then redeploy the app using the command
```
gcloud app deploy .
```

You will want to avoid leaving the project id and password in the app
configuration file, especially if you are adding the file to a code repository.
Revert the values back to their original symbols by editing app.yaml. To deploy
with Cloud Build, set the strings as environment variables using the command
```
DB_PASSWORD=[user db password]
gcloud builds submit \
  --substitutions=_DB_PASSWORD=$DB_PASSWORD \
  --config build/cb-deploy.yaml .
```

The [build/cb-deploy.yaml](build/cb-deploy.yaml) file contains a [custom build
step](https://cloud.google.com/cloud-build/docs/create-custom-build-steps) with
shell commands to replace the connection string settings in app.yaml with 
[user defined
substitutions](https://cloud.google.com/cloud-build/docs/configuring-builds/substitute-variable-values).

Check that the app is being served by App Engine by pointing your browser at
$PROJECT_ID.appspot.com.
```
curl -I $PROJECT_ID.appspot.com
```

## Testing and Performance Optimization
The next activities in your development lifecycle probably include integration
and load testing, profiling, and tracing.

### Integration and Load Testing

Our integration test uses curl to make sure that the application is properly
configured to talk to the database and test cases for sending and retrieving
messages are ok. 

```
gcloud builds submit --config build/cb-e2etest.yaml .
```

### Profiling
Now that we have generate some load on the app, we can look at the performance
profile. You can profile the code with [Stackdriver
Profiler](https://cloud.google.com/profiler/docs/) to identify areas
that can be optimized. The steps here are adapted from the page [Profiling Go
Code](https://cloud.google.com/profiler/docs/profiling-go).

If you have not enabling the Stackdriver Profiler API you can do that with this
command:
```
gcloud services enable cloudprofiler.googleapis.com
```

A few lines of code from package cloud.google.com/go/profiler have been added
near the entry point of the application to enable profiling. Navigate to the
Profiler menu in the Cloud Console to see the results.

### Tracing

The application has been instrumented for trace data collection with
[OpenCensus](https://opencensus.io/), with export [Stackdriver Cloud
Trace](https://opencensus.io/)for analysis and viewing of trace data.
The method here was adapted from the [OpenCensus Stackdriver
Codelab](https://opencensus.io/codelabs/stackdriver/#0). To see the trace data
navigate to the Trace menu in the Cloud Console. Because of the sparse sampling
of trace data, you many need to run a few iterations of the end-to-end to
generate sufficient data to be displayed.

## Deploy to a Kubernetes cluster
Use the command below to create a Kubernetes cluster:
```
CLUSTER_NAME=mycluster
ZONE=[Your zone]
gcloud container clusters create $CLUSTER_NAME \
  --num-nodes=1 \
  --zone=$ZONE
```

The instructions here to configure the Cloud SQL connection are adapted from 
[Connecting from Kubernetes
Engine](https://cloud.google.com/sql/docs/mysql/connect-kubernetes-engine). To 
connect to Cloud SQL from GKE you need to create and set access for a service
account. Create a service account from the Cloud Console add the role
Cloud SQL Editor, and download the JSON key file. Rename the file 
```credentials.json``` and place it in the top level of the project directory. 
Add a line for credentials.json to your ```.gitignore``` file to prevent it
from being checked into Git.

Create the instance credential secret
```
kubectl create secret generic cloudsql-instance-credentials \
  --from-file=credentials.json=credentials.json
```

Create the connection credential secret
```
kubectl create secret generic cloudsql-db-credentials \
 --from-literal=username=proxyuser \
 --from-literal=password=[PASSWORD]
```

Use Cloud Build to deploy the app to the cluster:
```
gcloud builds submit --config build/cb-k8s.yaml .
```

This will build and push a new container image, create a Kubernetes deployment, 
expose the deployment as a service of type NodePort, and create a Kubernetes
ingress.  Most of the details of the deployment are included in the files
deployment/k8s-deployment.yaml, deployment/k8s-service.yaml, and 
deployment/k8s-ingress.yaml.

Check that the app is running and that the service and ingress are available
```
kubectl get pods
kubectl get services
kubectl get ingress
```

You may need to wait for a few minutes before an external IP is assigned.
Note the external IP address and call it with curl:
```
IP=[Your IP]
curl -I "http://$IP/send?user=Friend1&friend=Friend2&text=We+miss+you!"
curl -I "http://$IP/messages?user=Friend2"
```

Run a test against the Kubernetes service using Cloud Build
```
export IP=[your ip
gcloud builds submit \
  --substitutions=_IP=$IP \
  --config build/cb-e2e-k8s.yaml .
```

## More
For more on Google Cloud Build see the recording of the [Cloud
Build](https://cloud.google.com/cloud-build/)
documentation page. Also see the [Go
Bookshelf](https://github.com/GoogleCloudPlatform/golang-samples/tree/master/getting-started/bookshelf)
example app for the use of a number of storage options with Go.
