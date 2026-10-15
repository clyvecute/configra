FROM golang:1.27-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/configra-api ./cmd/api

FROM alpine:3.22
RUN apk add --no-cache ca-certificates && adduser -D -H -u 10001 app
WORKDIR /app
COPY --from=builder /out/configra-api ./configra-api
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations
COPY --from=builder /app/openapi.yaml ./openapi.yaml
USER 10001
EXPOSE 8080

CMD ["./configra-api"]
