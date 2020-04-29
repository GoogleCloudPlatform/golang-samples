# Scheduled Backups System (SBS)

This shows an example how to use Cloud Scheduler and Cloud Function to configure
a schedule for creating Cloud Spanner backups.

## Run SpannerCreateBackup locally

Start a terminal for running a local server:

```bash
go run cmd/local_func_server/main.go
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

## Deploy scheduled jobs to Cloud Scheduler

Note: To use Cloud Scheduler, we must [create an App Engine app](https://cloud.google.com/scheduler/docs#supported_regions).

Make a copy of `schedule-template.yaml`, name it as `schedule.config.yaml` and
replace `PROJECT_ID`, `INSTANCE_ID`, `DATABASE_ID` with your configurations.

Deploy scheduled jobs for creating backups:

```bash
go run cmd/scheduler/main.go -config schedule.config.yaml
```
