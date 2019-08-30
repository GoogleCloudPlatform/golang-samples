FROM golang:latest as builder

WORKDIR /src/
COPY . /src/

ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags grpcping \
    -ldflags '-w -extldflags "-static"' \
    -mod vendor \
    -o ./server \
    ./cmd/server/

FROM gcr.io/distroless/static
COPY --from=builder /src/server .

ENTRYPOINT ["/server"]
