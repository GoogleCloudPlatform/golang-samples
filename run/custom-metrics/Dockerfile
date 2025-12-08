# Copyright 2023 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      https:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# [START cloudrun_mc_custom_metrics_dockerfile]
FROM golang:1.24 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sample-app

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/sample-app /sample-app
CMD ["/sample-app"]
# [END cloudrun_mc_custom_metrics_dockerfile]
