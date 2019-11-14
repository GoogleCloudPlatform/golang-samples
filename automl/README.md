# AutoML Samples

<a href="https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/java-docs-samples&page=editor&open_in_editor=vision/beta/cloud-client/README.md">
<img alt="Open in Cloud Shell" src ="http://gstatic.com/cloudssh/images/open-btn.png"></a>

This directory contains samples for the [Google Cloud AutoML APIs](https://cloud.google.com/automl/) - [docs](https://cloud.google.com/automl/docs/)

We highly reccommend that you refer to the official documentation pages:

* AutoML Natural Language
  * [Classification](https://cloud.google.com/natural-language/automl/docs)
  * [Entity Extraction](https://cloud.google.com/natural-language/automl/entity-analysis/docs)
  * [Sentiment Analysis](https://cloud.google.com/natural-language/automl/sentiment/docs)
* [AutoML Translation](https://cloud.google.com/translate/automl/docs)
* AutoML Video Intelligence
  * [Classification](https://cloud.google.com/video-intelligence/automl/docs)
  * [Object Tracking](https://cloud.google.com/video-intelligence/automl/object-tracking/docs)
* AutoML Vision
  * [Classification](https://cloud.google.com/vision/automl/docs)
  * [Edge](https://cloud.google.com/vision/automl/docs/edge-quickstart)
  * [Object Detection](https://cloud.google.com/vision/automl/object-detection/docs)
* [AutoML Tables](https://cloud.google.com/automl-tables/docs)

This API is part of the larger collection of Cloud Machine Learning APIs.

These Go samples demonstrates how to access the Cloud AutoML API using the
[Google Cloud Client Library for Go](https://github.com/googleapis/google-cloud-go).

## Sample Types

There are two types of samples: Base and API Specific

The base samples make up a set of samples that have code that
is identical or nearly identical for each AutoML Type.
Meaning that for "Base" samples you can use them with any AutoML Type.
However, for API Specific samples, there will be a unique sample for each AutoML type.
See the below list for more info.

## Base Samples

### Dataset Management

* [Import Dataset](import_dataset.go)
* [List Datasets](list_datasets.go) - For each AutoML Type the `metadata` field inside the dataset is unique, therefore each AutoML Type will have a
small section of code to print out the `metadata` field.
* [Get Dataset](get_dataset.go) - For each AutoML Type the `metadata` field inside the dataset is unique, therefore each AutoML Type will have a
small section of code to print out the `metadata` field.
* [Export Dataset](export_dataset.go)
* [Delete Dataset](delete_dataset.go)

### Model Management

* [List Models](list_models.go)
* [List Model Evaluations](list_model_evaluations.go)
* [Get Model](get_model.go)
* [Get Model Evaluation](get_model_evaluation)
* [Delete Model](delete_model.go)

### Operation Management

* [List Operation Statuses](list_operation_status.go)
* [Get Operation Status](get_operation_status.go)

## AutoML Type Specific Samples

### Translation

* [Translate Create Dataset](translate_create_dataset.go)
* [Translate Create Model](translate_create_model.go)
* [Translate Predict](translate_predict.go)
