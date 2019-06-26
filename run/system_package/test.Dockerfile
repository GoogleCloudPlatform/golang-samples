# Copyright 2019 Google LLC. All rights reserved.
# Use of this source code is governed by the Apache 2.0
# license that can be found in the LICENSE file.

# Use the offical Golang image to create a build artifact.
# https://hub.docker.com/_/golang
FROM golang

# Install dot. This positions in the same place as Ubuntu.
# [START run_system_package_alpine]
RUN apk --no-cache add graphviz ttf-ubuntu-font-family
# [END run_system_package_alpine]

# Copy local code to the container image.
WORKDIR /app
COPY . .

RUN go test -v .
