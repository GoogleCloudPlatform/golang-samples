# Testing

CI tests are run with a Google-internal system. We have 10 test projects
(golang-samples-tests, golang-samples-tests-{2..10}). Each has its own service
account. Some APIs don't support service accounts in a different project than
the `GOOGLE_CLOUD_PROJECT`.

Here are commands to set up each project and service account, wrapped in a for
loop to make it easier.

```
for i in {2..10}; do
    gcloud config set project golang-samples-tests-$i    
    gcloud iam service-accounts create kokoro-golang-samples-tests-$i --display-name "Kokoro golang-samples-tests $i"
    gcloud iam service-accounts keys create kokoro-golang-samples-tests-$i.json --iam-account kokoro-golang-samples-tests-$i@golang-samples-tests-$i.iam.gserviceaccount.com
    gcloud projects add-iam-policy-binding golang-samples-tests-$i --member serviceAccount:kokoro-golang-samples-tests-$i@golang-samples-tests-$i.iam.gserviceaccount.com --role roles/owner
    gcloud projects add-iam-policy-binding golang-samples-tests-$i --member serviceAccount:kokoro-golang-samples-tests-$i@golang-samples-tests-$i.iam.gserviceaccount.com --role roles/cloudkms.cryptoKeyEncrypterDecrypter

    # Every service account should have access to the main project (for
    # resources like Spanner).
    gcloud projects add-iam-policy-binding golang-samples-tests --member serviceAccount:kokoro-golang-samples-tests-$i@golang-samples-tests-$i.iam.gserviceaccount.com --role roles/owner
done
```

Each project needs to be setup for vision/product_search product set tests.