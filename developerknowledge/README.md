# Google Developer Knowledge API Go Samples

This directory contains Go code samples demonstrating how to use the
[Google Developer Knowledge API](https://developers.google.com/knowledge) client library (`cloud.google.com/go/developerknowledge/apiv1`).

## Setup

1. Enable the Developer Knowledge API on your Google Cloud project:
   ```bash
   gcloud services enable developerknowledge.googleapis.com
   ```

2. Download module dependencies:
   ```bash
   go mod download
   ```

## Samples

* **[Search Document Chunks](search_document_chunks.go)**: Search public developer documentation chunks by query (`developerknowledge_search_document_chunks`).
* **[Get Document](get_document.go)**: Retrieve a single documentation page with full markdown content (`developerknowledge_get_document`).
* **[Batch Get Documents](batch_get_documents.go)**: Fetch multiple documentation pages in one call (`developerknowledge_batch_get_documents`).
* **[Answer Query](answer_query.go)**: Get a grounded, cited answer to a technical question (`developerknowledge_answer_query`).

## Running Tests

```bash
go test -v .
```
