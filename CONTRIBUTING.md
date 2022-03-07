# Contributing

1. Sign one of the contributor license agreements below.
1. [Install Go](https://golang.org/doc/install).
1. Clone the repo:

   `git clone https://github.com/GoogleCloudPlatform/golang-samples.git`

1. Change into the checked out source:

   `cd golang-samples`

1. Fork the repo.
1. Set your fork as a remote:

   `git remote add fork https://github.com/GITHUB_USERNAME/golang-samples.git`

1. Make changes (see [Formatting](#formatting) and [Style](#style)) and commit
   to your fork. Initial commit messages should follow the
   [Conventional Commits](https://www.conventionalcommits.org/) style (e.g.
   `feat(functions): add gophers codelab`).
1. Send a pull request with your changes.
1. A maintainer will review the pull request and make comments. Prefer adding
   additional commits over amending and force-pushing since it can be difficult
   to follow code reviews when the commit history changes.

   Commits will be squashed when they're merged.

## Tested Go versions

We test using the oldest and newest supported Go versions. We do not test the
intermediate versions. See [testing/kokoro](testing/kokoro).

## Secrets

Sensitive secrets are provided to the testing environment via a secret of
environment variables that Kokoro sources for each execution of the test suite.
If you need to add any additional secrets, execute the following command to pull
down the secrets locally:

```sh
$ ./testing/kokoro/pull-secrets.sh
```

This will write a gitignored file to your local checkout at `./testing/kokoro/test-env.sh`.
Add any environment variables as necessary to this file (and don't forget to
mimick their empty versions in `./testing/kokoro/test-env.tmpl.sh` for local
testing), and then run the following command to upload the new secret values:

```
$ ./testing/kokoro/store-secrets.sh
```

## Dependencies

When you add a new dependency, it must be added to `go.mod` in the directory you
are working in (if one is present) or in the repo root. After adding to the root
`go.mod`, use the following to update all `go.sum` files appropriately:

```
find . -name go.mod -execdir go mod tidy \;
```

# Formatting

All code must be formatted with `gofmt` (with the latest Go version) and pass
`go vet`.

# Style

Please read and follow https://github.com/golang/go/wiki/CodeReviewComments for
all Go code in this repo.

The following style guidelines are specific to writing Go samples.

Canonical samples:

- Veneer client library with complex request: [`inspect_string.go`](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/dlp/snippets/inspect/inspect_string.go)
- Apiary client with normal request: [`dicom_store_create.go`](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/healthcare/dicom_store_create.go)
- Apiary client with complex request: [`fhir_resource_create.go`](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/healthcare/fhir_resource_create.go)
- Apiary client with file I/O: [`dicomweb_instance_store.go`](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/healthcare/dicomweb_instance_store.go)

## One file per sample

Each sample should be in its own file so the [imports used by the sample can
be included in the region tag](#include-imports-in-region-tags).

## Sample package name, file name, and directory

The top level directory should be the product the sample is for (e.g.
`functions` or `dlp`).

Sub-directories can be used to keep different groups of samples for the product
separate.

The package name should match the directory name, unless it's a quickstart.
[Quickstarts use `package main`](#only-quickstarts-have-package-main); the
default binary name is the name of the directory. See
https://golang.org/doc/effective_go.html#names.

Files should be named after the sample in them (e.g. `hello.go`). No need to
include the product name or "sample" in the filename.

If there are many samples to write in the same directory, use filename prefixes
to group the files acting on similar types (for example, when writing
create/update/delete type samples).

Hosting platform samples may require a different directory and file structure.
When possible, follow the pattern of existing samples for that product.

## Include imports and flags in region tags

The sample region (e.g. `[START foo]` and `[END foo]`) should include the import
block.

```go
// Package hello contains Hello samples.
package hello

// [START hello]
import "fmt"

func hello(w io.Writer) {
	fmt.Fprintln(w, "Hello, World")
}

// [END hello]
```

For quickstarts, the region should include the package declaration as well as any [flags](#function-arguments-for-quickstarts).

For snippets, the region should _not_ include the package declaration.

Also see [Imports](#imports).

## Print to an `io.Writer` for snippets

(Note: this doesn't apply to quickstarts)

Do not print to `stdout` or `stderr`. Pass `w io.Writer` as the first argument
to the sample function and print to it with `fmt.Fprintf(w, ...)`.

This pattern matches `http.Handler`s, which print to an `http.ResponseWriter`, normally named
`w`.

```go
func hello(w io.Writer) {
	fmt.Fprintln(w, "Hello, World.")
}
```

The output can be verified during testing using a buffer.

[inspect_test.go](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/dlp/snippets/inspect/inspect_test.go)

```go
func TestInspectString(t *testing.T) {
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	err := inspectString(buf, tc.ProjectID, "I'm Gary and my email is gary@example.com")
	if err != nil {
		t.Errorf("TestInspectFile: %v", err)
	}

	got := buf.String()
	if want := "Info type: EMAIL_ADDRESS"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

## Google Cloud Project ID

Quickstarts should use an example project ID or add a project ID flag.

If a project ID is needed, snippets should have a `projectID string` argument.

## Only quickstarts have `package main`

Sample code should not include a runnable binary. Binaries should only be
included for quickstarts (which should all be `package main` with the example
code in `func main`).

Quickstarts need to be in a separate directories from snippets because they need
to be in different packages.

## Declare a `context.Context` as needed

Don't pass a `context.Context` as an argument. New Go developers may not
understand where the `ctx` comes from.

```diff
- func hello(ctx context.Context, w io.Writer) { ... }
+ func hello(w io.Writer) {
+ 	ctx := context.Background()
+ 	// ...
+ }
```

## Function arguments for snippets

There should be as few function arguments as possible. An `io.Writer` and
project ID are the most common. If you need additional arguments (for example,
the ID of a resource to get or delete), there should be an example value in the
body of the sample function.

```go
// delete deletes the resource identified by name.
func delete(w io.Writer, name string) error {
	// name := "/projects/my-project/resources/my-resource"
	ctx := context.Background()
	client, err := foo.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("foo.NewClient: %v", err)
	}
	if err := client.Delete(name); err != nil {
		return fmt.Errorf("Delete: %v", err)
	}
	return nil
}
```

## Function arguments for quickstarts

Since [quickstarts use `package main`](#only-quickstarts-have-package-main), we use the `flag` package for
passing parameters into a quickstart, and use `testutil.BuildMain` to build and test your quickstart.

In your quickstart:

```go
func main() {
	projectID := flag.String("project_id", "", "Cloud Project ID")
	resourceName := flag.String("resourceName", "", "Name of resource")
	flag.Parse()

	fmt.Printf("projectID: %s, resource_name: %s", *projectID, *resourceName)
	// ...
}
```

In your quickstart test:

```go
func TestQuickstart(t *testing.T) {
	tc := testutil.SystemTest(t)
	m := testutil.BuildMain(t)

	if !m.Built() {
		t.Fatalf("failed to build app")
	}

	testResourceName := "my-resource-name"

	stdOut, stdErr, err := m.Run(nil, 10*time.Minute,
		"--project_id", tc.ProjectID,
		"--resource_name", testResourceName,
	)

	if err != nil {
		t.Errorf("stdout: %v", string(stdOut))
		t.Errorf("stderr: %v", string(stdErr))
		t.Errorf("execution failed: %v", err)
	}

	if got := string(stdOut); !strings.Contains(got, testResourceName) {
		t.Errorf("got %q, want to contain %q", got, testResourceName)
	}
}
```

## Don't export sample functions

Sample functions should not be
[exported](https://golang.org/ref/spec#Exported_identifiers). Users should not
be depending directly on this sample code. So, the function name should start
with a lower case letter.

## Prefer inline proto declarations

Where possible, prefer a single declaration when initializing a proto value.

Request values should usually be named `req` and be declared on their own so the
API call (which uses `req`) is easier to understand.

```diff
- myRequest := &pb.Request{}
- myRequest.Parent = projectID
+ req := &pb.Request{
+ 	Parent: projectID,
+ }
```

## Line length

Generally, Go code does not have a strict line length limit. See
[the Code Review Comments wiki](https://github.com/golang/go/wiki/CodeReviewComments#line-length).
However, sample code is embedded on cloud.google.com and very long lines can be
difficult to read in the embedded code viewer.

Keep lines under around 100 characters, keeping in mind the general rules in the
[wiki](https://github.com/golang/go/wiki/CodeReviewComments#line-length).

## Initialize clients and services in every sample

Don't initialize one client for the entire set of samples and pass it as an
argument. Each sample should initialize its own client/service.

```diff
- func hello(client foo.Client, w io.Writer) { ... }
+ func hello(w io.Writer) {
+	ctx := context.Background()
+	client, err := foo.NewClient(ctx)
+	// ...
+ }
```

## Return errors

If the sample can run into errors, return the errors with additional context.
Don't call `log.Fatal` or friends.

`log.Fatal` is difficult to test because it will stop the entire test suite.

Use `fmt.Errorf` to add information when returning errors. Usually, the name of
the `package.Function` or just `Function` that returned the error is enough. It
may also help to include any arguments that were passed to the function.

Prefer inline error declaration when they aren't needed outside the `if`
statement.

```diff
// delete deletes the resource identified by resourceID.
func delete(w io.Writer, resourceID string) error {
	// resource := fmt.Sprintf("/projects/my-project/resources/%s", resourceID)
	ctx := context.Background()
	client, err := foo.NewClient(ctx)
-	if err != nil {
-		log.Fatal(err)
-	}
-	err := client.Delete(resourceID)
-	if err != nil {
-		log.Fatal(err)
-	}
+	if err != nil {
+		return fmt.Errorf("foo.NewClient: %v", err)
+ 	}
+	if err := client.Delete(resourceID); err != nil {
+		return fmt.Errorf("Delete: %v", err)
+	}
	return nil
}
```

## Imports

Imports should be added and sorted by
[`goimports`](https://godoc.org/golang.org/x/tools/cmd/goimports). There should
be at least two groups, separated by a newline:

- Standard library
- Everything else

```go
import (
	"context"
	"fmt"
	"log"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"golang.org/x/exp/rand"
)
```

## Comment functions and packages

One file in the sample package should have a package comment. The comment is
shown as a description on
https://godoc.org/github.com/GoogleCloudPlatform/golang-samples. If there are
many files/samples in the package, it's common to create a `doc.go` file that
only has the package comment. The comment should start with
`Package packagename`.

Functions should have comments starting with the name of the function (with the
same capitalization, even if it's lower case).

```go
// Package foo contains samples for Foo.
package foo

// hello prints "Hello, World."
func hello(w io.Writer) { ... }
```

See
[Sample package name, file name, and directory](#sample-package-name-file-name-and-directory)
and https://golang.org/doc/effective_go.html#commentary.

## Common identifiers

- **`ctx`**: All `context.Context` values unless the original can't be shadowed.
- **`name`**: Fully-qualified resource names (e.g.
  `/projects/my-project/resource/my-resource`).
- **`parent`**: Partially-qualified resource names (e.g.
  `/projects/my-project`)
- **`req`**: Request value to send.
- **`projectID`**: Google Cloud project ID.

Names should always be camelCase, even if it's a constant. Initialisms/acronyms
should have consistent case (e.g. `createFHIRStore` and `fhirStoreID`). See
https://golang.org/doc/effective_go.html#names.

See [Don't export sample functions](#dont-export-sample-functions).

## Reflection

Do not use the `reflect` package.

## Go version in go.mod files

The Go version in `go.mod` files is the minimum version of Go supported by the
module. Generally, this should be the minimum version supported & tested by the
repo. There are some cases where we need a higher minimum version.

Do not update the minimum version unless required.

## Use `testutil` for tests

All tests should use `testutil.SystemTest` or variants. `testutil` checks the
`GOLANG_SAMPLES_PROJECT_ID` environment variable exists, and skips the test if
not.

If the test takes longer than ~2 minutes, use `testutil.EndToEndTest`.

If you can't use `testutil` for some reason, be sure to skip tests if
`GOLANG_SAMPLES_PROJECT_ID` is not set. This makes sure tests pass when someone
clones the repo and runs tests.

See [Print to an `io.Writer`](#print-to-an-iowriter) for a full test example.

See [Testing](#testing).

# Testing

Tests are required for all samples. When writing a pull request, be sure to
write and run the tests in any modified directories.

See [Use `testutil` for tests](#use-testutil-for-tests) and
[Print to an `io.Writer`](#print-to-an-iowriter).

## Creating resources for tests

When creating resources for tests, avoid using UUIDs. Instead, prefer
resource names that incorporate aspects of your test, such as `tc.ProjectID + -golang-test-mypai-mysnippet`.

The intention for this is that test can be run in parallel with different
versions of Go without race conditions. You may need a different or more
specific identifier, depending on the sample and test.

## Running system tests

1. To run the system test yourself, you need a Google Cloud Project and a service account. During the creation of the service account, you should download the JSON credential file.
1. Ensure the IAM roles needed by the tests are added to the service account.
1. Ensure you are logged in using `gcloud auth login`.
1. Set the `GOLANG_SAMPLES_PROJECT_ID` environment variable to a suitable test project.
1. Set the `GOOGLE_APPLICATION_CREDENTIALS` environment variable to the credential file path. Tests are authenticated using [Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials).
1. Install the test dependencies:

   `go get -t -d github.com/GoogleCloudPlatform/golang-samples/...`

1. Run the tests:

   `go test github.com/GoogleCloudPlatform/golang-samples/...`

Note: You may want to `cd` to the directory you're modifying and run
`go test -v ./...` to avoid running every test in the repo.

# Contributor License Agreements

Before we can accept your pull requests you'll need to sign a Contributor
License Agreement (CLA):

- **If you are an individual writing original source code** and **you own the
  intellectual property**, then you'll need to sign an [individual CLA][indvcla].
- **If you work for a company that wants to allow you to contribute your work**,
  then you'll need to sign a [corporate CLA][corpcla].

You can sign these electronically (just scroll to the bottom). After that,
we'll be able to accept your pull requests.

[gcloudcli]: https://developers.google.com/cloud/sdk/gcloud/
[indvcla]: https://developers.google.com/open-source/cla/individual
[corpcla]: https://developers.google.com/open-source/cla/corporate
