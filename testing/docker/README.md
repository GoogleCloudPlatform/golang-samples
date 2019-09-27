# Docker Containers for golang-samples

We run tests for all supported versions of Go in Docker containers. This
directory contains `Dockerfile`s for each version of Go. You can run tests in
these containers to make sure your setup is the same as the test environment.

See [`CONTRIBUTING.md`](/CONTRIBUTING.md),
[`system_tests.sh`](/testing/kokoro/system_tests.sh), and the `.cfg` files in
[`/testing/kokoro`](/testing/kokoro).

When new minor versions are released, we should build and push new versions of
these containers.

```
gcloud config set project golang-samples-tests
for v in go111 go112 go113;
  sudo docker build -f Dockerfile.$v --tag gcr.io/golang-samples-tests/$v .
  sudo docker push gcr.io/golang-samples-tests/$v
done
```