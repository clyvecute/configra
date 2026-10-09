FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/configra-api ./cmd/api

FROM alpine:3.22
RUN apk add --no-cache ca-certificates && adduser -D -H -u 10001 app
WORKDIR /app
COPY --from=build /out/configra-api ./configra-api
COPY --from=build /app/internal/db/migrations ./internal/db/migrations
COPY --from=build /app/openapi.yaml ./openapi.yaml
USER 10001
EXPOSE 8080

CMD ["./configra-api"]
