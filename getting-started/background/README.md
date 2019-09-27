Background Processing
---------------------

This directory contains an example of doing background processing with App
Engine, Cloud Pub/Sub, Cloud Functions, and Firestore.

Deploy commands:

```
$ GO111MODULE=on gcloud app deploy
$ gcloud functions deploy --runtime=go111 --trigger-topic=translate Translate --set-env-vars GOOGLE_CLOUD_PROJECT=my-project
```