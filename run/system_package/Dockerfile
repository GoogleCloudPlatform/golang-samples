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
# https://hub.docker.com/_/golang
FROM golang:1.12 as builder

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
