# ---- Build Stage ----
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN apk add --no-cache dumb-init build-base
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o /app/server ./cmd/server

# ---- Production Stage ----
FROM alpine:latest

WORKDIR /

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY .env .
COPY --from=builder /usr/bin/dumb-init /usr/bin/dumb-init
COPY --from=builder /app/server /app/server

USER appuser

EXPOSE 4200
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD /app/server
