# Building, Testing, and Optimizing Code on the Google Cloud
This example demonstrates a web app that sends messages from a user to their
friends. It demonstrates the use of Google [Cloud
Build](https://cloud.google.com/cloud-build/) for building and testing portable
Go code with an application that can be run locally, in unit tests, in 
[App Engine Flex](https://cloud.google.com/appengine/docs/flexible/), and a
Docker container.


## Project Setup
The prerequisites for running the example locally are 

1. [Go](https://golang.org/dl/)

2. the [gcloud](https://cloud.google.com/sdk/gcloud/) command line tool

3. Git

4. Docker community edition.

If you prefer not to install gcloud and Docker locally you
can use [Google Cloud Shell](https://cloud.google.com/shell/docs/). 

To run this example, you will need to create a Google Cloud project with billing
enabled. You will also need to enable the following APIs

1. [App Engine Admin API](https://cloud.google.com/appengine/docs/admin-api/)
   \- [enable](https://console.cloud.google.com/flows/enableapi?apiid=appengine)

2. [App Engine Flexible Environment 
   API](https://cloud.google.com/appengine/docs/flexible/go/)
   \- [enable](https://console.cloud.google.com/flows/enableapi?apiid=appengineflex.googleapis.com)

3. [Cloud SQL Admin API](https://cloud.google.com/sql/docs/mysql/admin-api/)
   \- [enable](https://console.cloud.google.com/flows/enableapi?apiid=sqladmin)

Set your project as the default in your local development environment with the
command
```
PROJECT_ID=[Your project id]
gcloud config set project $PROJECT_ID
```

Get the example code and artifacts with the command
```
git clone https://github.com/GoogleCloudPlatform/golang-samples.git
cd go/src/github.com/GoogleCloudPlatform/golang-samples/getting-started/devflowapp
```

## Local development, mocks, and unit testing
You can run locally and unit test using a mock of the ```MessageService```
interface.

### Running locally
To use a mock database when running the app locally, run::
```
export MESSAGE_SERVICE=mock
go get -d -v ./...
go run devflowapp.go
```

Check that the application successfully responds to a HTTP request:
```
curl -I http://localhost:8080
```

You can also use Curl to send a message to your friend a message, like
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

### Unit tests
To run unit tests use the two commands, for each of the main and services
packages:
```
go test
go test -v github.com/GoogleCloudPlatform/golang-samples/getting-started/devflowapp/services
```

Check that the tests pass. The Go
[httptest](https://golang.org/pkg/net/http/httptest/) package is used to
simulate HTTP requests and writers in some of the unit tests.

To use Cloud Build, run the build and execute unit tests using Cloud Build
```
gcloud builds submit --config build/cb-unittest.yaml .
```

View the build and test results in the cloud console under the [Cloud Build
menu](https://console.cloud.google.com/cloud-build).

## Database
### Setup
This section is based on the instructions at [Using Cloud SQL with
Go](https://cloud.google.com/go/getting-started/using-cloud-sql). 

Create a [Second Generation Cloud SQL 
instance](https://cloud.google.com/sql/docs/mysql/create-instance) using the 
commands below, naming the instance 'devflowapp'.
```
INSTANCE_NAME=devflowapp
MACHINE_TYPE=db-n1-standard-1
REGION=us-central1
gcloud sql instances create $INSTANCE_NAME --tier=$MACHINE_TYPE \
  --region=$REGION
```

To set the root password, use the cloud console or the command
```
gcloud sql users set-password root % --instance devflowapp --password [PASSWORD]  
```

To see the connection details for your instance use the command
```
gcloud sql instances describe $INSTANCE_NAME
```

Look for the ```connectionName```, which as the form 
```$PROJECT:$REGION:$INSTANCE_NAME```. 
Also, look for the IP address of the instance.

Connect to the your Cloud SQL instance with the command in the Cloud Shell
```
INSTANCE_NAME=devflowapp
gcloud sql connect $INSTANCE_NAME --user=root
```

Execute the statements in
[data/database_setup.sql](https://github.com/GoogleCloudPlatform/golang-samples/blob/master/getting-started/devflowapp/data/dastabase_setup.sql).

### Working with the Database in a Local Development Environment (Optional)

To install the Cloud SQL proxy locally follow [Using Cloud SQL
with Go](https://cloud.google.com/go/getting-started/using-cloud-sql), for
example on Mac, use the commands

```
curl -o cloud_sql_proxy https://dl.google.com/cloudsql/cloud_sql_proxy.darwin.386
chmod +x cloud_sql_proxy
```

Start the proxy with the commands

```
PROJECT_ID=[your project]
REGION=[your region]
CONNECTION_NAME=$PROJECT_ID:$REGION:$INSTANCE_NAME
./cloud_sql_proxy -instances="$CONNECTION_NAME"=tcp:3306
```

Run the web app in a separate terminal
```
DBPASSWORD=[your password]
DBUSER=proxyuser
DBHOST=127.0.0.1
DATABASE=messagesdb
export MYSQL_CONNECTION="$DBUSER:$DBPASSWORD@tcp($DBHOST:3306)/$DATABASE"
go run devflowapp.go
```

You can now use the same Curl commands as above to send your friend a message.

## Packaging in a Docker Container
The devflowapp example can package your application in a Docker container based
on the [golang image](https://hub.docker.com/_/golang/), with the commands
```
docker build -f Dockerfile -t devflowapp-image .
docker run -itd --rm --name devflowapp \
  -p 8080:8080 \
  --env MESSAGE_SERVICE=mock \
  devflowapp-image
```

This will run the container with the mock messaging service. We will describe
how to configure the app to connect to Cloud SQL below.

Check that the application is running OK use the same Curl commands as above.

Stop the app with the command
```
docker stop devflowapp
```

Build the Docker container with Cloud Build:
```
gcloud builds submit --config build/cb-docker.yaml .
```

Run in Docker, test that it works, and stop it
```
gcloud auth configure-docker
PROJECT_ID=[your project]
docker run -itd --rm -p 8080:8080 \
  --name devflowapp \
  --env MESSAGE_SERVICE=mock \
  gcr.io/$PROJECT_ID/devflowapp
docker ps
curl -I http://localhost:8080
docker stop devflowapp
```

## Deploy the app to App Engine Flexible
If you have not used App Engine in your project before you will need to enable
it. You can do that with the command
```
gcloud app create --region=us-central
```

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
Flex. You can do this by editing the app.yaml file then filling in the text for the
password and the project id. Then redeploy the app using the command
```
gcloud app deploy .
```

To do the equivalent operation in Cloud Build, follow the steps in [Deploying
artifacts](https://cloud.google.com/cloud-build/docs/configuring-builds/build-test-deploy-artifacts). 
In particular, you will need to grant App Engine Admin role to the Cloud Build
service account.

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

##  Integration and Load Testing

An integration test uses Curl to make sure that the application is properly
configured to talk to the database and test cases for sending and retrieving
messages are working ok end-to-end. 

```
gcloud builds submit --config build/cb-e2etest.yaml .
```

## More
For more on Google Cloud Build see the recording of the [Cloud
Build](https://cloud.google.com/cloud-build/)
documentation page. Also see the [Go
Bookshelf](https://github.com/GoogleCloudPlatform/golang-samples/tree/master/getting-started/bookshelf)
example app for the use of a number of storage options with Go.
