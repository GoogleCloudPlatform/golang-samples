# Accessing BigQuery from App Engine

This code has been written to support this [blog post][1].

You can run this application locally or deploy it to App Engine app servers.
To do so you will need the `goapp` tool included in the [Go App Engine SDK][2].

- get all the dependencies for this code snippet

	goapp get .

- deploy to App Engine servers

	goapp deploy --application=[your-application-id] .

[1]: https://medium.com/@francesc/accessing-bigquery-from-app-engine-d01823de81ee
[2]: https://cloud.google.com/appengine/downloads?hl=en
