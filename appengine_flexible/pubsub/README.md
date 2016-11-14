## Setup

Before you can run or deploy this sample, you will need to configure a Pub/Sub topic and subscription:

1. Enable the Cloud Pub/Sub API in the [Google Developers Console](https://console.developers.google.com/project/_/apiui/apiview/pubsub/overview).

2. Create a topic and subscription.

        $ gcloud alpha pubsub topics create [your-topic-name]
        $ gcloud alpha pubsub subscriptions create [your-subscription-name] \
            --topic [your-topic-name] \
            --push-endpoint \
                https://[your-app-id].appspot.com/pubsub/push \
            --ack-deadline 30

3. Update the environment variables in `app.yaml`.
