FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

RUN adduser -D -s /bin/sh appuser

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN mkdir -p /logs

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o bot cmd/bot/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /app/bot /bot

COPY --from=builder /logs /logs

USER appuser

WORKDIR /

EXPOSE 8080

ENV LOG_LEVEL=info
ENV LOG_CONSOLE=true
ENV LOG_FILE_PATH=/logs/bot.log
ENV LOG_MAX_SIZE=100
ENV LOG_MAX_BACKUPS=3
ENV LOG_MAX_AGE=28
ENV LOG_COMPRESS=true

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/bot", "--health-check"] || exit 1

ENTRYPOINT ["/bot"]
