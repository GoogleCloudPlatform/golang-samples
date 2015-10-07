# Accessing BigQuery from App Engine

A Google Cloud Platform customer asked me today how to list all the BigQuery
projects that you own from a Google App Engine app.

If you don’t know what BigQuery or App Engine are this post is probably not
for you … yet! Instead you should have a look at the docs for [BigQuery][1],
and [App Engine][2], two of my favorite products of [Google Cloud Platform][0].

<div style="text-align:center">
[![Google Cloud Platform Logo](https://cloud.google.com/_static/images/new-gcp-logo.png)](http://cloud.google.com)
</div>

The solution for this is quite simple, but I think there’s enough
moving pieces than a blog post was required. Let’s assume that the list will
be displayed as part of some the handling to a request, something like this:

```go
func handle(w http.ResponseWriter, r *http.Request) {

    // create a new App Engine context from the request.
    c := appengine.NewContext(r)

    // obtain the list of project names.
    names, err := projects(c)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // print it to the output.
    fmt.Fprintln(w, strings.Join(names, "\n"))

}
```

Every time a new HTTP request comes the handler define above will be executed.
It will create a new App Engine context that will then be used by the `projects`
function to return a list with all the names of the BigQuery projects visible
by the App Engine app. Finally, the list will be printed to the output `w`.

In order to register the handler so it is executed on every HTTP request we add
an `init` func.

```go
func init() {

    // all requests are handled by handler.
    http.HandleFunc(“/”, handle)

}
```

Now for the interesting part, how do we implement the `projects` function? 

```go
func projects(c context.Context) ([]string, error) {

    // some awesome code ...

}
```

Let’s implement the body of that function in three parts:

## Create an authenticated HTTP client 

Given a [`context.Context`][3] named `c` we can create an authenticated client
using this code.

```go
// create a new HTTP client.
client := &http.Client{
    Transport: &oauth2.Transport{
        Source: google.AppEngineTokenSource(c,
            bigquery.BigqueryScope),
        Base: &urlfetch.Transport{Context: c},
    },
}
```

In Go, HTTP clients use transports to communicate, and transports use …
transports! It’s transports all the way! But what is a HTTP transport
exactly?

The Transport field in an HTTP client is of type [`RoundTripper`][4]:

```go
type RoundTripper interface {
        RoundTrip(*Request) (*Response, error)
}
```

A [`RoundTripper`][4] is responsible for generating a response given a request,
but it also has the capacity of changing the request and response being set.
This is very useful for monitoring, instrumenting, and also authentication.

The [`oauth2.Transport`][5] type, given a request, adds authentication headers and
forwards the request through its base transport.

```go
&oauth2.Transport{
    Source: google.AppEngineTokenSource(c, bigquery.BigqueryScope),
    Base: &urlfetch.Transport{Context: c},
}
```

The snippet above creates a new [`oauth2.Transport`][5] that authenticates
requests using the default service account for the App Engine project, and then
uses an HTTP transport provided by [`urlfetch`][6]  —  the way of accessing external
resources from App Engine apps.

## Create a BigQuery service and list all projects 

Creating a BigQuery service with [`google.golang.org/api/bigquery/v2`][7] is
quite simple. Just import the package and create a new service given an HTTP
client with [`bigquery.New`][8].

```go
bq, err := bigquery.New(client)
if err != nil {
    return nil, fmt.Errorf("create service: %v", err)
}
```

It’s important to check that the error is not nil, since the operation could
fail as any other operation over the network. Once we have the service, we can
list all the projects by following the documentation:

```go
list, err := bq.Projects.List().Do()
if err != nil {
    return nil, fmt.Errorf("list projects: %v", err)
}
```

## Create a list with the project names 

We got list, a [`bigquery.ProjectList`][9],  so we can iterate over all the
projects in `Projects` and append their `FriendlyName` to our list before
returning it.

```go
var names []string

for _, p := range list.Projects {
    names = append(names, p.FriendlyName)
}

return names, nil
```

I hope this will be helpful for many of you!

## Questions?

Feel free to file issues or contact me on [Twitter][10].

[0]: https://cloud.google.com
[1]: https://cloud.google.com/bigquery/what-is-bigquery
[2]: https://cloud.google.com/appengine/docs
[3]: https://godoc.org/golang.org/x/net/context#Context
[4]: https://golang.org/pkg/net/http/#RoundTripper
[5]: https://godoc.org/golang.org/x/oauth2#Transport
[6]: https://cloud.google.com/appengine/docs/go/urlfetch/
[7]: https://godoc.org/google.golang.org/api/bigquery/v2
[8]: https://godoc.org/google.golang.org/api/bigquery/v2#New
[9]: https://godoc.org/google.golang.org/api/bigquery/v2#ProjectList
[10]: http://twitter.com/francesc
