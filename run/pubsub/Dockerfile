# Copyright 2019 Google LLC. All rights reserved.
# Use of this source code is governed by the Apache 2.0
# license that can be found in the LICENSE file.

# [START run_pubsub_dockerfile]

# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang as builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN CGO_ENABLED=0 GOOS=linux go build -v -o hellopubsub

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/hellopubsub .

# Run the web service on container startup.
CMD ["/hellopubsub"]

# [END run_pubsub_dockerfile]
