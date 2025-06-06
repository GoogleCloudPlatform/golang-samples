FROM golang:1.23
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY *.go ./
RUN go build -v -o /usr/local/bin/app .
CMD sh -c "app 2>&1 | tee /var/log/app.log"
