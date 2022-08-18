# Docker Containers for golang-samples

We run tests for all supported versions of Go in Docker containers. You can run
tests in these containers to make sure your setup is the same as the test
environment.

See [`CONTRIBUTING.md`](/CONTRIBUTING.md),
[`system_tests.sh`](/testing/kokoro/system_tests.sh), and the `.cfg` files in
[`/testing/kokoro`](/testing/kokoro).

When new Go versions are released, we should build and push new versions of
these containers.

Edit `testing/docker/cloudbuild.yaml` to use the new version, then locally in the
`testing/docker` directory, run:

```
gcloud builds submit \
    --timeout 1h \  
    --project=golang-samples-tests \
    --config=cloudbuild.yaml .
```