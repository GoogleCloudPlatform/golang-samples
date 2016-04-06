## Running tests

Prerequisites:

* Install the [Google Cloud SDK](https://cloud.google.com/sdk/).
* Install the `preview` and `app` commands. You can do this via:

        $ gcloud --quiet help preview app

* Install aedeploy:

        $ go get google.golang.org/appengine/cmd/aedeploy

Before running tests:

    $ gcloud auth login

Running without Docker:

    $ GOLANG_SAMPLES_E2E_TEST=1 go test -v
