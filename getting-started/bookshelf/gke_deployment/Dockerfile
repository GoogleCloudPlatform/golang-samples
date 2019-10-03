FROM golang:alpine

ENV GO111MODULE=off
ARG pkg=github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf

RUN apk add --no-cache ca-certificates

RUN set -ex \
      && apk add --no-cache --virtual .build-deps \
              git \
      && go get -v $pkg/... \
      && apk del .build-deps

RUN go get $pkg/...

RUN apk update && \
      apk add git && \
      cd $GOPATH/src/$pkg && \
      git checkout steps

RUN go install $pkg/...

# Needed for templates for the front-end app.
WORKDIR $GOPATH/src/$pkg/app

# Users of the image should invoke either of the commands.
CMD echo "Use the app or pubsub_worker commands."; exit 1
