# Cloud Run gRPC Ping Sample Application

For a general introduction to gRPC on Go, see the [official gRPC quickstart](https://grpc.io/docs/quickstart/go/).

## Deploying to Cloud Run

1. Build & Deploy the gRPC services:

   ```sh
   export GOOGLE_CLOUD_PROJECT=[PROJECT_ID]
   # Build and push container images.
   gcloud builds submit

   # Deploy ping service for private access.
   gcloud beta run deploy ping --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-ping

   # Get the host name of the ping service.
   PING_URL=$(gcloud beta run services describe ping --format='value(status.url)')
   PING_DOMAIN=${PING_URL#https://}

   # Deploy ping-relay service for public access.
   gcloud beta run deploy ping-relay --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-ping-relay \
       --set-env-vars GRPC_PING_HOST=${PING_DOMAIN} \
       --allow-unauthenticated
   ```

2. Make a ping request:

   Use the client CLI to send a request:

   ```
   go run . -server [RELAY-SERVICE-DOMAIN]:443  -message "Hello Friend"
   ```

   Use curl to make an HTTP request:
   ```
   curl [RELAY-SERVICE-URL] -d "Howdy HTTP Friends!"
   ```

If you later make some code changes, updating is more concise:

```sh
export GOOGLE_CLOUD_PROJECT=[PROJECT_ID]
gcloud beta run deploy ping --image gcr.io/$PROJECT_ID/grpc-ping
gcloud beta run deploy ping-relay --image gcr.io/$GOOGLE_CLOUD_PROJECT/grpc-ping-relay
```

See below for instructions on updating the proto.

## Environment Variable Configuration Options

* `GRPC_PING_HOST`: [relay: required] Ping service host nanme.
* `GRPC_PING_PORT`: [relay: `443`] Ping service port number.
* `GRPC_PING_INSECURE`: [relay: `false`] Use an insecure connection to the ping service. Primarily for local development.
* `GRPC_PING_UNAUTHENTICATED`: [relay: `false`] Make unauthenticated requests to the ping service. Primarily for local development.

## Building Locally

```sh
docker build --build-arg=SERVICE="ping" -t grpc-ping .
docker build --build-arg=SERVICE="relay" -t grpc-ping-relay .
```

## Running Locally

### Running client &rArr; ping

```sh
cd ping
go run .
```

Open another terminal at the grpc-ping directory:

```
go run . -server localhost:8080 -insecure -message "Hello Friend!"
```

### Running client &rArr; relay &rArr; ping

1. Start the ping service:

   ```sh
   cd ./ping
   PORT=9090 go run .
   ```

2. In another session, start the relay service:

   ```sh
   cd ./relay
   GRPC_PING_INSECURE=1 GRPC_PING_HOST=localhost GRPC_PING_PORT=9090 GRPC_PING_UNAUTHENTICATED=1 \
       go run .
   ```

3. From the grpc-ping directory use the grpc client to send a request:

   ```sh
   go run . -server localhost:8080 -insecure -message "Hello Relayed Friend!"
   ```

   Because **relay** also supports HTTP you can also use `curl`:

   ```sh
   curl http://localhost:8080/ -d "Howdy HTTP Friends!"
   ```

## Updating the Proto

1. Retrieve the protoc plugin for Go:

    ```
    go get -u github.com/golang/protobuf/protoc-gen-go
    ```

2. Modify the Protobuf by editing `api/v1/message.proto`.

3. Regenerate the go code:

    ```
    protoc \
        --proto_path api/v1 \
        --go_out "plugins=grpc:pkg/api/v1" \
        message.proto
    ```

## TODOS

* Confirm cmux approach
* Confirm use of errgroup
* Position region tags
* TODO: inline token creation
* Figure out what to do with streaming responses
* Add e2e tests
* Test CRfA variation with no external connectivity for the ping service.


## Temporary code notes

	"github.com/soheilhy/cmux"
```go
	log.Printf("grpg-ping: starting server...")

	if os.Getenv("GRPC_PING_HOST") == "" {
		log.Fatal("Must specify a ping host with 'GRPC_PING_HOST' environment variable. E.g., example.com")
	}

	var err error
	client, err = Client(os.Getenv("GRPC_PING_HOST"), os.Getenv("GRPC_PING_PORT"), os.Getenv("GRPC_PING_INSECURE") != "")
	if err != nil {
		log.Fatal(err)
	}

	// Determine port for gRPC service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("net.Listen: %v", err)
	}


	// Create a cmux.
	mux := cmux.New(listener)

	// Match connections in order:
	// First grpc, then HTTP, and otherwise Go RPC/TCP.
	httpListener := mux.Match(cmux.HTTP1Fast())
	grpcListener := mux.Match(cmux.Any())

	httpServer := &http.Server{
		Handler: http.HandlerFunc(indexHandler),
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPingServiceServer(grpcServer, &pingService{})

	go httpServer.Serve(httpListener)
	go grpcServer.Serve(grpcListener)

    mux.Serve()
    




	server := http.Server{
		Addr:    ":" + port,
		Handler: withProtocolMultiplex(grpcServer, http.HandlerFunc(indexHandler)),
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2"},
		},
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func withProtocolMultiplex(grpcServer *grpc.Server, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 1 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})
}