#!/bin/bash

set -ex

# Download and prepare Cloud SQL Proxy
wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64
mv cloud_sql_proxy.linux.amd64 /cloud_sql_proxy
chmod +x /cloud_sql_proxy
mkdir /cloudsql && chmod 0777 /cloudsql

# TODO: These variable names need to be populated.
/cloud_sql_proxy -instances="${MYSQL_INSTANCE}"=tcp:3306,${MYSQL_INSTANCE} -dir /cloudsql
/cloud_sql_proxy -instances="${POSTGRES_INSTANCE}"=tcp:5432,${POSTGRES_INSTANCE} -dir /cloudsql
/cloud_sql_proxy -instances="${SQLSERVER_INSTANCE}"=tcp:1433