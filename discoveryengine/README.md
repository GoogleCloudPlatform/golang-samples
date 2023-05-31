# Gen App builder / Enterprise Search - Go Sample

Go sample for Gen App Builder / Enterprise Search, equivalent of the [python
example](https://cloud.google.com/generative-ai-app-builder/docs/libraries#client-libraries-usage-python)


### Install the client library

```
go get cloud.google.com/go/discoveryengine
```

### Use the client library

Use the code provided with your PROJECT_ID and your search engine ID and a query

```
export PROJECT_ID=$(gcloud config get project)
export SEARCH_ENGINE_ID=partscatalog_1682542047684

go run *.go --project $PROJECT_ID --searchengine $SEARCH_ENGINE_ID --query "fuel tank"
```