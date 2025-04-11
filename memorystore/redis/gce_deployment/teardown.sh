#!/bin/bash
# [START memorystore_teardown_sh]
gcloud compute instances delete my-instance

gcloud compute firewall-rules delete allow-http-server-8080
# [END memorystore_teardown_sh]
