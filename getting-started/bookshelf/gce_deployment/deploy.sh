#! /bin/bash

# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

if [ -z "$BOOKSHELF_DEPLOY_LOCATION" ]; then
  echo "Must set \$BOOKSHELF_DEPLOY_LOCATION. For example: BOOKSHELF_DEPLOY_TARGET=gs://my-bucket/bookshelf-VERSION.tar"
  exit 1
fi

ZONE=us-central1-f

GROUP=frontend-group
TEMPLATE=$GROUP-tmpl
MACHINE_TYPE=f1-micro
IMAGE=debian-8
STARTUP_SCRIPT=startup-script.sh
SCOPES="userinfo-email,cloud-platform"
TAGS=http-server

MIN_INSTANCES=1
MAX_INSTANCES=10
TARGET_UTILIZATION=0.6

SERVICE=frontend-web-service

#
# Instance group setup
#

# First we have to create an instance template.
# This template will be used by the instance group
# to create new instances.

# [START create_template]
gcloud compute instance-templates create $TEMPLATE \
  --image $IMAGE \
  --machine-type $MACHINE_TYPE \
  --scopes $SCOPES \
  --metadata-from-file startup-script=$STARTUP_SCRIPT \
  --metadata app-location=$BOOKSHELF_DEPLOY_LOCATION \
  --tags $TAGS
# [END create_template]

# Create the managed instance group.

# [START create_group]
gcloud compute instance-groups managed \
  create $GROUP \
  --base-instance-name $GROUP \
  --size $MIN_INSTANCES \
  --template $TEMPLATE \
  --zone $ZONE
# [END create_group]

# [START create_named_port]
gcloud compute instance-groups managed set-named-ports \
    $GROUP \
    --named-ports http:8080 \
    --zone $ZONE
# [END create_named_port]

#
# Load Balancer Setup
#

# A complete HTTP load balancer is structured as follows:
#
# 1) A global forwarding rule directs incoming requests to a target HTTP proxy.
# 2) The target HTTP proxy checks each request against a URL map to determine the
#    appropriate backend service for the request.
# 3) The backend service directs each request to an appropriate backend based on
#    serving capacity, zone, and instance health of its attached backends. The
#    health of each backend instance is verified using either a health check.
#
# We'll create these resources in reverse order:
# service, health check, backend service, url map, proxy.

# Create a health check
# The load balancer will use this check to keep track of which instances to send traffic to.
# Note that health checks will not cause the load balancer to shutdown any instances.

# [START create_health_check]
gcloud compute http-health-checks create ah-health-check \
  --request-path /_ah/health \
  --port 8080
# [END create_health_check]

# Create a backend service, associate it with the health check and instance group.
# The backend service serves as a target for load balancing.

# [START create_backend_service]
gcloud compute backend-services create $SERVICE \
  --http-health-checks ah-health-check
# [END create_backend_service]

# [START add_backend_service]
gcloud compute backend-services add-backend $SERVICE \
  --instance-group $GROUP \
  --zone $ZONE
# [END add_backend_service]

# Create a URL map and web Proxy. The URL map will send all requests to the
# backend service defined above.

# [START create_url_map]
gcloud compute url-maps create $SERVICE-map \
  --default-service $SERVICE
# [END create_url_map]

# [START create_http_proxy]
gcloud compute target-http-proxies create $SERVICE-proxy \
  --url-map $SERVICE-map
# [END create_http_proxy]

# Create a global forwarding rule to send all traffic to our proxy

# [START create_forwarding_rule]
gcloud compute forwarding-rules create $SERVICE-http-rule \
  --global \
  --target-http-proxy $SERVICE-proxy \
  --port-range 80
# [END create_forwarding_rule]

#
# Autoscaler configuration
#
# [START set_autoscaling]
gcloud compute instance-groups managed set-autoscaling \
  $GROUP \
  --max-num-replicas $MAX_INSTANCES \
  --target-load-balancing-utilization $TARGET_UTILIZATION \
  --zone $ZONE
# [END set_autoscaling]

# [START create_firewall]
gcloud compute firewall-rules create default-allow-http-8080 \
    --allow tcp:8080 \
    --source-ranges 0.0.0.0/0 \
    --target-tags http-server \
    --description "Allow port 8080 access to http-server"
# [END create_firewall]
