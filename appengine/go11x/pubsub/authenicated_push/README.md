#  Pub/Sub sample for Google App Engine Standard Environment

[![Open in Cloud Shell][shell_img]][shell_link]

[shell_img]: http://gstatic.com/cloudssh/images/open-btn.png
[shell_link]: https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=appengine/go11x/pubsub/authenicated_push/README.md

This demonstrates how to receive messages using [Cloud Pub/Sub](https://cloud.google.com/pubsub) on
[Google App Engine Standard Environment](https://cloud.google.com/appengine/docs/standard/).

## Setup

Before you can run or deploy the sample, you will need to do the following:

1. Enable the Cloud Pub/Sub API in the 
[Google Developers Console](https://console.developers.google.com/project/_/apiui/apiview/pubsub/overview).

2. Allow Cloud Pub/Sub to create authentication tokens in your project.

        $ gcloud projects add-iam-policy-binding [your-project-id] \
            --member=serviceAccount:service-[your-project-number]@gcp-sa-pubsub.iam.gserviceaccount.com \
            --role=roles/iam.serviceAccountTokenCreator

3. Create a topic and subscription. The `--push-auth-service-account` flag activates the Pub/Sub push functionality for
Authentication and Authorization. Pub/Sub messages pushed to your endpoint will carry the identity of this service
account. You may use an existing service account or create a new one. The `--push-auth-token-audience` flag is optional;
if set, remember to modify the audience field check in `main.go`.

```
$ gcloud pubsub topics create [your-topic-name]
$ gcloud pubsub subscriptions create [your-subscription-name] \
    --topic=[your-topic-name] \
    --push-endpoint= https://[your-app-id].appspot.com/pubsub/message/receive?token=[your-verification-token] \
    --ack-deadline=30 \
    --push-auth-service-account=[your-service-account] \
    --push-auth-token-audience=http://example.com

1. Update the environment variables in ``app.yaml``.

## Running locally

When running locally, you can use the [Google Cloud SDK](https://cloud.google.com/sdk) to provide authentication to use
Google Cloud APIs:

```
gcloud init
```

Then set environment variables before starting your application:

```
$ export PUBSUB_VERIFICATION_TOKEN=[your-verification-token]
$ go run main.go
```

### Simulating push notifications

You can simulate a push message by making an HTTP request to the local push
notification endpoint.

```
$ curl --request POST \
--header 'Content-Type: application/json' \
--data '{"message": {"data": "This is an example message"}}' \
http://localhost:8080/pubsub/message/receive?token=[your-verification-token]
```

Response:

```
Missing Authorization header
```

The simulated push request fails because it does not have a Cloud Pub/Sub-generated JWT in the "Authorization" header.

## Running on App Engine

In the current directory, deploy using `gcloud`:

```
$ gcloud app deploy
```

Send a message with `gcloud`:

```
$ gcloud pubsub topics publish [your-topic-name] --message "This is a test"
```

View last 10 received messages from App Engine by visiting: 

```
https://[your-app-id].appspot.com/pubsub/message/list
```
