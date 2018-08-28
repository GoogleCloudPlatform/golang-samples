# Accessing BigQuery from App Engine

This code has been written to support this [blog post][1].

### Prerequisites

* [Go App Engine SDK][2].

### Deploy and run

You can run this application locally or deploy it to App Engine app servers.

1. Get the dependencies:

        goapp get .


1. Run the app locally:

        dev_appserver.py -A [your-project-id] .

1. Deploy the app:

        goapp deploy --application=[your-project-id] -version=[choose-a-version] .

[1]: https://medium.com/@francesc/accessing-bigquery-from-app-engine-d01823de81ee
[2]: https://cloud.google.com/appengine/downloads?hl=en
