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

# Use the offical Go image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.24 as builder

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
