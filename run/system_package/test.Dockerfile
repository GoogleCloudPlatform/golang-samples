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
FROM golang

# Install dot. This positions in the same place as Ubuntu.
# [START run_system_package_alpine]
RUN apk --no-cache add graphviz ttf-ubuntu-font-family
# [END run_system_package_alpine]

# Copy local code to the container image.
WORKDIR /app
COPY . .

RUN go test -v .
