# Connecting to Cloud SQL (MySQL) from a Go web app

This repo contains the Go source code for a simple web app that can be deployed to App Engine Standard. It is a demonstration of how to connect to a MySQL instance in Cloud SQL. The application is the "Tabs vs Spaces" web app used in the [Building Stateful Applications With Kubernetes and Cloud SQL](https://www.youtube.com/watch?v=qVgzP3PsXFw&t=1833s) session at Cloud Next '19.

## Before you begin

1. If you haven't already, set up a Go Development Environment by following the [Go setup guide](https://cloud.google.com/go/docs/setup) and
[create a project](https://cloud.google.com/resource-manager/docs/creating-managing-projects#creating_a_project).

1. Create a Cloud SQL for MySQL instance by following these
[instructions](https://cloud.google.com/sql/docs/mysql/create-instance).
Note the connection string, database user, and database password that you create.

1. Create a database for your application by following these
[instructions](https://cloud.google.com/sql/docs/mysql/create-manage-databases).
Note the database name.

1. Set up [Application Default Credentials][adc] and ensure you have
   added the 'Cloud SQL Client' role to your IAM principal.

[adc]: https://cloud.google.com/docs/authentication/provide-credentials-adc

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
export INSTANCE_HOST='127.0.0.1'
export DB_PORT='3306'
export DB_USER='<DB_USER_NAME>'
export DB_PASS='<DB_PASSWORD>'
export DB_NAME='<DB_NAME>'
```

Then use this command to launch the proxy in the background:
```bash
./cloud-sql-proxy <PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME> --port=3306 &
```

#### Windows/PowerShell
Use these PowerShell commands to initialize environment variables:
```powershell
$env:INSTANCE_HOST="127.0.0.1"
$env:DB_PORT="3306"
$env:DB_USER="<YOUR_DB_USER_NAME>"
$env:DB_PASS="<YOUR_DB_PASSWORD>"
$env:DB_NAME="<YOUR_DB_NAME>"
```

Then use this command to launch the proxy in a separate PowerShell session:
```powershell
Start-Process -filepath "C:\<path to proxy exe>" -ArgumentList "<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME> --port=3306"
```

### Launch proxy with Unix Domain Socket
NOTE: this option is currently only supported on Linux and Mac OS. Windows users should use the [Launch proxy with TCP](#launch-proxy-with-tcp) option.

To use a Unix socket, you'll need to create a directory and give write access to the user running the proxy. For example:
```bash
sudo mkdir ./cloudsql
sudo chown -R $USER ./cloudsql
```

Use these terminal commands to initialize environment variables:
```bash
export INSTANCE_UNIX_SOCKET='./cloudsql/<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>'
export DB_USER='<YOUR_DB_USER_NAME>'
export DB_PASS='<YOUR_DB_PASSWORD>'
export DB_NAME='<YOUR_DB_NAME>'
```

Then use this command to launch the proxy in the background:
```bash
./cloud-sql-proxy --unix-socket=./cloudsql \
    <PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME> &
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
  INSTANCE_UNIX_SOCKET: /cloudsql/<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>
  DB_USER: <YOUR_DB_USER_NAME>
  DB_PASS: <YOUR_DB_PASSWORD>
  DB_NAME: <YOUR_DB_NAME>
```

Note: Saving credentials in environment variables is convenient, but not secure - consider a more
secure solution such as [Cloud Secret Manager](https://cloud.google.com/secret-manager) to help keep secrets safe.

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
  INSTANCE_UNIX_SOCKET: /cloudsql/<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>
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
## Deploy to Cloud Run

See the [Cloud Run documentation](https://cloud.google.com/sql/docs/mysql/connect-run)
for more details on connecting a Cloud Run service to Cloud SQL.

1. Build the container image:

```sh
gcloud builds submit --tag gcr.io/[YOUR_PROJECT_ID]/run-sql
```

2. Deploy the service to Cloud Run:

```sh
gcloud run deploy run-sql --image gcr.io/[YOUR_PROJECT_ID]/run-sql \
  --add-cloudsql-instances '<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>' \
  --update-env-vars INSTANCE_UNIX_SOCKET='/cloudsql/<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>' \
  --update-env-vars DB_USER='<DB_USER_NAME>' \
  --update-env-vars DB_PASS='<DB_PASSWORD>' \
  --update-env-vars DB_NAME='<DB_NAME>'
```

Take note of the URL output at the end of the deployment process.

Replace environment variables with the correct values for your Cloud SQL
instance configuration.

It is recommended to use the [Secret Manager integration](https://cloud.google.com/run/docs/configuring/secrets) for Cloud Run instead
of using environment variables for the SQL configuration. The service injects the SQL credentials from
Secret Manager at runtime via an environment variable.

Create secrets via the command line:
```sh
echo -n $INSTANCE_UNIX_SOCKET | \
    gcloud secrets create [INSTANCE_UNIX_SOCKET_SECRET] --data-file=-
```

Deploy the service to Cloud Run specifying the env var name and secret name:
```sh
gcloud beta run deploy SERVICE --image gcr.io/[YOUR_PROJECT_ID]/run-sql \
    --add-cloudsql-instances <PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME> \
    --update-secrets INSTANCE_UNIX_SOCKET=[INSTANCE_UNIX_SOCKET_SECRET]:latest,\
      DB_USER=[DB_USER_SECRET]:latest, \
      DB_PASS=[DB_PASS_SECRET]:latest, \
      DB_NAME=[DB_NAME_SECRET]:latest
```

3. Navigate your browser to the URL noted in step 2.

For more details about using Cloud Run see http://cloud.run.

## Deploy to Cloud Functions

To deploy the service to [Cloud Functions](https://cloud.google.com/functions/docs) run the following command:

```sh
gcloud functions deploy votes --gen2 --runtime go120 --trigger-http \
  --allow-unauthenticated \
  --entry-point Votes \
  --region <INSTANCE_REGION> \
  --set-env-vars INSTANCE_UNIX_SOCKET=/cloudsql/<PROJECT_ID>:<INSTANCE_REGION>:<INSTANCE_NAME> \
  --set-env-vars DB_USER=$DB_USER \
  --set-env-vars DB_PASS=$DB_PASS \
  --set-env-vars DB_NAME=$DB_NAME
```

Note: If the function fails to deploy or returns a `500: Internal service error`,
this may be due to a known limitation with Cloud Functions gen2 not being able
to configure the underlying Cloud Run service with a Cloud SQL connection.

A workaround command to fix this is is to manually revise the Cloud Run
service with the Cloud SQL Connection:

```sh
gcloud run deploy votes --source . \
  --region <INSTANCE_REGION> \
  --add-cloudsql-instances <PROJECT_ID>:<INSTANCE_REGION>:<INSTANCE_NAME>
```

The Cloud Function command above can now be re-run with a successful deployment.

## Running Integration Tests

The integration tests depend on a Unix socket and a TCP listener provided by the
Cloud SQL Auth Proxy. To run the tests, you will need to start two instances of
the Cloud SQL Auth Proxy, one for a TCP connection and one for a Unix socket
connection.

```
cloud-sql-proxy <PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME> --port=3306
cloud-sql-proxy <PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME> \
    --unix-socket=/cloudsql
```

To run integration tests, use the following command, setting environment
variables to the correct values:

```
GOLANG_SAMPLES_E2E_TEST="yes" \
  MYSQL_USER=some-user \
  MYSQL_PASSWORD=some-pass
  MYSQL_DATABASE=some-db \
  MYSQL_PORT=3307 \
  MYSQL_HOST=127.0.0.1
  MYSQL_UNIX_SOCKET='/cloudsql/<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>'
  MYSQL_INSTANCE='<PROJECT-ID>:<INSTANCE-REGION>:<INSTANCE-NAME>' \
  go test -v
```
