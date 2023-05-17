# Connecting to Cloud SQL (SQL Server) from a Go web app

This repo contains the Go source code for a simple web app that can be deployed to App Engine Standard. It is a demonstration of how to connect to a SQL Server instance in Cloud SQL. The application is the "Tabs vs Spaces" web app used in the [Building Stateful Applications With Kubernetes and Cloud SQL](https://www.youtube.com/watch?v=qVgzP3PsXFw&t=1833s) session at Cloud Next '19.

## Before you begin

1. If you haven't already, set up a Go Development Environment by following the [Go setup guide](https://cloud.google.com/go/docs/setup) and
[create a project](https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project).

1. Create a Cloud SQL for SQL Server instance by following these
[instructions](https://cloud.google.com/sql/docs/sqlserver/create-instance).
Note the connection string, database user, and database password that you create.

1. Create a database for your application by following these
[instructions](https://cloud.google.com/sql/docs/sqlserver/create-manage-databases).
Note the database name.

1. Create a service account with the 'Cloud SQL Client' permissions by following these
[instructions](https://cloud.google.com/sql/docs/mysql/connect-external-app#4_if_required_by_your_authentication_method_create_a_service_account).
Download a JSON key to use to authenticate your connection.

## Running locally

To run this application locally, download and install the `cloud_sql_proxy` by
following the instructions
[here](https://cloud.google.com/sql/docs/sqlserver/sql-proxy#install).

Instructions are provided below for using the proxy with a TCP connection or a Unix Domain Socket. On Linux or Mac OS you can use either option, but on Windows the proxy currently requires a TCP connection.

### Launch proxy with TCP

To run the sample locally with a TCP connection, set environment variables and launch the proxy as shown below.

#### Linux / Mac OS
Use these terminal commands to initialize environment variables:
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service/account/key.json
export INSTANCE_HOST='127.0.0.1'
export DB_PORT='1433'
export DB_USER='<YOUR_DB_USER_NAME>'
export DB_PASS='<YOUR_DB_PASSWORD>'
export DB_NAME='<YOUR_DB_NAME>'
```

Then use this command to launch the proxy in the background:
```bash
./cloud_sql_proxy -instances=<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>=tcp:1433 -credential_file=$GOOGLE_APPLICATION_CREDENTIALS &
```

#### Windows/PowerShell
Use these PowerShell commands to initialize environment variables:
```powershell
$env:GOOGLE_APPLICATION_CREDENTIALS="<CREDENTIALS_JSON_FILE>"
$env:INSTANCE_HOST="127.0.0.1"
$env:DB_PORT="1433"
$env:DB_USER="<YOUR_DB_USER_NAME>"
$env:DB_PASS="<YOUR_DB_PASSWORD>"
$env:DB_NAME="<YOUR_DB_NAME>"
```

Then use this command to launch the proxy in a separate PowerShell session:
```powershell
Start-Process -filepath "C:\<path to proxy exe>" -ArgumentList "-instances=<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>=tcp:1433 -credential_file=<CREDENTIALS_JSON_FILE>"
```

### Testing the application

To test the application locally, follow these steps after the proxy is running:

* Install dependencies: `go get ./...`
* Run the application: `go run cloudsql.go`
* Navigate to `http://127.0.0.1:8080` in a web browser to verify your application is running correctly.

## Deploying to App Engine Standard

To run the sample on GAE-Standard, create an App Engine project by following the setup for these
[instructions](https://cloud.google.com/appengine/docs/standard/go/quickstart#before-you-begin).

First, update [`app.standard.yaml`](cmd/app/app.standard.yaml) with the correct values to pass the environment
variables into the runtime. Your `app.standard.yaml` file should look like this:

```yaml
runtime: go116
env_variables:
  INSTANCE_CONNECTION_NAME: <PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>
  DB_USER: <YOUR_DB_USER_NAME>
  DB_PASS: <YOUR_DB_PASSWORD>
  DB_NAME: <YOUR_DB_NAME>
```

Note: Saving credentials in environment variables is convenient, but not secure - consider a more
secure solution such as [Cloud KMS](https://cloud.google.com/kms/) to help keep secrets safe.

Next, the following command will deploy the application to your Google Cloud project:
```bash
gcloud app deploy cmd/app/app.standard.yaml
```

## Deploying to App Engine Flexible

To run the sample on GAE-Flex, create an App Engine project by following the setup for these
[instructions](https://cloud.google.com/appengine/docs/standard/go/quickstart#before-you-begin).

First, update [`app.flexible.yaml`](app.flexible.yaml) with the correct values to pass the environment
variables into the runtime. Your `app.flexible.yaml` file should look like this:

```yaml
runtime: custom
env: flex

env_variables:
  INSTANCE_CONNECTION_NAME: <PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>
  DB_USER: <YOUR_DB_USER_NAME>
  DB_PASS: <YOUR_DB_PASSWORD>
  DB_NAME: <YOUR_DB_NAME>

beta_settings:
  cloud_sql_instances: <PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>
```

Note: Saving credentials in environment variables is convenient, but not secure - consider a more
secure solution such as [Cloud Secret Manager](https://cloud.google.com/secret-manager) to help keep secrets safe.

Next, the following command will deploy the application to your Google Cloud project:

```bash
gcloud app deploy app.flexible.yaml
```
