FROM golang:alpine

ARG pkg=github.com/GoogleCloudPlatform/golang-samples/getting-started/bookshelf

RUN apk add --no-cache ca-certificates

COPY . $GOPATH/src/$pkg

RUN set -ex \
      && apk add --no-cache --virtual .build-deps \
              git \
      && go get -v $pkg/... \
      && apk del .build-deps

RUN go install $pkg/...

# Needed for templates for the front-end app.
WORKDIR $GOPATH/src/$pkg/app

# Users of the image should invoke either of the commands.
CMD echo "Use the app or pubsub_worker commands."; exit 1
