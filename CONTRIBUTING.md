# Contributing

* [Ways to Contribute](#ways-to-contribute)
* [Getting Ready to Contribute](#getting-ready-to-contribute)
    * [Contributor License Agreements](#contributor-license-agreements)
    * [Development Environment Setup](#development-environment-setup)
* [Code Style](#code-style)
* [Testing](#testing)
* [Pull Request Lifecycle](#pull-request-lifecycle)

## Ways To Contribute

Thank you for your interest in `golang-samples`!

This repository hosts code samples linked in cloud.google.com documentation.
Because samples will be accompanied by separate documentation, we do not
typically welcome unsolicited new samples. However, if you feel a specific
sample is missing or incorrect, please [file an issue](https://github.com/GoogleCloudPlatform/golang-samples/issues/new/choose), and we can discuss the
available options.

If you'd like to contribute to existing samples, have a look at our [issues
list](https://github.com/GoogleCloudPlatform/golang-samples/issues) to see where we could use your help. Leave a comment on the issue to let
others know you are interested.

## Getting ready to contribute

### Contributor License Agreements

Before we can accept your contributions, you'll need to sign a Contributor
License Agreement (CLA):

- **If you are an individual writing original source code** and **you own the
  intellectual property**, then you'll need to sign an [individual CLA][indvcla].
- **If you work for a company that wants to allow you to contribute your work**,
  then you'll need to sign a [corporate CLA][corpcla].

You can sign these electronically (just scroll to the bottom). After that,
we'll be able to accept your pull requests.

### Development environment setup

1. [Install Go](https://golang.org/doc/install).

1. To contribute your changes, you'll most likely need to fork the repository.
   This can be done from the "Fork" menu in the Github UI, or with the [Github
   CLI](http://cli.github.com) command: `gh repo fork
   GoogleCloudPlatform/golang-samples`.

1. Clone the repo. Replace `${GITHUB_OWNER}` with your own github user name to
   clone your fork.

   `git clone https://github.com/${GITHUB_OWNER}/golang-samples.git`

1. Change into the checked out source:

   `cd golang-samples`

1. You are now ready to make your changes. See [Pull Request
   Lifecycle](#pull-request-lifecycle) to learn how to send your changes for
   review.

## Code Style

All code must be formatted with `gofmt` (with the latest Go version) and pass
`go vet`. To run these tools on samples in the `iam` directory, you would run `make lint dir=iam` from the root of the repository.

The [Google Cloud Samples Style Guide][style-guide] is considered the primary
guidelines for all Google Cloud samples. This section details some additional,
Go-specific rules that will be merged into the Samples Style Guide in the near
future.

[style-guide]: https://googlecloudplatform.github.io/samples-style-guide/

Please read and follow https://github.com/golang/go/wiki/CodeReviewComments for
all Go code in this repo.

The following style guidelines are specific to writing Go samples.

### Google Cloud Project ID

If a project ID is needed, snippets should have a `projectID string` argument.

Tests that require a Project ID should use [`testutil`](https://pkg.go.dev/github.com/GoogleCloudPlatform/golang-samples/internal/testutil) helper functions, or
consult the `GOLANG_SAMPLES_PROJECT_ID` environment variable.

### Sample package name, file name, and directory

The top level directory should be the product the sample is for (e.g.
`functions` or `dlp`).

Sub-directories can be used to keep different groups of samples for the product
separate.

The package name should match the directory name, as is standard Go practice.

Files should be named after the sample in them (e.g. `hello.go`). No need to
include the product name or "sample" in the filename.

If there are many samples to write in the same directory, use filename prefixes
to group the files acting on similar types (for example, when writing
create/update/delete type samples).

Hosting platform samples may require a different directory and file structure.
When possible, follow the pattern of existing samples for that product.

For snippets, the region should _not_ include the package declaration.

### Print to an `io.Writer` for snippets

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

### Declare a `context.Context` as needed

Don't pass a `context.Context` as an argument. New Go developers may not
understand where the `ctx` comes from.

```diff
- func hello(ctx context.Context, w io.Writer) { ... }
+ func hello(w io.Writer) {
+ 	ctx := context.Background()
+ 	// ...
+ }
```

### Function arguments for snippets

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

### Don't export sample functions

Sample functions should not be
[exported](https://golang.org/ref/spec#Exported_identifiers). Users should not
be depending directly on this sample code. So, the function name should start
with a lower case letter.

### Prefer inline proto declarations

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

### Line length

Generally, Go code does not have a strict line length limit. See
[the Code Review Comments wiki](https://github.com/golang/go/wiki/CodeReviewComments#line-length).
However, sample code is embedded on cloud.google.com and very long lines can be
difficult to read in the embedded code viewer.

Keep lines under around 100 characters, keeping in mind the general rules in the
[wiki](https://github.com/golang/go/wiki/CodeReviewComments#line-length).

### Return errors

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
+		return fmt.Errorf("foo.NewClient: %w", err)
+ 	}
+	if err := client.Delete(resourceID); err != nil {
+		return fmt.Errorf("Delete: %w", err)
+	}
	return nil
}
```

## Go version in go.mod files

The Go version in `go.mod` files is the minimum version of Go supported by the
module. Generally, this should be the minimum version supported & tested by the
repo. There are some cases where we need a higher minimum version.

Do not update the minimum version unless required.


## Testing

Tests are required for all samples. When writing a pull request, be sure to
write and run the tests in any modified directories.

### Running tests

To run the system test yourself, you will need a Google Cloud Project, and valid
authentication.

1. Ensure you are logged in using `gcloud auth login --update-adc`.
    * the `--update-adc` flag refreshes [Application Default
      Credentials](https://developers.google.com/identity/protocols/application-default-credentials).
1. Set the `GOLANG_SAMPLES_PROJECT_ID` environment variable to a suitable test project.
1. To run all tests in a directory, run `make test dir=relative/dir`

### Use `testutil` for tests

All tests should use `testutil.SystemTest` or variants. `testutil` checks the
`GOLANG_SAMPLES_PROJECT_ID` environment variable exists, and skips the test if
not.

If the test takes longer than ~2 minutes, use `testutil.EndToEndTest`.

If you can't use `testutil` for some reason, be sure to skip tests if
`GOLANG_SAMPLES_PROJECT_ID` is not set. This makes sure tests pass when someone
clones the repo and runs tests.

### Creating resources for tests

Tests are responsible for creating any resources they require, and destroying
them once testing is complete. Names of these
resources should be unique enough to avoid conflicts in the event of multiple
concurrent test runs - this typically means suffixing them with some thing
unique, like a timestamp.


## Pull Request lifecycle

1. Before creating a Pull Request, ensure that your code meets style guidelines,
   and the tests pass. From the root of the repository, run `make lint test
   dir=relative/dir`.
1. [Create a pull
request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request)
with your changes.
    * PR titles should follow [Conventional
      Commits](https://www.conventionalcommits.org/) style (e.g.
      `feat(functions): add gophers codelab`).
    * You may wish to [enable
      automerge](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/incorporating-changes-from-a-pull-request/automatically-merging-a-pull-request#enabling-auto-merge)
      on your PR, so it submits when all PR checks are passing (including
      review).

1. Within 2-5 days, a reviewer will review your PR. They may approve it, or
   request changes. When requesting changes, reviewers should self-assign the
   PR to ensure they are aware of any updates.

1. If additional changes are needed, push additional commits to your PR branch -
   this helps the reviewer know which parts of the PR have changed.  Commits
   will be squashed when merged.

   - Please follow up with changes promptly. If a PR is awaiting changes by the
     author for more than 10 days, maintainers may mark that PR as Draft. PRs
     that are inactive for more than 30 days may be closed.

[gcloudcli]: https://developers.google.com/cloud/sdk/gcloud/
[indvcla]: https://developers.google.com/open-source/cla/individual
[corpcla]: https://developers.google.com/open-source/cla/corporate
