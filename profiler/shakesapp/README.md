# Profiler Shakespeare Application in Go

This application starts a server which will return the number of times a word or
phrase appears in the works of Shakespeare, and sends queries to that server to
create load.

Profiler is enabled for this application, and can be used to identify
opportunities to optimize the server code.

## Running the application

This application should be run with Go 1.14. Once that version of Go is
installed:

*   Get the application code:

    ```sh
    git clone https://github.com/GoogleCloudPlatform/golang-samples.git
    cd golang-samples/getting-started/profiler/shakesapp
    ```

*   Install dependencies:

    ```sh
    go get .
    ```

*   Run the application:

    ```sh
    go run . -project_id="your-project-id" -version="application-version"
    ```
