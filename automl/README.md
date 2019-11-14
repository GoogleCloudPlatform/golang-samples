# AutoML Samples

<a href="https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/GoogleCloudPlatform/golang-samples&page=editor&open_in_editor=automl/README.md">
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
* [Get Model Evaluation](get_model_evaluation.go)
* [Delete Model](delete_model.go)
* [Deploy Model](deploy_model.go) - Not supported by Translation
* [Uneploy Model](undeploy_model.go) - Not supported by Translation

### Operation Management

* [List Operation Statuses](list_operation_status.go)
* [Get Operation Status](get_operation_status.go)

## AutoML Type Specific Samples

### Translation

* [Translate Create Dataset](translate_create_dataset.go)
* [Translate Create Model](translate_create_model.go)
* [Translate Predict](translate_predict.go)

### Natural Language Entity Extraction

* [Entity Extraction Create Dataset](language_entity_extraction_create_dataset.go)
* [Entity Extraction Create Model](language_entity_extraction_create_model.go)
* [Entity Extraction Predict](language_entity_extraction_predict.go)
* [Entity Extraction Batch Predict](language_batch_predict.go)

### Natural Language Sentiment Analysis

* [Sentiment Analysis Create Dataset](language_sentiment_analysis_create_dataset.go)
* [Sentiment Analysis Create Model](language_sentiment_analysis_create_model.go)
* [Sentiment Analysis Predict](language_sentiment_analysis_predict.go)

### Natural Language Text Classification

* [Text Classification Create Dataset](language_text_classification_create_dataset.go)
* [Text Classification Create Model](language_text_classification_create_model.go)
* [Text Classification Predict](language_text_classification_predict.go)

### Vision Classification

* [Classification Create Dataset](vision_classification_create_dataset.go)
* [Classification Create Model](vision_classification_create_model.go)
* [Classification Predict](vision_classification_predict.go)
* [Classification Batch Predict](vision_batch_predict.go)
* [Deploy Node Count](vision_classification_deploy_model_node_count.go)

### Vision Object Detection

* [Object Detection Create Dataset](vision_object_detection_create_dataset.go)
* [Object Detection Create Model](vision_object_detection_create_model.go)
* [Object Detection Predict](vision_object_detection_predict.go)
* [Object Detection Batch Predict](vision_batch_predict.go)
* [Deploy Node Count](vision_object_detection_deploy_model_node_count.go)
