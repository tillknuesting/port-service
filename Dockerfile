FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum 

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o manager cmd/server/main.go

# Use distroless 
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /app/manager .

USER 65532:65532  

ENTRYPOINT ["/manager"]
