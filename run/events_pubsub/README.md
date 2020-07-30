# Run – Events Pub/Sub

A simple HTTP server using the CloudEvents SDK.

## Commands

Run locally:

```sh
go run main.go
```

Run in container:

```sh
docker build . -t pubsub-event && docker run --rm -p 8080 --expose 8080 pubsub-event
```

Test locally:

```sh
go test
```