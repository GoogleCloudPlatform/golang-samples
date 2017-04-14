FROM golang:alpine

RUN apk add --no-cache ca-certificates

COPY . /go/src/bookshelf

RUN set -ex \
      && apk add --no-cache --virtual .build-deps \
              git \
      && go get -v bookshelf/... \
      && apk del .build-deps

RUN go build -o /app bookshelf/app
RUN go build -o /pubsub_worker bookshelf/pubsub_worker

WORKDIR /go/src/bookshelf/app

CMD ["/app"]
