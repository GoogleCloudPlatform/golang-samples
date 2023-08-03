FROM golang:1.20.7-alpine

WORKDIR /go/src
COPY ./go.mod .
COPY ./go.sum .
RUN mkdir ./hotmid/
COPY ./hotmid/*.go ./hotmid/

RUN go install ./...

CMD ["hotmid"]
