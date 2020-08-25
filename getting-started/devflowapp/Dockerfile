FROM golang:1.13

COPY devflowapp.go .
COPY services/ /go/src/github.com/GoogleCloudPlatform/golang-samples/getting-started/devflowapp/services/
RUN go get -d -v ./...
RUN go build devflowapp.go 

CMD ["./devflowapp"]
