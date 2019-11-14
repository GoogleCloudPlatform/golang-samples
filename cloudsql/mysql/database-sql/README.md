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

1. Use the information noted in the previous steps:
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service/account/key.json
export INSTANCE_CONNECTION_NAME='<MY-PROJECT>:<INSTANCE-REGION>:<INSTANCE-NAME>'
export DB_USER='my-db-user'
export DB_PASS='my-db-pass'
export DB_NAME='my_db'
```
Note: Saving credentials in environment variables is convenient, but not secure - consider a more
secure solution such as [Cloud KMS](https://cloud.google.com/kms/) to help keep secrets safe.

## Running locally

To run this application locally, download and install the `cloud_sql_proxy` by
following the instructions
[here](https://cloud.google.com/sql/docs/mysql/sql-proxy#install). Once the
proxy has been downloaded, use the following commands start it running in the background for local testing on either Linux/Mac OS or Windows.

### launch the proxy on Linux/Mac OS

Linux environments use a Unix socket for the Cloud SQL proxy, so you'll need to create the `/cloudsql`
directory and give the user running the proxy the appropriate permissions:
```bash
sudo mkdir /cloudsql
sudo chown -R $USER /cloudsql
```

Once the `/cloudsql` directory is ready, use the following command to start the proxy in the
background:
```bash
./cloud_sql_proxy -dir=/cloudsql --instances=$CLOUD_SQL_CONNECTION_NAME --credential_file=$GOOGLE_APPLICATION_CREDENTIALS
```
Note: Make sure to run the command under a user with write access in the 
`/cloudsql` directory. This proxy will use this folder to create a unix socket
the application will use to connect to Cloud SQL. 

### launch the proxy on Windows

Windows environments use a TCP connection for the Cloud SQL proxy. Use this PowerShell command to launch the proxy:

```powershell
Start-Process -filepath "C:\<path to proxy exe>" -ArgumentList "-instances=<project-id>:<region>:<instance-name>=tcp:3306"
```

### testing the application

To test the application locally, follow these steps after the proxy is running:

* Install the MySQL driver: `go get github.com/go-sql-driver/mysql`
* Run the application: `go run cloudsql.go`
* Navigate to `http://127.0.0.1:8080` in a web browser to verify your application is running correctly.

## Google App Engine Standard

To run on GAE-Standard, create an App Engine project by following the setup for these 
[instructions](https://cloud.google.com/appengine/docs/standard/python3/quickstart#before-you-begin).

First, create an `app.yaml` with the correct values to pass the environment 
variables into the runtime. Your app.yaml file should look like this:

```yaml
runtime: go111
env_variables:
  INSTANCE_CONNECTION_NAME: <project-id>:<region>:<instance-name>
  DB_USER: root
  DB_PASS: <password for root user>
  DB_NAME: my_db
```

Next, the following command will deploy the application to your Google Cloud project:
```bash
gcloud app deploy
```
