## Running tests

Prerequisites:

* Install the [Google Cloud SDK](https://cloud.google.com/sdk/).
* Install the `app` command. You can do this via:

        $ gcloud --quiet help app

Before running tests:

    $ gcloud auth login

Running without Docker:

    $ GOLANG_SAMPLES_E2E_TEST=1 go test -v
