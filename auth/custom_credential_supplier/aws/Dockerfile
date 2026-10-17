# --- Stage 1: Build the application ---
FROM golang:1.21 AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a static binary
# CGO_ENABLED=0 ensures the binary is statically linked and works on Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# --- Stage 2: Run the application ---
FROM alpine:3.18

WORKDIR /app

# Install CA certificates for HTTPS requests to Google APIs
RUN apk add --no-cache ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Run the binary
CMD ["./app"]
