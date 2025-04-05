
# stage-1: building app binary file
FROM golang:1.24-alpine AS builder 

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/http-server ./cmd/http-server/

# stage-2: creating minimal image size for executable
FROM alpine:3.19

# needed for https/tls
# RUN apk --no-cache add ca-certificates

RUN adduser -D -g '' appuser
USER appuser

COPY --from=builder --chown=appuser:appuser /app/bin/http-server /usr/local/bin/http-server

CMD [ "http-server" ]
