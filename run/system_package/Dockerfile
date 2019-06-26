# Copyright 2019 Google LLC. All rights reserved.
# Use of this source code is governed by the Apache 2.0
# license that can be found in the LICENSE file.

# Use the offical Golang image to create a build artifact.
# https://hub.docker.com/_/golang
FROM golang as builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build the release artifact.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o graphviz-web

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
# https://hub.docker.com/_/alpine
FROM ubuntu

# [START run_system_package_ubuntu]
RUN apt-get update -y && apt-get install -y \
  graphviz \
  && apt-get clean
# [END run_system_package_ubuntu]

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/graphviz-web .

# Run the web service on container startup.
CMD ["/graphviz-web"]
