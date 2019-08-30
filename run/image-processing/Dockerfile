# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    https://www.apache.org/licenses/LICENSE-2.0
# 
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Use the offical Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.12 as builder

# Copy local code to the container image.
WORKDIR /app
COPY . .

# Build hellopubsub inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o hellopubsub

# Use a Docker multi-stage build to create a lean production image.
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM alpine:3.10

# [START run_imageproc_dockerfile_imagemagick]

# Install Imagemagick into the container image.
# For more on system packages review the system packages tutorial.
# https://cloud.google.com/run/docs/tutorials/system-packages#dockerfile
RUN apk add --no-cache imagemagick

# [END run_imageproc_dockerfile_imagemagick]

# Install certificates for secure communication with network services.
# For production containers, a single RUN statement should install all system packages.
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/hellopubsub .

# Run the web service on container startup.
CMD ["/hellopubsub"]
