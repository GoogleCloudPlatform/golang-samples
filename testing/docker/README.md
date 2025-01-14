# Docker Containers for golang-samples

We run tests for all supported versions of Go in Docker containers. You can run
tests in these containers to make sure your setup is the same as the test
environment.

See [`CONTRIBUTING.md`](/CONTRIBUTING.md),
[`system_tests.sh`](/testing/kokoro/system_tests.sh), and the `.cfg` files in
[`/testing/kokoro`](/testing/kokoro).

When new Go versions are released, we should build and push new versions of
these containers.

Go language version and resulting image name are controlled by the cloud build
substitutions `_GO_VERSION` and `_IMAGE_NAME` respectively. The command below
will build Go 1.21 and push the resulting image to
`gcr.io/golang-samples-tests/go121`

```
gcloud builds submit . \
    --project=golang-samples-tests \
    --substitutions "_GO_VERSION=1.21,_IMAGE_NAME=go121"
```
