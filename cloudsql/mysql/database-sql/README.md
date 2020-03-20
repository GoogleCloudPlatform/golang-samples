# Connecting to Cloud SQL (MySQL) from a Go web app

This repo contains the Go source code for a simple web app that can be deployed to App Engine Standard. It is a demonstration of how to connect to a MySQL instance in Cloud SQL. The application is the "Tabs vs Spaces" web app used in the [Building Stateful Applications With Kubernetes and Cloud SQL](https://www.youtube.com/watch?v=qVgzP3PsXFw&t=1833s) session at Cloud Next '19.

## Before you begin

1. If you haven't already, set up a Go Development Environment by following the [Go setup guide](https://cloud.google.com/go/docs/setup) and 
[create a project](https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project).

1. Create a 2nd Gen Cloud SQL Instance by following these 
[instructions](https://cloud.google.com/sql/docs/mysql/create-instance). Note the connection string,
database user, and database password that you create.

1. Create a database for your application by following these 
[instructions](https://cloud.google.com/sql/docs/mysql/create-manage-databases). Note the database
name. 

1. Create a service account with the 'Cloud SQL Client' permissions by following these 
[instructions](https://cloud.google.com/sql/docs/mysql/connect-external-app#4_if_required_by_your_authentication_method_create_a_service_account).
Download a JSON key to use to authenticate your connection.

## Running locally

To run this application locally, download and install the `cloud_sql_proxy` by
following the instructions
[here](https://cloud.google.com/sql/docs/mysql/sql-proxy#install).

Instructions are provided below for using the proxy with a TCP connection or a Unix Domain Socket. On Linux or Mac OS you can use either option, but on Windows the proxy currently requires a TCP connection.

### Launch proxy with TCP

To run the sample locally with a TCP connection, set environment variables and launch the proxy as shown below.

#### Linux / Mac OS
Use these terminal commands to initialize environment variables:
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service/account/key.json
export DB_TCP_HOST='127.0.0.1:3306'
export DB_USER='<DB_USER_NAME>'
export DB_PASS='<DB_PASSWORD>'
export DB_NAME='<DB_NAME>'
```

Then use this command to launch the proxy in the background:
```bash
./cloud_sql_proxy -instances=<project-id>:<region>:<instance-name>=tcp:3306 -credential_file=$GOOGLE_APPLICATION_CREDENTIALS &
```

#### Windows/PowerShell
Use these PowerShell commands to initialize environment variables:
```powershell
$env:GOOGLE_APPLICATION_CREDENTIALS="<CREDENTIALS_JSON_FILE>"
$env:DB_TCP_HOST="127.0.0.1:3306"
$env:DB_USER="<DB_USER_NAME>"
$env:DB_PASS="<DB_PASSWORD>"
$env:DB_NAME="<DB_NAME>"
```

Then use this command to launch the proxy in a separate PowerShell session:
```powershell
Start-Process -filepath "C:\<path to proxy exe>" -ArgumentList "-instances=<project-id>:<region>:<instance-name>=tcp:3306 -credential_file=<CREDENTIALS_JSON_FILE>"
```

### Launch proxy with Unix Domain Socket
NOTE: this option is currently only supported on Linux and Mac OS. Windows users should use the [Launch proxy with TCP](#launch-proxy-with-tcp) option.

To use a Unix socket, you'll need to create the `/cloudsql` directory and give write access to the user running the proxy. Use these commands to create the directory and set permissions:
```bash
sudo mkdir /cloudsql
sudo chown -R $USER /cloudsql
```

Use these terminal commands to initialize environment variables:
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service/account/key.json
export INSTANCE_CONNECTION_NAME='<MY-PROJECT>:<INSTANCE-REGION>:<INSTANCE-NAME>'
export DB_USER='<DB_USER_NAME>'
export DB_PASS='<DB_PASSWORD>'
export DB_NAME='<DB_NAME>'
```

Then use this command to launch the proxy in the background:
```bash
./cloud_sql_proxy -dir=/cloudsql --instances=$INSTANCE_CONNECTION_NAME --credential_file=$GOOGLE_APPLICATION_CREDENTIALS &
```

### Testing the application

To test the application locally, follow these steps after the proxy is running:

* Install dependencies: `go get ./...`
* Run the application: `go run cloudsql.go`
* Navigate to `http://127.0.0.1:8080` in a web browser to verify your application is running correctly.

## Deploying to App Engine Standard

To run the sample on GAE-Standard, create an App Engine project by following the setup for these 
[instructions](https://cloud.google.com/appengine/docs/standard/go/quickstart#before-you-begin).

First, create an `app.yaml` with the correct values to pass the environment 
variables into the runtime. Your app.yaml file should look like this:

```yaml
runtime: go111
env_variables:
  INSTANCE_CONNECTION_NAME: <project-id>:<region>:<instance-name>
  DB_USER: YOUR_DB_USER
  DB_PASS: YOUR_DB_PASS
  DB_NAME: YOUR_DB
```

Note: Saving credentials in environment variables is convenient, but not secure - consider a more
secure solution such as [Cloud KMS](https://cloud.google.com/kms/) to help keep secrets safe.

Next, the following command will deploy the application to your Google Cloud project:
```bash
gcloud app deploy
```

## Deploy to Google App Engine Flexible

First, update `app.flexible.yaml` with the correct values to pass the environment 
variables into the runtime.

Next, the following command will deploy the application to your Google Cloud project:
```bash
gcloud app deploy app.flexible.yaml
```

To launch your browser and view the app at https://[YOUR_PROJECT_ID].appspot.com, run the following
command:
```bash
gcloud app browse
```

## Deploy to Cloud Run

See the [Cloud Run documentation](https://cloud.google.com/sql/docs/mysql/connect-run)
for more details on connecting a Cloud Run service to Cloud SQL.

1. Build the container image:

```sh
gcloud builds submit --tag gcr.io/[YOUR_PROJECT_ID]/run-sql
```

2. Deploy the service to Cloud Run:

```sh
gcloud run deploy run-sql --image gcr.io/[YOUR_PROJECT_ID]/run-sql
```

Take note of the URL output at the end of the deployment process.

3. Configure the service for use with Cloud Run

```sh
gcloud beta run services update run-sql \
    --add-cloudsql-instances [INSTANCE_CONNECTION_NAME] \
    --update-env-vars INSTANCE_CONNECTION_NAME=[INSTANCE_CONNECTION_NAME] \
    --update-env-vars DB_USER=[YOUR_DB_USER] \
    --update-env-vars DB_PASS=[YOUR_DB_PASS] \
    --update-env-vars DB_NAME=[YOUR_DB]
```

Replace environment variables with the correct values for your Cloud SQL
instance configuration.

This step can be done as part of deployment but is separated for clarity.

4. Navigate your browser to the URL noted in step 2.

For more details about using Cloud Run see http://cloud.run.
Review other [Go on Cloud Run samples](../../../run/).
