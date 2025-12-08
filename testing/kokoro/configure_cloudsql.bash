#!/bin/bash

# Copyright 2021 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Command gimmeproj provides access to a pool of projects.
#
# The metadata about the project pool is stored in Cloud Datastore in a meta-project.
# Projects are leased for a certain duration, and automatically returned to the pool when the lease expires.
# Projects should be returned before the lease expires.

set -ex

# Download and prepare Cloud SQL Proxy
wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64
mv cloud_sql_proxy.linux.amd64 /cloud_sql_proxy
chmod +x /cloud_sql_proxy
mkdir /cloudsql && chmod 0777 /cloudsql

/cloud_sql_proxy -instances="${MYSQL_INSTANCE}"=tcp:3306,${MYSQL_INSTANCE} -dir /cloudsql &
/cloud_sql_proxy -instances="${POSTGRES_INSTANCE}"=tcp:5432,${POSTGRES_INSTANCE} -dir /cloudsql &
/cloud_sql_proxy -instances="${SQLSERVER_INSTANCE}"=tcp:1433 &

# Give proxies a second to connect before moving on. If future restructuring of Go's Kokoro
# test suite ever means this isn't enough time, reordering or increasing the sleep is reasonable.
sleep 5
