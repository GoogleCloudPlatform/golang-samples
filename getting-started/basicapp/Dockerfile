FROM golang:1.10

COPY basicwebapp.go .
COPY services/ /go/src/github.com/GoogleCloudPlatform/golang-samples/getting-started/basicapp/services/
RUN go get -d -v ./...
RUN go build basicwebapp.go 

CMD ["./basicwebapp"]
