# Cloud Audit Logs: Service Account Key creation

This sample illustrates how to receive and process Cloud Audit Logs when a new
Service Account Key is created.

Note: Some organizations disable service account key creation. This sample will
not produce the expected outcome in these organizations, since the trigger
action cannot performed.

## Prerequisites

You should be familiar with Cloud Audit Logs and Service Account Keys.

You will need permission to create new service accounts, and service account
keys to follow these instructions.

## Procedure

### Provision a service account to invoke Cloud Run

  A service account is required when delivering Eventarc messages to Cloud Run.
  Create a new service account, and grant it the necessary roles to receive
  Eventarc messages, and invoke Cloud Run.
 
  - `gcloud iam service-accounts create iam-audit-account`
  - `gcloud projects add-iam-policy-binding ${project} --member=iam-audit-account@${project}.iam.gserviceaccount.com --role=roles/run.invoker`
  - `gcloud projects add-iam-policy-binding ${project} --member=iam-audit-account@${project}.iam.gserviceaccount.com --role=roles/eventarc.eventReceiver`

### Enable Cloud Audit Logs for Cloud IAM

Open the [Audit Logs page in cloud
console](https://console.cloud.google.com/iam-admin/audit), and enable "Data
Write" audit logs for "Identity and Access Management (IAM) API".

### Deploy the Cloud Run Service

This service processes the events from Eventarc when new service account keys
are created.

  `gcloud run deploy audit-iam-keys --source .`

### Create the Eventarc Trigger

```
gcloud eventarc triggers create audit-iam-trigger \
    --location=us-central1 \
    --event-filters=type=google.cloud.audit.log.v1.written \
    --event-filters=serviceName=iam.googleapis.com \
    --event-filters=methodName=google.iam.admin.v1.CreateServiceAccountKey \
    --destination-run-region=${REGION} \
    --destination-run-service=audit-iam-keys \
    --service-account=iam-audit-account@${project}.iam.gserviceaccount.com
```

### Wait ~10m

Eventarc Triggers may take up to 10 minutes to begin delivering events.

### Create a new Service Account Key

This triggers an event that our Cloud Run service will process.

```
gcloud iam service-accounts keys create ./sample.key --iam-account="iam-audit-account@${project}.iam.gserviceaccount.com"
```

### Watch for logs from the Cloud Run service

`gcloud alpha run services logs tail audit-iam-keys`

