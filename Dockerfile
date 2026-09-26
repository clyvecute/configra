FROM golang:alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/api

# Ensure migrations are available to the binary
# The binary runs from /app, so ./internal/db/migrations will resolve to /app/internal/db/migrations

CMD ["./main"]
