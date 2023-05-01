FROM golang:1.20 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o sample-app

FROM alpine:3
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/sample-app /sample-app
CMD ["/sample-app"]
