#! /bin/bash

# Copyright 2015 Google Inc. All rights reserved.
# Use of this source code is governed by the Apache 2.0
# license that can be found in the LICENSE file.

set -x

ZONE=us-central1-f
gcloud config set compute/zone $ZONE

GROUP=frontend-group
TEMPLATE=$GROUP-tmpl
SERVICE=frontend-web-service

gcloud compute instance-groups managed stop-autoscaling $GROUP --zone $ZONE

gcloud compute forwarding-rules delete $SERVICE-http-rule --global 

gcloud compute target-http-proxies delete $SERVICE-proxy 

gcloud compute url-maps delete $SERVICE-map 

gcloud compute backend-services delete $SERVICE 

gcloud compute http-health-checks delete ah-health-check

gcloud compute instance-groups managed delete $GROUP  

gcloud compute instance-templates delete $TEMPLATE 

gcloud compute firewall-rules delete default-allow-http-8080
