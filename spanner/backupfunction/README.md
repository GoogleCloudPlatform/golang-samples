# Scheduled Backups System (SBS)

This shows an example how to use Cloud Scheduler and Cloud Function to configure
a schedule for creating Cloud Spanner backups.

## Run SpannerCreateBackup locally

Start a terminal for running a local server:

```bash
cd cmd
go run main.go
```

Start another terminal for calling the function:

```bash
DATA=$(printf '{"database":"projects/[PROJECT_ID]/instances/[INSTANCE_ID]/databases/[DATABASE_ID]", "expire": "6h"}'|base64) && curl --data '{"data":"'$DATA'"}' localhost:8080
```

## Run SpannerCreateBackup in Cloud

Create a pub/sub topic: 

```bash
gcloud pubsub topics create cloud-spanner-scheduled-backups
```

Deploy the `SpannerCreateBackup` function that subscribes the above topic: 

```bash
gcloud functions deploy SpannerCreateBackup --trigger-topic cloud-spanner-scheduled-backups --runtime go113
```

Call the `SpannerCreateBackup` function from command-line:

```bash
DATA=$(printf '{"database":"projects/[PROJECT_ID]/instances/[INSTANCE_ID]/databases/[DATABASE_ID]", "expire": "6h"}'|base64) && gcloud functions call SpannerCreateBackup --data '{"data":"'$DATA'"}'
```
