# Profiler Shakespeare application

This application starts a server which will return the number of times a word
or phrase appears in the works of Shakespeare, and sends queries to that server
to create load.

This application's server is intentionally non-optimal. Profiler is enabled for
this application, and can be used identify how to optimize the server.

## Running the application

*   Get the application code:

    ```sh
    git clone https://github.com/GoogleCloudPlatform/golang-samples.git
    cd golang-samples/profiler/shakesapp
    ```

*   Install dependencies:

    ```sh
    go get .
    ```

*   Run the application:

    ```sh
    go run . -version="application-version"
    ```
