FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum 

RUN go mod download

COPY . .

# Build the application and name the binary as 'service'
RUN CGO_ENABLED=0 GOOS=linux go build -o service cmd/server/main.go

# Use distroless for a minimal runtime environment
FROM gcr.io/distroless/static:nonroot
WORKDIR /

# Copy the 'service' binary from the builder stage
COPY --from=builder /app/service .

USER 65532:65532

# Set the entrypoint to the 'service' binary
ENTRYPOINT ["/service"]
