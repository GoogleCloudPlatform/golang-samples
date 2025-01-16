# Use the offical Go image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.21 as builder

# Copy local code to the container image.
RUN mkdir /app
WORKDIR /app
COPY . .

# Build the handle_task command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep"go .)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o handle_task

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3.14

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/handle_task /handle_task

# Run the web service on container startup.
CMD ["/handle_task"]
