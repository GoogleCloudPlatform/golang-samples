# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: bin cross upload

VERSION=`git rev-parse HEAD`
DATE=`date +%Y-%m-%d_%H:%M`
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildDate=${DATE}"

bin:
	rm -f ${GOPATH}/bin/gimme-acc
	go install ${LDFLAGS}

cross:
	GOOS=linux   GOARCH=amd64 go build ${LDFLAGS} -o bin/linux_amd64/gimme-acc
	GOOS=darwin  GOARCH=amd64 go build ${LDFLAGS} -o bin/darwin_amd64/gimme-acc
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o bin/windows_amd64/gimme-acc.exe

upload: cross
	gsutil -m cp -a public-read -r bin/* gs://gimme-acc/
