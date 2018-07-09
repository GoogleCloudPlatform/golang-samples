FROM golang:1.8-alpine

RUN apk update && apk add git

RUN go get github.com/garyburd/redigo/redis

ADD . /go/src/visit-counter
RUN go install visit-counter

ENV REDISHOST redis
ENV REDISPORT 6379

ENTRYPOINT /go/bin/visit-counter

EXPOSE 8080
