FROM golang:alpine

WORKDIR /go/src/docdemo
COPY *.go .

RUN apk update \
    && apk add --no-cache git \
    && go get -d ./... \
    && apk del git

RUN go install ./...

CMD ["docdemo"]
