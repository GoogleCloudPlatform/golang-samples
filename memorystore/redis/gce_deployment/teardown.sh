#!/bin/bash
# [START teardown_sh]
gcloud compute instances delete my-instance

gcloud compute firewall-rules delete allow-http-server-8080
# [END teardown_sh]
