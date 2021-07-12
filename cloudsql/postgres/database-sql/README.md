# Connecting to Cloud SQL (Postgres) from a Go web app

This repo contains the Go source code for a simple web app that can be deployed to App Engine Standard. It is a demonstration of how to connect to a Postgres instance in Cloud SQL. The application is the "Tabs vs Spaces" web app used in the [Building Stateful Applications With Kubernetes and Cloud SQL](https://www.youtube.com/watch?v=qVgzP3PsXFw&t=1833s) session at Cloud Next '19.

## Before you begin

1. If you haven't already, set up a Go Development Environment by following the [Go setup guide](https://cloud.google.com/go/docs/setup) and
   [create a project](https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project).

1. Create a Cloud SQL for Postgres instance by following these
   [instructions](https://cloud.google.com/sql/docs/postgres/create-instance).
   Note the connection string, database user, and database password that you create.

1. Create a database for your application by following these
   [instructions](https://cloud.google.com/sql/docs/postgres/create-manage-databases).
   Note the database name.

1. Create a service account with the 'Cloud SQL Client' permissions by following these
   [instructions](https://cloud.google.com/sql/docs/mysql/connect-external-app#4_if_required_by_your_authentication_method_create_a_service_account).
   Download a JSON key to use to authenticate your connection.

## Running locally

To run this application locally, download and install the `cloud_sql_proxy` by
following the instructions
[here](https://cloud.google.com/sql/docs/postgres/sql-proxy#install).

Instructions are provided below for using the proxy with a TCP connection or a Unix Domain Socket. On Linux or Mac OS you can use either option, but on Windows the proxy currently requires a TCP connection.

### Launch proxy with TCP

To run the sample locally with a TCP connection, set environment variables and launch the proxy as shown below.

#### Linux / Mac OS

Use these terminal commands to initialize environment variables:

```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service/account/key.json
export DB_HOST='127.0.0.1'
export DB_PORT='5432'
export DB_USER='<DB_USER_NAME>'
export DB_PASS='<DB_PASSWORD>'
export DB_NAME='<DB_NAME>'
```

Then use this command to launch the proxy in the background:

```bash
./cloud_sql_proxy -instances=<project-id>:<region>:<instance-name>=tcp:5432 -credential_file=$GOOGLE_APPLICATION_CREDENTIALS &
```

#### Windows/PowerShell

Use these PowerShell commands to initialize environment variables:

```powershell
$env:GOOGLE_APPLICATION_CREDENTIALS="<CREDENTIALS_JSON_FILE>"
$env:DB_HOST="127.0.0.1"
$env:DB_PORT="5432"
$env:DB_USER="<DB_USER_NAME>"
$env:DB_PASS="<DB_PASSWORD>"
$env:DB_NAME="<DB_NAME>"
```

Then use this command to launch the proxy in a separate PowerShell session:

```powershell
Start-Process -filepath "C:\<path to proxy exe>" -ArgumentList "-instances=<project-id>:<region>:<instance-name>=tcp:5432 -credential_file=<CREDENTIALS_JSON_FILE>"
```

### Launch proxy with Unix Domain Socket

NOTE: this option is currently only supported on Linux and Mac OS. Windows users should use the [Launch proxy with TCP](#launch-proxy-with-tcp) option.

To use a Unix socket, you'll need to create a directory and give write access to the user running the proxy. For example:

```bash
sudo mkdir ./cloudsql
sudo chown -R $USER ./cloudsql
```

You'll also need to initialize an environment variable containing the directory you just created:

```bash
export DB_SOCKET_DIR=./cloudsql
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
./cloud_sql_proxy -dir=$DB_SOCKET_DIR --instances=$INSTANCE_CONNECTION_NAME --credential_file=$GOOGLE_APPLICATION_CREDENTIALS &
```

### Testing the application

To test the application locally, follow these steps after the proxy is running:

- Install dependencies: `go get ./...`
- Run the application: `go run cloudsql.go`
- Navigate to `http://127.0.0.1:8080` in a web browser to verify your application is running correctly.

## Deploying to App Engine Standard

To run the sample on GAE-Standard, create an App Engine project by following the setup for these
[instructions](https://cloud.google.com/appengine/docs/standard/go/quickstart#before-you-begin).

First, update `app.standard.yaml` with the correct values to pass the environment
variables into the runtime. Your `app.standard.yaml` file should look like this:

```yaml
runtime: go113
env_variables:
  INSTANCE_CONNECTION_NAME: <project-id>:<region>:<instance-name>
  DB_USER: YOUR_DB_USER
  DB_PASS: YOUR_DB_PASS
  DB_NAME: YOUR_DB
```

Note: Saving credentials in environment variables is convenient, but not secure - consider a more
secure solution such as [Cloud Secret Manager](https://cloud.google.com/secret-manager) to help keep secrets safe.

Next, the following command will deploy the application to your Google Cloud project:

```bash
gcloud app deploy app.standard.yaml
```

## Deploying to App Engine Flexible

To run the sample on GAE-Flex, create an App Engine project by following the setup for these
[instructions](https://cloud.google.com/appengine/docs/standard/go/quickstart#before-you-begin).

First, update `app.flexible.yaml` with the correct values to pass the environment
variables into the runtime. Your `app.flexible.yaml` file should look like this:

```yaml
runtime: custom
env: flex

env_variables:
  INSTANCE_CONNECTION_NAME: <project>:<region>:<instance>
  DB_USER: <your_database_username>
  DB_PASS: <your_database_password>
  DB_NAME: <your_database_name>

beta_settings:
  cloud_sql_instances: <project>:<region>:<instance>
```

Note: Saving credentials in environment variables is convenient, but not secure - consider a more
secure solution such as [Cloud Secret Manager](https://cloud.google.com/secret-manager) to help keep secrets safe.

Next, the following command will deploy the application to your Google Cloud project:

```bash
gcloud app deploy app.flexible.yaml
```

## Deploy to Cloud Run

See the [Cloud Run documentation](https://cloud.google.com/sql/docs/postgres/connect-run)
for more details on connecting a Cloud Run service to Cloud SQL.

1. Build the container image:

```sh
gcloud builds submit --tag gcr.io/[YOUR_PROJECT_ID]/run-sql
```

1. Deploy the service to Cloud Run:

```sh
gcloud run deploy run-sql --image gcr.io/[YOUR_PROJECT_ID]/run-sql \
  --add-cloudsql-instances '<MY-PROJECT>:<INSTANCE-REGION>:<INSTANCE-NAME>' \
  --set-env-vars INSTANCE_CONNECTION_NAME='<MY-PROJECT>:<INSTANCE-REGION>:<INSTANCE-NAME>' \
  --set-env-vars DB_USER='<DB_USER_NAME>' \
  --set-env-vars DB_PASS='<DB_PASSWORD>' \
  --set-env-vars DB_NAME='<DB_NAME>'
```

Take note of the URL output at the end of the deployment process.

Replace environment variables with the correct values for your Cloud SQL
instance configuration.

It is recommended to use the [Secret Manager integration](https://cloud.google.com/run/docs/configuring/secrets) for Cloud Run instead
of using environment variables for the SQL configuration. The service injects the SQL credentials from
Secret Manager at runtime via an environment variable.

Create secrets via the command line:
```sh
echo -n $INSTANCE_CONNECTION_NAME | \
    gcloud secrets create [CLOUD_SQL_CONNECTION_NAME_SECRET] --data-file=-
```

Deploy the service to Cloud Run specifying the env var name and secret name:
```sh
gcloud beta run deploy SERVICE --image gcr.io/[YOUR_PROJECT_ID]/run-sql \
    --add-cloudsql-instances $INSTANCE_CONNECTION_NAME \
    --update-secrets CLOUD_SQL_CONNECTION_NAME=[CLOUD_SQL_CONNECTION_NAME_SECRET]:latest,\
      DB_USER=[DB_USER_SECRET]:latest, \
      DB_PASS=[DB_PASS_SECRET]:latest, \
      DB_NAME=[DB_NAME_SECRET]:latest
```

4. Navigate your browser to the URL noted in step 2.

For more details about using Cloud Run see http://cloud.run.
