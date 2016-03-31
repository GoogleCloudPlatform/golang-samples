## Running tests

Prerequisites:

* Install the [Google Cloud SDK](gcloud).
* Install the `preview` and `app` commands. You can do this via:

        $ gcloud --quiet help preview app
* Install aedeploy:

        $ go get google.golang.org/appengine/cmd/aedeploy

Before running tests:

    $ gcloud config project set $PROJECT_ID
    $ gcloud auth login

Running via Docker:

    $ docker build -f integration_test_Dockerfile -t goaev2itest .
    $ docker run -v ~/.config/gcloud:/root/.config/gcloud goaev2itest

Running without Docker:

    $ go test

[gcloud]: https://cloud.google.com/sdk/
